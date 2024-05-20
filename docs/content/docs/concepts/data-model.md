---
weight: 4
bookHidden: true
---

## **Carbonaut Data Model**

![carbonaut data domain structure](/docs/concepts/data-domain-structure.drawio.png)

{{< mermaid class="optional" >}}
classDiagram
    LocationData <|-- EnergyMixData : energymix data depends \n on location data
    ResourceData <|-- UtilizationData : utilization data depends \n on resource data
    ResourceData <|-- LocationData : location data depends \n on resource data 
    
    class LocationData{
      +Region
      +Country
      +...
    }
    class EnergyMixData{
      +SolarPercentage
      +CoalPercentage
      +HydroPercentage
      +...
    }
    class UtilizationData{
      -CPUFrequency
      -EnergyUsage
      -...
    }
    class ResourceData{
      +CPUCores
      +Memory
      +IP
      +...
    }

    style LocationData fill:#D5E8D4
    style EnergyMixData fill:#D5E8D4
    style ResourceData fill:#DAE8FC
    style UtilizationData fill:#DAE8FC
{{< /mermaid >}}