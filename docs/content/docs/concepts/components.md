---
weight: 2
---

## **Carbonaut Components**

Carbonaut has a couple of internal components which are explained in this document.

![carbonaut context](/docs/concepts/context.drawio.png)

This document focuses on the center component "Connector". Information about the User Interface around the Carbonaut Server is documented [here](/docs/concepts/server/). The integration of external data sources about providers is documented [here](/docs/concepts/data-providers/). 

### Connector: Internal Building Blocks

At a higher level, Carbonaut integrates data over providers (see [Data Provider docs](/docs/concepts/data-providers)), and exposes collected data over a server (see [Server API docs](/docs/reference/server-api/)). Between these two components is the **Connector** component which contains the main lifecycle of the system.

![carbonaut building blocks](/docs/concepts/building-blocks.drawio.png)



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

