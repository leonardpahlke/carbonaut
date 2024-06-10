# Carbonaut Testing Szenarios

This folder contains information about executing manual testing scenarios. The first scenario focuses on testing the functional workings in carbonaut. The second one about the deployment with Kubernetes. This needs to be executed from project root. Make sure to have environment variables set to be able to use equinix & electricity map APIs.

Execute scenarions by running `make test-scenario-1` and `test-scenario-2`.
These scenarios will create files in the `test-scenario/results-1` and `test-scenario/results-2` folders.

**Requirements**:
* Run on Linux or Macos (not tested on Windows, might work with WSL)
* Likely alredy installed dependencies:
  * [`OpenSSH`](https://www.openssh.com/)
  * [`curl`](https://curl.se/)
  * [`jq`](https://jqlang.github.io/jq/)
  * [`Go`](https://go.dev/)
  * [`Make`](https://www.gnu.org/software/make/)
  * `osascript` - this may already be installed or cause problems. It is used to start another terminal window in the background. This makes it easy to automate the entire process. For example, during the execution, processes like port forwarding need to run until the scenario ends. The port forwarding would get started in another terminal. osascript does this.
* Additional requirements:
  * [`OpenTofu`](https://opentofu.org/)
  * [`kubectl`](https://kubernetes.io/docs/reference/kubectl/)
  * [`minikube`](https://minikube.sigs.k8s.io/)
  * [`Docker`](https://www.docker.com/) 
* Environment variables:
  * `METAL_AUTH_TOKEN` to access Equinix infrastructure ([equinix metal].(https://deploy.equinix.com/))
  * `ELECTRICITY_MAP_AUTH_TOKEN` to access [electricity map](https://app.electricitymaps.com/map).
