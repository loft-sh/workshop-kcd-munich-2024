using TypeSpec.Versioning;
using TypeSpec.Http;
using TypeSpec.Rest;

namespace ComputePlatformService.Models;

model KubernetesNamespace {
  @minLength(1)
  @maxLength(64)
  @key
  name: string;
}

model Kubeconfig {
  @header contentType: "text/yaml";
  @body kubeconfig: bytes;
}

model KubeconfigOptions {
  @query insecure?: boolean;
  @query publicK8sEndpoint?: string;

  /** Base64-encoded certificate data */
  @query certificateAuthorityData?: bytes;
}

@added(Versions.v2)
model vCluster {
  @minLength(1)
  @maxLength(64)
  @key
  name: string;

  version: "0.19.5" | "0.20.0-beta.9";
}

@error
model errorTemplate<Kind extends ApiErrorKind, Code = 400> {
  @statusCode
  status: Code;

  type: Kind;
  title: string;
  detail: string;
  instance: string;
}

enum ApiErrorKind {
  Generic,
  AlreadyExists,
  NotFound,
  Unauthorized,
}

model ErrorGeneric is errorTemplate<ApiErrorKind.Generic, 400>;
model ErrorNotFound is errorTemplate<ApiErrorKind.NotFound, 404>;
model ErrorAlreadyExists is errorTemplate<ApiErrorKind.AlreadyExists, 409>;
model ErrorUnauthorized is errorTemplate<ApiErrorKind.Unauthorized, 401>;
