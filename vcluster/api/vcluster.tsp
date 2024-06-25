using TypeSpec.Versioning;
using TypeSpec.Http;
using TypeSpec.Rest;
using ComputePlatformService.Models;

@route("/vcluster")
namespace ComputePlatformService.Kubernetes.Namespace.vClusterService {
  @added(Versions.v2)
  @readsResource(vCluster)
  op get(
    @path("namespace") ns: string,
  ): vCluster | ErrorUnauthorized | ErrorNotFound;

  @added(Versions.v2)
  @createsResource(vCluster)
  op create(
    @path("namespace") ns: string,
  ): vCluster | ErrorUnauthorized | ErrorNotFound;

  @added(Versions.v2)
  @deletesResource(vCluster)
  op delete(
    @path("namespace") ns: string,
  ): boolean | ErrorUnauthorized | ErrorNotFound;

  /**
   * Get a Kubeconfig to access your vCluster.
   *
   * Will return a Kubeconfig yaml file with
   * access to the namespace.
   */
  @added(Versions.v2)
  @route("/kubeconfig")
  op kubeconfig(
    @path("namespace") ns: string,
    ...KubeconfigOptions,
  ): Kubeconfig | ErrorUnauthorized | ErrorNotFound;
}
