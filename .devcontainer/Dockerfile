FROM buildpack-deps:noble-curl

LABEL dev.containers.features="common"

ARG TARGETOS
ARG TARGETARCH

# Install tools
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/${TARGETARCH}/kubectl" && install -c -m 0755 kubectl /usr/local/bin && rm kubectl
RUN curl -L -o helm3.tar.gz https://github.com/loft-sh/helm/releases/download/v3.14-loft/helm-v3.14-loft-${TARGETOS}-${TARGETARCH}.tar.gz && tar -zxvf helm3.tar.gz linux-${TARGETARCH}/helm && install -c -m 0755 linux-${TARGETARCH}/helm /usr/local/bin/helm && rm helm3.tar.gz && rm -R linux-${TARGETARCH}
RUN curl -L -o devspace "https://github.com/loft-sh/devspace/releases/latest/download/devspace-linux-${TARGETARCH}" && install -c -m 0755 devspace /usr/local/bin && rm devspace
RUN curl -L -o vcluster "https://github.com/loft-sh/vcluster/releases/latest/download/vcluster-linux-${TARGETARCH}" && install -c -m 0755 vcluster /usr/local/bin && rm vcluster
RUN curl -L -o devpod "https://github.com/loft-sh/devpod/releases/latest/download/devpod-linux-${TARGETARCH}" && install -c -m 0755 devpod /usr/local/bin && rm devpod
RUN curl -L -o kind "https://kind.sigs.k8s.io/dl/v0.17.0/kind-linux-${TARGETARCH}" && install -c -m 0755 kind /usr/local/bin && rm kind
RUN curl -L -o loft "https://github.com/loft-sh/loft/releases/latest/download/loft-linux-${TARGETARCH}" && install -c -m 0755 loft /usr/local/bin && rm loft

RUN apt-get update && apt-get -y install --no-install-recommends \
        jq \
        yq \
        xdg-utils
