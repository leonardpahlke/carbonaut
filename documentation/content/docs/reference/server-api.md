---
title: Server API
weight: 1
slug: server-api
---

## **Carbonaut Server API**

<!-- TODO: automatically generate this information (like schema.md) -->

The Carbonaut Server is a lightweight HTTP server designed to provide metrics and static data management using a caching mechanism to optimize data delivery. This server is built using Go and utilizes a custom `connector` for data retrieval and a `cache` for data storage.

### Configuration

#### Struct: Config

- **Port** (`int`): The network port on which the server listens. Default is `8088`. The port can be specified in the `config.yaml`.

### HTTP Endpoints

{{% details title="`/`" closed="true" %}}
- **Method**: GET
- **Description**: Returns a message indicating that the Carbonaut Server is running and displays the metrics collection endpoint.
{{% /details %}}

{{% details title="`/static-data`" closed="true" %}}
- **Method**: GET
- **Description**: Provides static data retrieved via the connector. If available, data is served from the cache; otherwise, it is retrieved from the state and cached.
{{% /details %}}

{{% details title="`/metrics-json`" closed="true" %}}
- **Method**: GET
- **Description**: Returns collected data metrics. If available, data is served from the cache; otherwise, it is freshly collected, serialized to JSON, and cached.
{{% /details %}}

{{% details title="`/stop`" closed="true" %}}
- **Method**: GET
- **Description**: Stops the server. It sends a signal through the `ExitChan` to shut down the server.
{{% /details %}}

### Caching Strategy

The server employs a simple caching mechanism with a fixed expiration time for entries. It purges expired items every 5 minutes and caches both metrics data and static data for quick access.
