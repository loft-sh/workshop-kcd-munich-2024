package v2

import (
	"bytes"
	"context"

	"github.com/loft-sh/workshop-kcd-munich-2024/vcluster-platform-scenario/pkg/namespace"
	"github.com/loft-sh/workshop-kcd-munich-2024/vcluster-platform-scenario/pkg/vcluster"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=./config.yaml ../../tsp-output/@typespec/openapi3/openapi.v2.yaml

// Make sure we conform to ServerInterface
var _ StrictServerInterface = (*ClusterService)(nil)

type ClusterService struct {
	namespaceService namespace.Service
	vClusterService  vcluster.Service
}

// VClusterServiceCreate implements StrictServerInterface.
func (c *ClusterService) VClusterServiceCreate(ctx context.Context, request VClusterServiceCreateRequestObject) (VClusterServiceCreateResponseObject, error) {
	installOptions := vcluster.VClusterInstallOptions{
		Namespace: request.Namespace,
	}

	if request.Body.Name != nil {
		installOptions.Name = *request.Body.Name
	}

	if request.Body.Values != nil {
		installOptions.Values = *request.Body.Values
	}

	if request.Body.UseLocalChart != nil && *request.Body.UseLocalChart {
		installOptions.UseLocalChart = true
	} else {
		if request.Body.Version != nil {
			switch *request.Body.Version {
			case Stable, N0196:
				installOptions.Version = string(N0196)
			case Beta, N0200Beta11:
				installOptions.Version = string(N0200Beta11)
			}
		}
	}

	if request.Params.Wait != nil {
		installOptions.UseLocalChart = *request.Params.Wait
	}

	vClusterRelease, err := c.vClusterService.Install(ctx, installOptions)
	if err != nil {
		return nil, err
	}

	return VClusterServiceCreate200JSONResponse{
		Name:    vClusterRelease.Name,
		Version: ModelsVClusterVersion(vClusterRelease.Version),
	}, nil
}

// VClusterServiceDelete implements StrictServerInterface.
func (c *ClusterService) VClusterServiceDelete(ctx context.Context, request VClusterServiceDeleteRequestObject) (VClusterServiceDeleteResponseObject, error) {
	name := ""
	if request.Body.Name != nil {
		name = *request.Body.Name
	}

	wait := request.Params.Wait != nil && *request.Params.Wait

	if err := c.vClusterService.Uninstall(ctx, request.Namespace, name, wait); err != nil {
		return nil, err
	}

	return VClusterServiceDelete200JSONResponse(true), nil
}

// VClusterServiceKubeconfig implements StrictServerInterface.
func (c *ClusterService) VClusterServiceKubeconfig(ctx context.Context, request VClusterServiceKubeconfigRequestObject) (VClusterServiceKubeconfigResponseObject, error) {
	name := ""
	if request.Params.Name != nil {
		name = *request.Params.Name
	}

	kubeconfig, err := c.vClusterService.Kubeconfig(ctx, request.Namespace, name)
	if err != nil {
		return nil, err
	}

	return VClusterServiceKubeconfig200TextyamlResponse{
		Body:          bytes.NewReader(kubeconfig),
		ContentLength: int64(len(kubeconfig)),
	}, nil
}

// VClusterServiceGet implements StrictServerInterface.
func (c *ClusterService) VClusterServiceGet(ctx context.Context, request VClusterServiceGetRequestObject) (VClusterServiceGetResponseObject, error) {
	name := ""
	if request.Params.Name != nil {
		name = *request.Params.Name
	}

	release, values, err := c.vClusterService.Get(ctx, request.Namespace, name)
	if err != nil {
		return nil, err
	}

	return VClusterServiceGet200JSONResponse{
		Name:    release.Name,
		Version: ModelsVClusterVersion(release.Version),
		Values:  &values,
	}, nil
}

// KubernetesList implements StrictServerInterface.
func (c *ClusterService) KubernetesList(ctx context.Context, request KubernetesListRequestObject) (KubernetesListResponseObject, error) {
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

// NamespaceCreate implements StrictServerInterface.
func (c *ClusterService) NamespaceCreate(ctx context.Context, request NamespaceCreateRequestObject) (NamespaceCreateResponseObject, error) {
	ns, err := c.namespaceService.Create(ctx, request.Namespace)
	if err != nil {
		return nil, err
	}

	return NamespaceCreate200JSONResponse{
		Name: ns.Name,
	}, nil
}

// NamespaceDelete implements StrictServerInterface.
func (c *ClusterService) NamespaceDelete(ctx context.Context, request NamespaceDeleteRequestObject) (NamespaceDeleteResponseObject, error) {
	if err := c.namespaceService.Delete(ctx, request.Namespace); err != nil {
		return nil, err
	}

	return NamespaceDelete200JSONResponse(true), nil
}

// NamespaceGet implements StrictServerInterface.
func (c *ClusterService) NamespaceGet(ctx context.Context, request NamespaceGetRequestObject) (NamespaceGetResponseObject, error) {
	ns, err := c.namespaceService.Get(ctx, request.Namespace)
	if err != nil {
		return nil, err
	}

	return NamespaceGet200JSONResponse{
		Name: ns.Name,
	}, nil
}

// NamespaceKubeconfig implements StrictServerInterface.
func (c *ClusterService) NamespaceKubeconfig(ctx context.Context, request NamespaceKubeconfigRequestObject) (NamespaceKubeconfigResponseObject, error) {
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

func NewClusterService(clientset *kubernetes.Clientset, config *rest.Config) *ClusterService {
	return &ClusterService{
		namespaceService: namespace.Service{
			Clientset: clientset,
			Config:    config,
		},
		vClusterService: vcluster.Service{
			Clientset: clientset,
		},
	}
}
