using TypeSpec.Versioning;
using TypeSpec.Http;
using TypeSpec.Rest;
using ComputePlatformService.Models;

@route("/{namespace}")
namespace ComputePlatformService.Kubernetes.Namespace {
  /** Get details about a Kubernetes namespace you own */
  @TypeSpec.Rest.readsResource(KubernetesNamespace)
  op get(@path("namespace") ns: string):
    | KubernetesNamespace
    | ErrorNotFound
    | ErrorUnauthorized
    | ErrorGeneric;

  /** Create a Kubernetes namespace */
  @createsResource(KubernetesNamespace)
  op create(
    /** The name of your namespace to be created */
    @path("namespace") name: string,
  ):
    | KubernetesNamespace
    | ErrorAlreadyExists
    | ErrorGeneric
    | ErrorUnauthorized;

  /** Delete a Kubernetes namespace you own */
  @deletesResource(KubernetesNamespace)
  op delete(@path("namespace") ns: string):
    | boolean
    | ErrorNotFound
    | ErrorUnauthorized
    | ErrorGeneric;

  /**
   * Get a Kubeconfig to access your namespace.
   *
   * Will return a Kubeconfig yaml file with
   * access to the namespace.
   */
  @route("/kubeconfig")
  op kubeconfig(
    @path("namespace") ns: string,
    ...KubeconfigOptions,
  ): Kubeconfig | ErrorUnauthorized | ErrorNotFound;
}
