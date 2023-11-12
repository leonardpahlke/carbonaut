# Carbonaut API schema

The Carbonaut schema is structured in three folders.

* **Plugin API definition** (`plugin/*.proto`): defines the gRPC API the Carbonaut Connector component and external plugins implement to transport data. The data transported are data transport objects (DTOs).
* **Carbonaut Schema**: (`schema/*.proto`): defines the integrated and analyzed data collected from connected plugins.
* **Generic Models** (`models/*.proto`): defines general models that are used in the `plugin` and `schema` definitions (e.g. model: Location).

Consider the diagram below to understand how the data schemas and api definitions are used in Carbonaut.
For a more comprehensive overview have a look to the [data schema](https://carbonaut.cloud/docs/concepts/schema) and [architecture](https://carbonaut.cloud/docs/concepts/architecture) documentation.

```mermaid
---
title: Carbonaut Data Schema Usage
---
%%{init: {'theme':'neutral'}}%%
flowchart LR
    models[/Generic Models/] ---|inherit| plugin[/Plugin API definition/] & schema[/Carbonaut Schema/]

    subgraph PluginSys[Carbonaut Plugin]
        pluginComponent[Carbonaut Plugin]
    end
    PluginSys --> exp[(External Data Provider)]
    User{{User}} ------> CarbonautSys
    subgraph CarbonautSys[Carbonaut]
        connector["Carbonaut Connector"] -->|gRPC| pluginComponent
        analyzer["Carbonaut Analyzer"] --> connector
        metrics["Carbonaut Metrics Server"] --> analyzer
    end

    plugin -...-|use| pluginComponent & connector & analyzer
    schema -.-|use| metrics & analyzer

    
    classDef Components fill:#9598C6,stroke:#24233d,stroke-width:4px
    classDef Definitions fill:#fdbe56,stroke:#24233d,stroke-width:4px
    class pluginComponent,connector,analyzer,metrics Components;
    class models,plugin,schema Definitions;
```

The `Plugin API definition` defines a gRPC service and with that a DTO Carbonaut data schema.
The `Carbonaut Plugin` implements the gRPC server and transforms the data scraped from the `External Data Provider` into schema.

The `Carbonaut Analyzer` translates the `Plugin API definition` schema to the `Carbonaut Schema` while enriching and mixing the data scraped.
The `Carbonaut Metrics Server` translates the `Carbonaut Schema` into OpenTelemetry metrics (Gauge, Histograms, etc.).

For a more comprehensive overview have a look to the [data schema](https://carbonaut.cloud/docs/concepts/schema) and [architecture](https://carbonaut.cloud/docs/concepts/architecture) documentation.
