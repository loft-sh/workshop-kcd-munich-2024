import "@typespec/versioning";
import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";
import "./models.tsp";
import "./namespace.tsp";
import "./vcluster.tsp";

using TypeSpec.Versioning;
using TypeSpec.Http;
using TypeSpec.Rest;
using ComputePlatformService.Models;

/**
 * API enabling you to provision compute capacity in a Kubernetes cluster
 *
 * RFC 9457 compatible
 */
@service({
  title: "Compute Platform Service",
})
@versioned(Versions)
@server(
  "http://localhost:8080/{version}",
  "API version",
  {
    /** API version */
    version: Versions = Versions.v1,
  }
)
@useAuth(BearerAuth)
namespace ComputePlatformService;

enum Versions {
  v1,
  v2,
}

@route("/kubernetes")
namespace Kubernetes {
  /**
   * List all Kubernetes namespaces you own
   */
  @listsResource(KubernetesNamespace)
  op list(): KubernetesNamespace[] | ErrorUnauthorized | ErrorGeneric;
}
