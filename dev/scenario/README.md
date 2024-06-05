# Carbonaut Testing Szenarios

This folder contains information about executing manual testing scenarios. The first scenario focuses on testing the functional workings in carbonaut. The second one about the deployment with Kubernetes. This needs to be executed from project root. Make sure to have environment variables set to be able to use equinix & electricity map APIs.

## Scenario 1:

### Step 1 - Initialization
![step 1 recording](./scenario%201%20results/recordings/s1-1%20recording.mp4)
* `make tf-apply` and hit - enter, enter, enter, yes, navigate to http://localhost:8088/
* `make tf-connection-verify` - we exepct that this fails since resources are not configured yet.
* `make tf-configure` - uses the ansible playbooks 
* `make tf-connection-verify` - should work
* `make tf-stress-test`
* `go run main.go -c dev/config.yaml`
* http://localhost:8088
* `curl localhost:8088/static-data > tmp/s1-s1-state.json`
* `curl localhost:8088/metrics-json > tmp/s1-s1-metrics.json`

### Step 2 Update Empty Config (we do not delete the infrastructure)
![step 2 recording](./scenario%201%20results/recordings/s1-2%20recording.mp4)
* `curl -X POST -H "Content-Type: application/x-yaml" --data-binary @dev/empty-config.yaml http://localhost:8088/load-config`
* `curl localhost:8088/static-data > tmp/s1-s2-state.json`
* `curl localhost:8088/metrics-json > tmp/s1-s2-metrics.json`

### Step 3 Use old config again
![step 3 recording](./scenario%201%20results/recordings/s1-3%20recording.mp4)
* `curl -X POST -H "Content-Type: application/x-yaml" --data-binary @dev/config.yaml http://localhost:8088/load-config`
* `curl localhost:8088/static-data > tmp/s1-s3-state.json`
* `curl localhost:8088/metrics-json > tmp/s1-s3-metrics.json`

### Step 4 - Delete of one ressource
![step 4 recording](./scenario%201%20results/recordings/s1-4%20recording.mp4)
1. `make tf-apply`
   1. enter, enter, 0, yes, navigate equinix console and check that the resource was deleted but the project still exists.
2. `make tf-connection-verify` - nothing should get listed since there are no resources we need to contact
3. `curl localhost:8088/metrics-json > tmp/s1-s4-metrics.json`
4. `curl localhost:8088/static-data > tmp/s1-s4-state.json`

### Step 5 - Creating of two resources
![step 5 recording](./scenario%201%20results/recordings/s1-5%20recording.mp4)
1. `make tf-apply`
   1. enter, enter, 2, yes, navigate equinix console and check that the resources were created.
2. `make tf-configure`
3. `make tf-connection-verify`
4. `curl localhost:8088/metrics-json > tmp/s1-s5-metrics.json`
5. `curl localhost:8088/static-data > tmp/s1-s5-state.json`
