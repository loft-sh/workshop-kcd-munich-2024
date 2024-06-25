package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	apiv1 "github.com/loft-sh/workshop-kcd-munich-2024/vcluster/api/v1"
	apiv2 "github.com/loft-sh/workshop-kcd-munich-2024/vcluster/api/v2"
	"github.com/loft-sh/workshop-kcd-munich-2024/vcluster/pkg/auth"
	middleware "github.com/oapi-codegen/nethttp-middleware"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type apiVersion struct {
	GetSwagger    func() (swagger *openapi3.T, err error)
	ServiceRouter func(*http.ServeMux)
	Version       string
}

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	port := flag.String("port", "8080", "port for HTTP server")
	flag.Parse()

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		if !errors.Is(err, rest.ErrNotInCluster) {
			panic(err.Error())
		}

		// use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	h := http.NewServeMux()

	for _, apiVersion := range []apiVersion{
		{
			Version:    "v1",
			GetSwagger: apiv1.GetSwagger,
			ServiceRouter: func(r *http.ServeMux) {
				strictClusterServiceHandler := apiv1.NewStrictHandlerWithOptions(apiv1.NewApiService(clientset, config), nil, apiv1.StrictHTTPServerOptions{
					RequestErrorHandlerFunc:  errorHandlerFunc,
					ResponseErrorHandlerFunc: errorHandlerFunc,
				})

				apiv1.HandlerFromMux(strictClusterServiceHandler, r)
			},
		},
		{
			Version:    "v2",
			GetSwagger: apiv2.GetSwagger,
			ServiceRouter: func(r *http.ServeMux) {
				// Create an instance of our handler which satisfies the generated interface
				strictClusterServiceHandler := apiv2.NewStrictHandler(apiv2.NewClusterService(clientset, config), nil)

				// We now register our petStore above as the handler for the interface
				apiv2.HandlerFromMux(strictClusterServiceHandler, r)
			},
		},
	} {
		if err := registerApiVersion(h, apiVersion); err != nil {
			fmt.Fprintf(os.Stderr, "Error registering api version %s: %v", apiVersion.Version, err)
			os.Exit(1)
		}
	}

	h.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello there")
	})

	s := &http.Server{
		Handler: h,
		Addr:    net.JoinHostPort("0.0.0.0", *port),
	}

	// And we serve HTTP until the world ends.
	log.Fatal(s.ListenAndServe())
}

func registerApiVersion(h *http.ServeMux, api apiVersion) error {
	version := api.Version
	getSwagger := api.GetSwagger
	serviceRouter := api.ServiceRouter

	swagger, err := getSwagger()
	if err != nil {
		return fmt.Errorf("loading openapi spec: %w", err)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	r := http.NewServeMux()

	serviceRouter(r)

	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	validatedApiHandler := middleware.OapiRequestValidatorWithOptions(swagger, &middleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
				if input == nil ||
					input.RequestValidationInput == nil ||
					input.RequestValidationInput.Request == nil ||
					input.RequestValidationInput.Request.Header == nil {
					return errors.New("invalid input")
				}

				authorizationHeader := strings.Split(input.RequestValidationInput.Request.Header.Get("Authorization"), " ")
				if len(authorizationHeader) != 2 {
					return errors.New("incorrect authorization header")
				}

				user, err := auth.UserFromBearerToken(authorizationHeader[1])
				if err != nil {
					return err
				}

				authCtx := auth.NewContext(input.RequestValidationInput.Request.Context(), &user)

				*input.RequestValidationInput.Request = *input.RequestValidationInput.Request.WithContext(authCtx)

				return nil
			},
		},
		ErrorHandler: func(w http.ResponseWriter, message string, statusCode int) {
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(statusCode)

			errorType := string(apiv1.Generic)

			if statusCode == 401 {
				errorType = string(apiv1.Unauthorized)
			}

			response := struct {
				Type     string `json:"type"`
				Title    string `json:"title"`
				Detail   string `json:"detail"`
				Instance string `json:"instance"`
				Status   int    `json:"status"`
			}{
				Type:     errorType,
				Title:    message,
				Status:   statusCode,
				Detail:   message,
				Instance: "",
			}

			json.NewEncoder(w).Encode(response)
		},
	})(r)

	swagger, err = getSwagger()
	if err != nil {
		return fmt.Errorf("load openapi: %w", err)
	}

	swagger.Servers[0].Variables["version"].Default = version

	specContent, err := json.Marshal(swagger)
	if err != nil {
		return fmt.Errorf("marshal openapi: %w", err)
	}

	handler := http.NewServeMux()

	handler.HandleFunc("GET /openapi.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		fmt.Fprint(w, string(specContent))
	})

	handler.HandleFunc("GET /reference", func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecContent: string(specContent),
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Compute Platform Service",
			},
			DarkMode: true,
		})
		if err != nil {
			fmt.Printf("%v", err)
		}

		fmt.Fprintln(w, htmlContent)
	})
	handler.Handle("/", validatedApiHandler)

	h.Handle(fmt.Sprintf("/%v/", version), http.StripPrefix(fmt.Sprintf("/%v", version), handler))

	return nil
}

func errorHandlerFunc(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/problem+json")

	status := 400
	errorType := string(apiv1.Generic)

	if errors.Is(err, auth.ErrUnauthorized) {
		status = 401
		errorType = string(apiv1.Unauthorized)
	}

	if kerrors.IsAlreadyExists(err) {
		status = 409
		errorType = string(apiv1.AlreadyExists)
	}

	if kerrors.IsNotFound(err) {
		status = 404
		errorType = string(apiv1.NotFound)
	}

	w.WriteHeader(status)

	response := struct {
		Type     string `json:"type"`
		Title    string `json:"title"`
		Detail   string `json:"detail"`
		Instance string `json:"instance"`
		Status   int64  `json:"status"`
	}{
		Type:     errorType,
		Title:    err.Error(),
		Status:   int64(status),
		Detail:   err.Error(),
		Instance: r.URL.Path,
	}

	json.NewEncoder(w).Encode(response)
}
