# Kubebuilder

## Scenario

We want to contribute a fix to kubebuilder. To help ourselves and others in the future we'll write a devcontainer.

## Getting started

0. Ensure you have a container engine on your machine. In case you don't, let one of us know
1. Setup DevPod Desktop or CLI, see [docs](https://devpod.sh/docs/getting-started/install)
2. Clone/Fork the [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) project
3. Take a look at the [CONTRIBUTING.md](https://github.com/kubernetes-sigs/kubebuilder/blob/master/CONTRIBUTING.md) as starting point

`devpod up . --recreate --debug`

# Example setup

Check out ./.devcontainer/ for an example configuration including a custome feature.
Copy the folder over into the root of the kubebuilder directory to run it.
