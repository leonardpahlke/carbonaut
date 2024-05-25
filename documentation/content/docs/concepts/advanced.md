---
weight: 9
---

## **Carbonaut Advanced Topics**

This page covers several topics to the carbonaut project which offer greater depth to the project.

### **INFO**: Internal State

Carbonaut maintains an internal state which includes data which does not change until a resource was destroyed. Information about how much CPU cores or which Chip Architecture is considered static resource information. Information about the geolocation which indicate where the resource is hosted is considered static environment information. The data schema is defined [here](/docs/reference/schema/#type-staticresdata).

![carbonaut context](/docs/concepts/state.drawio.png)

### **HOW TO**: Add new Plugins

Plugins need to implement a provider interface. There are three providers defined as described in the components documentation [here](/docs/concepts/components/#provider--plugins). As of now, the plugins need to be part of Carbonaut and integrated so they can be discovered. Available plugins are located in the `pkg/plugin/*` folder. If a new plugin should be implemented which integrates another cloud provider like GCP, a new `staticres` plugin needs to be implemented. The plugin is implemented in `pkg/plugin/staticresplugins/gcp` and referenced in `pkg/plugin/staticresplugins/staticresplugin.go`. There are example plugins for each provider type which can be used as a starting point.

### **HOW TO**: Manually build the project

The project pushes binaries to the GitHub release and Docker containers to Docker hub. If you prefer to build binaries manually, you can follow these instructions. In general the entire build process pushing to GitHub and Docker Hub is public and part of the Carbonaut repository.

**Prerequisits install project dependencies**:
* Install dependencies: Install Go and run `make install`


#### Build Go Binaries

#### Build Container Images

#### Setup Testing Infrastructure

See next section "**HOW TO**: Run the end to end scenario" 

### **HOW TO**: Run the end to end scenario

{{< hint info >}}
**TBD**  
documentation not yet added
{{< /hint >}}