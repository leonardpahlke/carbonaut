---
title: Carbonaut Components
type: docs
weight: 1
slug: components
---

<div style="display: flex; justify-content: space-between;">
  <div style="flex: 2;">
    <p>
      Carbonaut's internal components are explained in this document. At a high level, Carbonaut serves metrics and integration functionallity over a server component which calls endpoints in the connector to collect metrics. The connector is the central component which controls the main application life cycle, it inhales the logic registering plugins and managing state. Plugins implement providers which deliver different kinds of data which is required to get the holistic overview of cloud native environmental sustainability.
    </p>
  </div>
  <div style="flex: 1; padding-left: 20px;">
    <img src="/logo.png" alt="Alternative Text" style="width: 100%;">
  </div>
</div>

![carbonaut context](/concepts/context.drawio.png)

<!-- The next sections will focus on each building block as well as the data schema. -->

### Server

The Carbonaut Server hosts an HTTP server which serves collected metrics and serves an access point to configure the deployment at runtime. Detailed information can be found in the [api documentation](/docs/reference/server-api/). As of now the server exposes metrics in `json` format. These metrics could be mapped to [prometheus metric types](https://prometheus.io/docs/concepts/metric_types/) down the road for better integration in the cloud native ecosystem.

### Connector

At a higher level, Carbonaut integrates data over providers, and exposes collected data over a server (see [Server API docs](/docs/reference/server-api/)). Between these two components is the **Connector** component which contains the main lifecycle of the system.

![carbonaut building blocks](/concepts/building-blocks.drawio.png)

#### Internal Runtime

The connector, parses the configuration, starts and stopps plugins, updates the local state which contains the topology of the IT infrastructure and collects data from all connected providers which the API server exposes in different formats. A simplified version of the runtime is visualized below.

```mermaid
sequenceDiagram
    autonumber
    actor A as Alice the Platform Engineer
    actor J as John the DevOps Engineer
    participant CMain as Carbonaut CMD
    participant CServer as Carbonaut Server
    participant CConn as Carbonaut Connector
    participant EInfra as Infrastructure Data Sources
    participant EEnv as Environment Data Sources

    A->>CMain: Start Carbonaut
    activate CMain

    CMain->>CConn: Run (main lifecycle)
    activate CConn

    CMain->>CServer: Start Listening
    activate CServer

    A-->>CServer: Update Configuration file
    Note over CConn,EInfra: Carbonaut syncronizes the static resource state of <br> the infrastructure and maintains a mirrored state
    note right of CConn: Look up configuration file
    note right of CConn: Update State with Accounts

    par [Mirror Static Resource & Environment Data]
        loop
            activate CConn
            note right of CConn: Parse configured 'Accounts' <br> (new, old and remaining)
            CConn->>EInfra: Discover 'Projects' <br> (new, old and remaining)
            note right of CConn: Update State with Projects
            loop
                CConn->>EInfra: Discover 'Resources' by 'Project' <br> (new, old and remaining)
                note right of CConn: Update State with Resources
            end
            deactivate CConn
        end

    and [Serve Static & Dynamic Data]
        J->>CServer: I would like to get Carbonaut's metrics
        activate J
        activate CServer
        note right of CServer: Serve from Cache if set
        CServer->>CConn: Collect Static and Dynamic Data
        activate CConn
        loop Over all Accounts
            loop Over all Projects
                loop Over all Resources
                    CConn->>EInfra: Collect Dynamic Data <br> of the resource (dynres)
                    activate EInfra
                    EInfra-->>CConn: Energy Usage and CPU Frequency data
                    deactivate EInfra

                    CConn->>EEnv: Collect Dynamic Data (dynenv)
                    activate EEnv
                    EEnv-->>CConn: Energy Mix data
                    deactivate EEnv
                end
            end
        end
        note right of CConn: Collected dynamic data for each resource
        note right of CConn: Use static data from state
        CConn-->>CServer: Return Dynamic and Static Data
        deactivate CConn
        CServer-->>J: Visualize state data
        deactivate J
        note right of CServer: Update Cache
        deactivate CServer
    end

    deactivate CConn
    deactivate CServer
    CMain-->>A: Carbonaut Stopped
    deactivate CMain
```

#### Carbonaut Internal State

Carbonaut maintains an internal state which includes data which does not change until a resource was destroyed. Information about how much CPU cores or which Chip Architecture is considered static resource information. Information about the geolocation which indicate where the resource is hosted is considered static environment information. The data schema is defined [here](/docs/reference/schema/#type-staticresdata).

![carbonaut context](/concepts/state.drawio.png)

### Provider & Plugins

Carbonaut collects data of your infrastructure over data providers. Providers are interfaces which are implemented as plugins (see Plugin section in the sidebar). There are three different kinds of providers:

- `Dynamic Environment Provider` collects data about the Energy Mix (may be extended).
- `Dynamic Resource Provider` collects data about Energy Usage, CPU Frequency etc. (may be extended).
- `Static Resource Provider` collects data about CPU, Memory, IP etc. and also data about the geolocation and region of the IT resource.


```mermaid
classDiagram
    direction LR
    ResourceData <|-- EnergyMixData : energymix data depends on location \n data which is part of resource data
    ResourceData <|-- UtilizationData : utilization data depends \n on resource data
    class EnergyMixData{
      +SolarPercentage
      +CoalPercentage
      +...
    }
    class UtilizationData{
      +CPUFrequency
      +EnergyUsage
      +...
    }
    class ResourceData{
      +CPUCores
      +IP
      +Country
      +...
    }
```


![carbonaut data domain structure](/concepts/data-domain-structure.drawio.png)

These providers depend on each other. To collect energy usage of a system you first need to be aware of the system's topology. In the cloud environment we have heterogenous systems which changes dynamically. Therefore Resources are captured in projects (like K8s namespace's) and accounts (like K8s cluster's). The static resource provider just references the account and further discovery of projects and resources are done at runtime.

Providers expose interfaces which are defined in the schema reference document. The [`dynenv` provider](/docs/reference/schema/#type-provider) integrates [data](/docs/reference/schema/#type-dynamicenvdata) about the energy mix. The [`dynres` provider](/docs/reference/schema/#type-provider-1) integrates [data](/docs/reference/schema/#type-dynamicresdata) umong other things about energy usage. The [`staticres` provider](/docs/reference/schema/#type-provider-2) integrates [data](/docs/reference/schema/#type-staticresdata) umong other things about Operating System, Memory, Region.
