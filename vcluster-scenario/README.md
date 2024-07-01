# vCluster Workshops

## Scenario

You are tasked with creating a platform for your developers/tenants/customer to deliver
self-service capabilities while maintaining efficient resource usage.

## Prerequisites

- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [just](https://github.com/casey/just)
- [kind](https://kind.sigs.k8s.io/)
- [cloud-provider-kind](https://github.com/kubernetes-sigs/cloud-provider-kind)
- [jq](https://jqlang.github.io/jq/)
- optionally: [mkcert](https://github.com/FiloSottile/mkcert)

## Setup

### Kind

Ensure you have a local container runtime (Docker Desktop, Orbstack, ...) installed.

Run the following commands:

<details>
  <summary>kind create cluster</summary>

```text
Creating cluster "kind" ...
✓ Ensuring node image (kindest/node:v1.30.0) 🖼
✓ Preparing nodes 📦
✓ Writing configuration 📜
✓ Starting control-plane 🕹️ 
✓ Installing CNI 🔌
✓ Installing StorageClass 💾
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind

Have a nice day! 👋
```

</details>

<details>
  <summary>sudo cloud-provider-kind</summary>

```text
I0701 10:49:53.805843   12214 app.go:41] FLAG: --enable-log-dumping="false"
I0701 10:49:53.805891   12214 app.go:41] FLAG: --logs-dir=""
I0701 10:49:53.805897   12214 app.go:41] FLAG: --v="2"
I0701 10:49:53.954681   12214 controller.go:166] probe HTTP address https://127.0.0.1:62309
I0701 10:49:53.970284   12214 controller.go:83] Creating new cloud provider for cluster kind
I0701 10:49:53.982015   12214 controller.go:90] Starting cloud controller for cluster kind
I0701 10:49:53.982287   12214 controller.go:231] Starting service controller
I0701 10:49:53.982297   12214 node_controller.go:164] Sending events to api server.
I0701 10:49:53.982322   12214 shared_informer.go:313] Waiting for caches to sync for service
I0701 10:49:53.982709   12214 node_controller.go:173] Waiting for informer caches to sync
I0701 10:49:53.987194   12214 reflector.go:359] Caches populated for *v1.Service from pkg/mod/k8s.io/client-go@v0.30.2/tools/cache/reflector.go:232
I0701 10:49:53.987208   12214 reflector.go:359] Caches populated for *v1.Node from pkg/mod/k8s.io/client-go@v0.30.2/tools/cache/reflector.go:232
I0701 10:49:54.083395   12214 shared_informer.go:320] Caches are synced for service
I0701 10:49:54.083461   12214 controller.go:733] Syncing backends for all LB services.
I0701 10:49:54.083472   12214 controller.go:737] Successfully updated 0 out of 0 load balancers to direct traffic to the updated set of nodes
I0701 10:49:54.085527   12214 instances.go:47] Check instance metadata for kind-control-plane
I0701 10:49:54.113552   12214 instances.go:75] instance metadata for kind-control-plane: &cloudprovider.InstanceMetadata{ProviderID:"kind://kind/kind/kind-control-plane", InstanceType:"kind-node", NodeAddresses:[]v1.NodeAddress{v1.NodeAddress{Type:"Hostname", Address:"kind-control-plane"}, v1.NodeAddress{Type:"InternalIP", Address:"192.168.228.2"}, v1.NodeAddress{Type:"InternalIP", Address:"fc00:f853:ccd:e793::2"}}, Zone:"", Region:"", AdditionalLabels:map[string]string(nil)}
I0701 10:49:54.128881   12214 node_controller.go:267] Update 1 nodes status took 45.472667ms.
```

</details>

### ArgoCD

Run the following to setup a local ArgoCD installation.

<details>
  <summary>just setup</summary>

```text
mkdir tmp
mkdir: tmp: File exists
mkcert -install -cert-file tmp/argocd-server.crt -key-file tmp/argocd-server.key argocd-server argocd-server.k8s.orb.local argocd.k8s.orb.local argocd-server.argocd.svc.cluster.local
The local CA is already installed in the system trust store! 👍
The local CA is already installed in the Firefox trust store! 👍


Created a new certificate valid for the following names 📜
 - "argocd-server"
 - "argocd-server.k8s.orb.local"
 - "argocd.k8s.orb.local"
 - "argocd-server.argocd.svc.cluster.local"

The certificate is at "tmp/argocd-server.crt" and the key at "tmp/argocd-server.key" ✅

It will expire on 1 October 2026 🗓

kubectl create namespace argocd
namespace/argocd created
kubectl create -n argocd secret tls argocd-server-tls --cert=./tmp/argocd-server.crt --key=./tmp/argocd-server.key --dry-run=client -o yaml | kubectl apply -f -
secret/argocd-server-tls created
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
zcustomresourcedefinition.apiextensions.k8s.io/applications.argoproj.io created
customresourcedefinition.apiextensions.k8s.io/applicationsets.argoproj.io created
customresourcedefinition.apiextensions.k8s.io/appprojects.argoproj.io created
serviceaccount/argocd-application-controller created
serviceaccount/argocd-applicationset-controller created
serviceaccount/argocd-dex-server created
serviceaccount/argocd-notifications-controller created
serviceaccount/argocd-redis created
serviceaccount/argocd-repo-server created
serviceaccount/argocd-server created
role.rbac.authorization.k8s.io/argocd-application-controller created
role.rbac.authorization.k8s.io/argocd-applicationset-controller created
role.rbac.authorization.k8s.io/argocd-dex-server created
role.rbac.authorization.k8s.io/argocd-notifications-controller created
role.rbac.authorization.k8s.io/argocd-redis created
role.rbac.authorization.k8s.io/argocd-server created
clusterrole.rbac.authorization.k8s.io/argocd-application-controller created
clusterrole.rbac.authorization.k8s.io/argocd-applicationset-controller created
clusterrole.rbac.authorization.k8s.io/argocd-server created
rolebinding.rbac.authorization.k8s.io/argocd-application-controller created
rolebinding.rbac.authorization.k8s.io/argocd-applicationset-controller created
rolebinding.rbac.authorization.k8s.io/argocd-dex-server created
rolebinding.rbac.authorization.k8s.io/argocd-notifications-controller created
rolebinding.rbac.authorization.k8s.io/argocd-redis created
rolebinding.rbac.authorization.k8s.io/argocd-server created
clusterrolebinding.rbac.authorization.k8s.io/argocd-application-controller created
clusterrolebinding.rbac.authorization.k8s.io/argocd-applicationset-controller created
clusterrolebinding.rbac.authorization.k8s.io/argocd-server created
configmap/argocd-cm created
configmap/argocd-cmd-params-cm created
configmap/argocd-gpg-keys-cm created
configmap/argocd-notifications-cm created
configmap/argocd-rbac-cm created
configmap/argocd-ssh-known-hosts-cm created
configmap/argocd-tls-certs-cm created
secret/argocd-notifications-secret created
secret/argocd-secret created
service/argocd-applicationset-controller created
service/argocd-dex-server created
service/argocd-metrics created
service/argocd-notifications-controller-metrics created
service/argocd-redis created
service/argocd-repo-server created
service/argocd-server created
service/argocd-server-metrics created
deployment.apps/argocd-applicationset-controller created
deployment.apps/argocd-dex-server created
deployment.apps/argocd-notifications-controller created
deployment.apps/argocd-redis created
deployment.apps/argocd-repo-server created
deployment.apps/argocd-server created
statefulset.apps/argocd-application-controller created
networkpolicy.networking.k8s.io/argocd-application-controller-network-policy created
networkpolicy.networking.k8s.io/argocd-applicationset-controller-network-policy created
networkpolicy.networking.k8s.io/argocd-dex-server-network-policy created
networkpolicy.networking.k8s.io/argocd-notifications-controller-network-policy created
networkpolicy.networking.k8s.io/argocd-redis-network-policy created
networkpolicy.networking.k8s.io/argocd-repo-server-network-policy created
networkpolicy.networking.k8s.io/argocd-server-network-policy created
```

</details>

## Next Steps

### Namespace as a service

1. Fork the GitHub repo at [ThomasK33/kcd-munich-workshop-gitops](https://github.com/ThomasK33/kcd-munich-workshop-gitops).
2. Edit the applicationset [naas-application-set.yaml](./naas-application-set.yaml)
   to use your own GitHub repo.

3. Run the following commands:

    <details>
      <summary>just namespace-as-a-service</summary>

    ```text
    kubectl apply -f naas-application-set.yaml
    applicationset.argoproj.io/namespace-as-a-service created
    ```

    </details>

4. Access ArgoCD UI using a method of your choice:
   - via [LoadBalancer](https://argo-cd.readthedocs.io/en/stable/getting_started/#service-type-load-balancer)
   - via [Ingress](https://argo-cd.readthedocs.io/en/stable/getting_started/#ingress)
   - via [Port-Forwarding](https://argo-cd.readthedocs.io/en/stable/getting_started/#port-forwarding)

5. Updated the target ref in the applicationset to reference a newer commit (271dd595267b).

6. Reapply the applicationset using `just namespace-as-a-service` and sync via
   the UI for brevity.

### vCluster as a service

1. Apply the vCluster application set using

    <details>
    <summary>just vcluster-as-a-service</summary>

    ```text
    kubectl apply -f vcaas-application-set.yaml
    applicationset.argoproj.io/vcluster-as-a-service created
    ```

    </details>

2. Grab the kubeconfig from the `vc-vcluster` secret in the `dev-thomas` namespace.

    `kubectl get secret vc-vcluster -n dev-thomas -o json | jq -rc '.data.config | @base64d' > vc-kubeconfig.yaml`

3. Access the vCluster using the kubeconfig file via kubectl

    `KUBECOFIG=./vc-kubeconfig.yaml kubectl get nodes`

4. Apply a simple deployment into the vCluster using kubectl

    `KUBECONFIG=./vc-kubeconfig.yaml kubectl apply -f ./vcluster/simple-deployment.yaml`

5. List all pods in the host cluster namespace using

    `kubectl get pods -A`

6. Hand over the `vc-kubeconfig.yaml` to your developer/customer/tenant
