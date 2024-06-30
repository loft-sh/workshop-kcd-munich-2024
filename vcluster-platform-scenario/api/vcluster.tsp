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

    @maxLength(64)
    @query
    name?: string,
  ): vCluster | ErrorUnauthorized | ErrorNotFound;

  /** Create a new vCluster or update an existing one */
  @added(Versions.v2)
  @createsOrUpdatesResource(vCluster)
  op create(
    @path("namespace") ns: string,

    @maxLength(64)
    name?: string,

    useLocalChart?: boolean,
    version?: vClusterVersion = vClusterVersion.Beta,
    values?: Record<unknown>,

    /** Wait for the installation to succeed */
    @query wait?: boolean = true,
  ): vCluster | ErrorUnauthorized | ErrorNotFound;

  @added(Versions.v2)
  @deletesResource(vCluster)
  op delete(
    @path("namespace") ns: string,

    @maxLength(64)
    name?: string,

    /** Wait for the deletion to succeed */
    @query wait?: boolean = true,
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

    @maxLength(64)
    @query
    name?: string,
  ): Kubeconfig | ErrorUnauthorized | ErrorNotFound;
}
