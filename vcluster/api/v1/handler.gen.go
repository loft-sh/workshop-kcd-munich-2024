//go:build go1.22

// Package v1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package v1

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oapi-codegen/runtime"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for ModelsErrorAlreadyExistsType.
const (
	AlreadyExists ModelsErrorAlreadyExistsType = "AlreadyExists"
)

// Defines values for ModelsErrorGenericType.
const (
	Generic ModelsErrorGenericType = "Generic"
)

// Defines values for ModelsErrorNotFoundType.
const (
	NotFound ModelsErrorNotFoundType = "NotFound"
)

// Defines values for ModelsErrorUnauthorizedType.
const (
	Unauthorized ModelsErrorUnauthorizedType = "Unauthorized"
)

// ModelsErrorAlreadyExists defines model for Models.ErrorAlreadyExists.
type ModelsErrorAlreadyExists struct {
	Detail   string                       `json:"detail"`
	Instance string                       `json:"instance"`
	Title    string                       `json:"title"`
	Type     ModelsErrorAlreadyExistsType `json:"type"`
}

// ModelsErrorAlreadyExistsType defines model for ModelsErrorAlreadyExists.Type.
type ModelsErrorAlreadyExistsType string

// ModelsErrorGeneric defines model for Models.ErrorGeneric.
type ModelsErrorGeneric struct {
	Detail   string                 `json:"detail"`
	Instance string                 `json:"instance"`
	Title    string                 `json:"title"`
	Type     ModelsErrorGenericType `json:"type"`
}

// ModelsErrorGenericType defines model for ModelsErrorGeneric.Type.
type ModelsErrorGenericType string

// ModelsErrorNotFound defines model for Models.ErrorNotFound.
type ModelsErrorNotFound struct {
	Detail   string                  `json:"detail"`
	Instance string                  `json:"instance"`
	Title    string                  `json:"title"`
	Type     ModelsErrorNotFoundType `json:"type"`
}

// ModelsErrorNotFoundType defines model for ModelsErrorNotFound.Type.
type ModelsErrorNotFoundType string

// ModelsErrorUnauthorized defines model for Models.ErrorUnauthorized.
type ModelsErrorUnauthorized struct {
	Detail   string                      `json:"detail"`
	Instance string                      `json:"instance"`
	Title    string                      `json:"title"`
	Type     ModelsErrorUnauthorizedType `json:"type"`
}

// ModelsErrorUnauthorizedType defines model for ModelsErrorUnauthorized.Type.
type ModelsErrorUnauthorizedType string

// ModelsKubernetesNamespace defines model for Models.KubernetesNamespace.
type ModelsKubernetesNamespace struct {
	Name string `json:"name"`
}

// ModelsKubeconfigOptionsCertificateAuthorityData defines model for Models.KubeconfigOptions.certificateAuthorityData.
type ModelsKubeconfigOptionsCertificateAuthorityData = []byte

// ModelsKubeconfigOptionsInsecure defines model for Models.KubeconfigOptions.insecure.
type ModelsKubeconfigOptionsInsecure = bool

// ModelsKubeconfigOptionsPublicK8sEndpoint defines model for Models.KubeconfigOptions.publicK8sEndpoint.
type ModelsKubeconfigOptionsPublicK8sEndpoint = string

// NamespaceKubeconfigParams defines parameters for NamespaceKubeconfig.
type NamespaceKubeconfigParams struct {
	Insecure          *ModelsKubeconfigOptionsInsecure          `form:"insecure,omitempty" json:"insecure,omitempty"`
	PublicK8sEndpoint *ModelsKubeconfigOptionsPublicK8sEndpoint `form:"publicK8sEndpoint,omitempty" json:"publicK8sEndpoint,omitempty"`

	// CertificateAuthorityData Base64-encoded certificate data
	CertificateAuthorityData *ModelsKubeconfigOptionsCertificateAuthorityData `form:"certificateAuthorityData,omitempty" json:"certificateAuthorityData,omitempty"`
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /kubernetes)
	KubernetesList(w http.ResponseWriter, r *http.Request)

	// (DELETE /kubernetes/{namespace})
	NamespaceDelete(w http.ResponseWriter, r *http.Request, namespace string)

	// (GET /kubernetes/{namespace})
	NamespaceGet(w http.ResponseWriter, r *http.Request, namespace string)

	// (POST /kubernetes/{namespace})
	NamespaceCreate(w http.ResponseWriter, r *http.Request, namespace string)

	// (GET /kubernetes/{namespace}/kubeconfig)
	NamespaceKubeconfig(w http.ResponseWriter, r *http.Request, namespace string, params NamespaceKubeconfigParams)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// KubernetesList operation middleware
func (siw *ServerInterfaceWrapper) KubernetesList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.KubernetesList(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// NamespaceDelete operation middleware
func (siw *ServerInterfaceWrapper) NamespaceDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "namespace" -------------
	var namespace string

	err = runtime.BindStyledParameterWithOptions("simple", "namespace", r.PathValue("namespace"), &namespace, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "namespace", Err: err})
		return
	}

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.NamespaceDelete(w, r, namespace)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// NamespaceGet operation middleware
func (siw *ServerInterfaceWrapper) NamespaceGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "namespace" -------------
	var namespace string

	err = runtime.BindStyledParameterWithOptions("simple", "namespace", r.PathValue("namespace"), &namespace, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "namespace", Err: err})
		return
	}

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.NamespaceGet(w, r, namespace)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// NamespaceCreate operation middleware
func (siw *ServerInterfaceWrapper) NamespaceCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "namespace" -------------
	var namespace string

	err = runtime.BindStyledParameterWithOptions("simple", "namespace", r.PathValue("namespace"), &namespace, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "namespace", Err: err})
		return
	}

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.NamespaceCreate(w, r, namespace)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// NamespaceKubeconfig operation middleware
func (siw *ServerInterfaceWrapper) NamespaceKubeconfig(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "namespace" -------------
	var namespace string

	err = runtime.BindStyledParameterWithOptions("simple", "namespace", r.PathValue("namespace"), &namespace, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "namespace", Err: err})
		return
	}

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params NamespaceKubeconfigParams

	// ------------- Optional query parameter "insecure" -------------

	err = runtime.BindQueryParameter("form", true, false, "insecure", r.URL.Query(), &params.Insecure)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "insecure", Err: err})
		return
	}

	// ------------- Optional query parameter "publicK8sEndpoint" -------------

	err = runtime.BindQueryParameter("form", true, false, "publicK8sEndpoint", r.URL.Query(), &params.PublicK8sEndpoint)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "publicK8sEndpoint", Err: err})
		return
	}

	// ------------- Optional query parameter "certificateAuthorityData" -------------

	err = runtime.BindQueryParameter("form", true, false, "certificateAuthorityData", r.URL.Query(), &params.CertificateAuthorityData)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "certificateAuthorityData", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.NamespaceKubeconfig(w, r, namespace, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{})
}

type StdHTTPServerOptions struct {
	BaseURL          string
	BaseRouter       *http.ServeMux
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, m *http.ServeMux) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseRouter: m,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, m *http.ServeMux, baseURL string) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseURL:    baseURL,
		BaseRouter: m,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options StdHTTPServerOptions) http.Handler {
	m := options.BaseRouter

	if m == nil {
		m = http.NewServeMux()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	m.HandleFunc("GET "+options.BaseURL+"/kubernetes", wrapper.KubernetesList)
	m.HandleFunc("DELETE "+options.BaseURL+"/kubernetes/{namespace}", wrapper.NamespaceDelete)
	m.HandleFunc("GET "+options.BaseURL+"/kubernetes/{namespace}", wrapper.NamespaceGet)
	m.HandleFunc("POST "+options.BaseURL+"/kubernetes/{namespace}", wrapper.NamespaceCreate)
	m.HandleFunc("GET "+options.BaseURL+"/kubernetes/{namespace}/kubeconfig", wrapper.NamespaceKubeconfig)

	return m
}

type KubernetesListRequestObject struct {
}

type KubernetesListResponseObject interface {
	VisitKubernetesListResponse(w http.ResponseWriter) error
}

type KubernetesList200JSONResponse []ModelsKubernetesNamespace

func (response KubernetesList200JSONResponse) VisitKubernetesListResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type KubernetesList400JSONResponse ModelsErrorGeneric

func (response KubernetesList400JSONResponse) VisitKubernetesListResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type KubernetesList401JSONResponse ModelsErrorUnauthorized

func (response KubernetesList401JSONResponse) VisitKubernetesListResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type NamespaceDeleteRequestObject struct {
	Namespace string `json:"namespace"`
}

type NamespaceDeleteResponseObject interface {
	VisitNamespaceDeleteResponse(w http.ResponseWriter) error
}

type NamespaceDelete200JSONResponse bool

func (response NamespaceDelete200JSONResponse) VisitNamespaceDeleteResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type NamespaceDelete400JSONResponse ModelsErrorGeneric

func (response NamespaceDelete400JSONResponse) VisitNamespaceDeleteResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type NamespaceDelete401JSONResponse ModelsErrorUnauthorized

func (response NamespaceDelete401JSONResponse) VisitNamespaceDeleteResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type NamespaceDelete404JSONResponse ModelsErrorNotFound

func (response NamespaceDelete404JSONResponse) VisitNamespaceDeleteResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type NamespaceGetRequestObject struct {
	Namespace string `json:"namespace"`
}

type NamespaceGetResponseObject interface {
	VisitNamespaceGetResponse(w http.ResponseWriter) error
}

type NamespaceGet200JSONResponse ModelsKubernetesNamespace

func (response NamespaceGet200JSONResponse) VisitNamespaceGetResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type NamespaceGet400JSONResponse ModelsErrorGeneric

func (response NamespaceGet400JSONResponse) VisitNamespaceGetResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type NamespaceGet401JSONResponse ModelsErrorUnauthorized

func (response NamespaceGet401JSONResponse) VisitNamespaceGetResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type NamespaceGet404JSONResponse ModelsErrorNotFound

func (response NamespaceGet404JSONResponse) VisitNamespaceGetResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type NamespaceCreateRequestObject struct {
	Namespace string `json:"namespace"`
}

type NamespaceCreateResponseObject interface {
	VisitNamespaceCreateResponse(w http.ResponseWriter) error
}

type NamespaceCreate200JSONResponse ModelsKubernetesNamespace

func (response NamespaceCreate200JSONResponse) VisitNamespaceCreateResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type NamespaceCreate400JSONResponse ModelsErrorGeneric

func (response NamespaceCreate400JSONResponse) VisitNamespaceCreateResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type NamespaceCreate401JSONResponse ModelsErrorUnauthorized

func (response NamespaceCreate401JSONResponse) VisitNamespaceCreateResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type NamespaceCreate409JSONResponse ModelsErrorAlreadyExists

func (response NamespaceCreate409JSONResponse) VisitNamespaceCreateResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(409)

	return json.NewEncoder(w).Encode(response)
}

type NamespaceKubeconfigRequestObject struct {
	Namespace string `json:"namespace"`
	Params    NamespaceKubeconfigParams
}

type NamespaceKubeconfigResponseObject interface {
	VisitNamespaceKubeconfigResponse(w http.ResponseWriter) error
}

type NamespaceKubeconfig200TextyamlResponse struct {
	Body          io.Reader
	ContentLength int64
}

func (response NamespaceKubeconfig200TextyamlResponse) VisitNamespaceKubeconfigResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/yaml")
	if response.ContentLength != 0 {
		w.Header().Set("Content-Length", fmt.Sprint(response.ContentLength))
	}
	w.WriteHeader(200)

	if closer, ok := response.Body.(io.ReadCloser); ok {
		defer closer.Close()
	}
	_, err := io.Copy(w, response.Body)
	return err
}

type NamespaceKubeconfig401JSONResponse ModelsErrorUnauthorized

func (response NamespaceKubeconfig401JSONResponse) VisitNamespaceKubeconfigResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type NamespaceKubeconfig404JSONResponse ModelsErrorNotFound

func (response NamespaceKubeconfig404JSONResponse) VisitNamespaceKubeconfigResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {

	// (GET /kubernetes)
	KubernetesList(ctx context.Context, request KubernetesListRequestObject) (KubernetesListResponseObject, error)

	// (DELETE /kubernetes/{namespace})
	NamespaceDelete(ctx context.Context, request NamespaceDeleteRequestObject) (NamespaceDeleteResponseObject, error)

	// (GET /kubernetes/{namespace})
	NamespaceGet(ctx context.Context, request NamespaceGetRequestObject) (NamespaceGetResponseObject, error)

	// (POST /kubernetes/{namespace})
	NamespaceCreate(ctx context.Context, request NamespaceCreateRequestObject) (NamespaceCreateResponseObject, error)

	// (GET /kubernetes/{namespace}/kubeconfig)
	NamespaceKubeconfig(ctx context.Context, request NamespaceKubeconfigRequestObject) (NamespaceKubeconfigResponseObject, error)
}

type StrictHandlerFunc = strictnethttp.StrictHTTPHandlerFunc
type StrictMiddlewareFunc = strictnethttp.StrictHTTPMiddlewareFunc

type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}}
}

func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}

// KubernetesList operation middleware
func (sh *strictHandler) KubernetesList(w http.ResponseWriter, r *http.Request) {
	var request KubernetesListRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.KubernetesList(ctx, request.(KubernetesListRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "KubernetesList")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(KubernetesListResponseObject); ok {
		if err := validResponse.VisitKubernetesListResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// NamespaceDelete operation middleware
func (sh *strictHandler) NamespaceDelete(w http.ResponseWriter, r *http.Request, namespace string) {
	var request NamespaceDeleteRequestObject

	request.Namespace = namespace

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.NamespaceDelete(ctx, request.(NamespaceDeleteRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "NamespaceDelete")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(NamespaceDeleteResponseObject); ok {
		if err := validResponse.VisitNamespaceDeleteResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// NamespaceGet operation middleware
func (sh *strictHandler) NamespaceGet(w http.ResponseWriter, r *http.Request, namespace string) {
	var request NamespaceGetRequestObject

	request.Namespace = namespace

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.NamespaceGet(ctx, request.(NamespaceGetRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "NamespaceGet")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(NamespaceGetResponseObject); ok {
		if err := validResponse.VisitNamespaceGetResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// NamespaceCreate operation middleware
func (sh *strictHandler) NamespaceCreate(w http.ResponseWriter, r *http.Request, namespace string) {
	var request NamespaceCreateRequestObject

	request.Namespace = namespace

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.NamespaceCreate(ctx, request.(NamespaceCreateRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "NamespaceCreate")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(NamespaceCreateResponseObject); ok {
		if err := validResponse.VisitNamespaceCreateResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// NamespaceKubeconfig operation middleware
func (sh *strictHandler) NamespaceKubeconfig(w http.ResponseWriter, r *http.Request, namespace string, params NamespaceKubeconfigParams) {
	var request NamespaceKubeconfigRequestObject

	request.Namespace = namespace
	request.Params = params

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.NamespaceKubeconfig(ctx, request.(NamespaceKubeconfigRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "NamespaceKubeconfig")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(NamespaceKubeconfigResponseObject); ok {
		if err := validResponse.VisitNamespaceKubeconfigResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xYzXLbNhB+FQzaIysqqZs6vDmO4/EkTT1NOz3YPkDgUkIKAgywUMxq9O4dABJ/RqSV",
	"2KNxDzppSGJ3v/35dqFdUa7LSitQaGm2ohUzrAQEE55+0zlIO3nvZsC1KsT89wqFVnbCwaAoBGcIZw4X",
	"2gis3zJkXigHy40IB2lG3zALr05+AsV1DjnpCJLcCyRU+GNfHJiaJlSxEmhGR/Un1PIFlMFQoU3JkGZ0",
	"ViPQhGJdeVmLRqg5Xa+TcfxCWeDOgFczZL/53rW30T/TWgJTDxuo3EwK/v7UXqi80kLhmKXdgwMmG5fW",
	"24/d9FwYo82ZNMDy+uJe2E0mja58FMHGpCATckClj79FpjgMfkSBcuRLeLGioFxJsxvaB3A3lA4DX5ww",
	"kPvj4etWf7LF10HTatCzz8CRtuEODl+CAiP4c7i6NX14Jz9qfKedyp/Dy8b24d38S7HI8n/hWVzt2T+o",
	"u75TGAUI9iMrwVYsYu87HFvDipbs/gOoOS5o9uokoaVQ28cX+zAGFbs4fP/wjU1g/cn3kWjvDTADxjda",
	"/xQaTGh04XXbWBeIVexBQhV6t9WfXV8RUGwmhZqTWjuCmlRGL4UVWhE/ZxwC4axiXGBNhCKMtPEgXDqL",
	"YG7Vrfrj3Tl5ffLLr0GIoZjJNuAZPd9oupYM/Qwgn8AsBfdHlmBsBLN84aOuK1CsEjSjP0+mkylNaMVw",
	"EZxO/2lM+8c54K5HH4RFwqTswlTbvNngo/6qaLBjmBe6ymlG29NeAfWZsZVWNlp6OZ36H64VQpwLrKqk",
	"n3VCq/Sz9ZZXnRkgEMog+KOBgmb0h7Qd2elmGqQPlFdT75QZw+qYwb6ffy6A+OoBi2TBLLGOc4Ac8omX",
	"PvlOwN+As9fBRwBZMEswhGsnc6I0EqdyMJ5iOcEO4NyBrzShlkyKnNhaIbvfAH9xEOC9djGA/oxzsJYI",
	"S1zn5MQfXSfdwktXTTWtY/VJQNitw7fhfZ8ujeRoGTYlEMVD8be3u5vNfcQTor2ONFppt5+gcfDQteTu",
	"iSW+c7M6VugzVaiHdHIQSM11Yk8wmfKRLEQ/hpATA1Y7w2ESNQw27EtAEuexJWymHT6aM5eA/2PCPHYU",
	"HIl1JNY+YlXaDjDr3AAbHUHjNIpiu0zahemVEV14ZpoOUVGTGRAe1OTbbcGRgkcKPo6Crw8Cqb/72JNi",
	"rlUhBUdLvgpchGByZwwoJBY9x3QRXsYM7Lk1htdx9TT6L8YPxcjbeNDni8Ug9ck28f+9/hZSEgPojOpL",
	"1ayUpBASAu5btVGBOqBtlYw3g1bZAUdrMpy21ly6fyv4FCW7C72naBtdhH5D/0K4x9RnrV/U7dpUKBZW",
	"kQNbxu/rUMdZ+6hZ213HBBJ0FzE3dz7FUeHQzDy7viLbdUdCnZGbDU2WplJzJhfaYnY6PZ2mq82xNU3o",
	"khnBZjKWTLMt8boL5iTGxUnyoKnt0iwcXL6kd+v1+m79XwAAAP//cW2fFkkYAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}