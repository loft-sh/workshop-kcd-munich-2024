package v1

import (
	"bytes"
	"context"

	"github.com/loft-sh/workshop-kcd-munich-2024/vcluster/pkg/namespace"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=./config.yaml ../../tsp-output/@typespec/openapi3/openapi.v1.yaml

// Make sure we conform to ServerInterface
var _ StrictServerInterface = (*ApiService)(nil)

type ApiService struct {
	namespaceService namespace.Service
}

// NamespaceDelete implements StrictServerInterface.
func (c *ApiService) NamespaceDelete(ctx context.Context, request NamespaceDeleteRequestObject) (NamespaceDeleteResponseObject, error) {
	if err := c.namespaceService.Delete(ctx, request.Namespace); err != nil {
		return nil, err
	}

	return NamespaceDelete200JSONResponse(true), nil
}

// NamespaceCreate implements StrictServerInterface.
func (c *ApiService) NamespaceCreate(ctx context.Context, request NamespaceCreateRequestObject) (NamespaceCreateResponseObject, error) {
	ns, err := c.namespaceService.Create(ctx, request.Namespace)
	if err != nil {
		return nil, err
	}

	return NamespaceCreate200JSONResponse{
		Name: ns.Name,
	}, nil
}

// NamespaceKubeconfig implements StrictServerInterface.
func (c *ApiService) NamespaceKubeconfig(ctx context.Context, request NamespaceKubeconfigRequestObject) (NamespaceKubeconfigResponseObject, error) {
	kubeconfigOptions := namespace.KubeconfigOptions{}

	if request.Params.PublicK8sEndpoint != nil {
		kubeconfigOptions.Server = *request.Params.PublicK8sEndpoint
	}
	if request.Params.CertificateAuthorityData != nil {
		kubeconfigOptions.CaData = *request.Params.CertificateAuthorityData
	}
	if request.Params.Insecure != nil {
		kubeconfigOptions.Insecure = *request.Params.Insecure
	}

	kubeconfig, err := c.namespaceService.GetKubeconfig(ctx, request.Namespace, kubeconfigOptions)
	if err != nil {
		return nil, err
	}

	return NamespaceKubeconfig200TextyamlResponse{
		Body:          bytes.NewReader(kubeconfig),
		ContentLength: int64(len(kubeconfig)),
	}, nil
}

// NamespaceGet implements StrictServerInterface.
func (c *ApiService) NamespaceGet(ctx context.Context, request NamespaceGetRequestObject) (NamespaceGetResponseObject, error) {
	ns, err := c.namespaceService.Get(ctx, request.Namespace)
	if err != nil {
		return nil, err
	}

	return NamespaceGet200JSONResponse{
		Name: ns.Name,
	}, nil
}

// KubernetesList implements StrictServerInterface.
func (c *ApiService) KubernetesList(ctx context.Context, request KubernetesListRequestObject) (KubernetesListResponseObject, error) {
	namespaces, err := c.namespaceService.List(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]ModelsKubernetesNamespace, 0, len(namespaces))

	for _, ns := range namespaces {
		response = append(response, ModelsKubernetesNamespace{
			Name: ns.Name,
		})
	}

	return KubernetesList200JSONResponse(response), nil
}

func NewApiService(clientset *kubernetes.Clientset, config *rest.Config) *ApiService {
	return &ApiService{
		namespaceService: namespace.Service{
			Clientset: clientset,
			Config:    config,
		},
	}
}
