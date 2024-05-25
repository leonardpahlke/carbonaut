---
title: Documentation
next: first-page
---


![carbonaut banner image](carbonaut-banner.png)

Carbonaut is project to collect and refine environmental sustainability data from your IT infrastructure and make it available in a common data schema. The project does not implement a scraper to scan your IT infrastructure or measure your virtual machine energy, but integrates over providers which are implemented as plugins to integrate external data sources. The project is a POC and therefore not production ready. If you like to [contribute](/docs/reference/contributing) you are very welcome to do so!

{{< callout type="warning" >}}
  Carbonaut is a proof of concept project which may be turned in the future into a proper open source project. Right now this is not the case.
{{< /callout >}}

{{% details title="Click me to reveal" closed="true" %}}

The aim of the project is to collect and integrate data around environmental sustainability and make it accessible to users so that they can make automatic or manual decisions. The project integrates heterogenous data sources over plugins and transforms the data in a simple minimalistic data schema. Carbonaut discovers created IT infrastructure, maintains a minimalistic internal static resource state and chains dynamic changing data together to collect a holistic overview of cloud native sustainability data.


{{% /details %}}


![carbonaut context](/docs/concepts/context.drawio.png)


Carbonaut is used with other tools and is due to the plugin design platform agnostic. It's developed for a cloud native usecase. There are several other tools available which are also about cloud native sustainability.

- **[Kepler](https://github.com/sustainable-computing-io/kepler)** measures and exposes energy and carbon emission metrics of Kubernetes clusters. Unlike Carbonaut, Kepler employs its measurement tools, which include RAPL, eBPF approximations, and others. It actively discovers new nodes and pods, updating its measurements accordingly. On the other hand, Carbonaut does not directly measure metrics but integrates with projects that deliver these metrics. Carbonaut is not restricted to Kubernetes, offering a more agnostic approach. Carbonaut is capable of integrating IT infrastructure metrics across platforms like GCP, Equinix, and AWS. It can utilize RAPL measurements on one platform while employing different providers on another, and use EnergyMap to obtain energy mix data for all. Note that not all of these plugins have been implemented yet. If the system is entirely based on Kubernetes, Kepler is likely a better choice (especially at the current stage of development of the Carbonaut project).
- **[Scaphandre](https://github.com/hubblo-org/scaphandre)** uses various interfaces to collect energy and resource usage data from a computer. Scaphandre is a lightweight binary that offers various ways to export collected data for further use. Carbonaut implements a Scaphandre provider to collect dynamic resource data (energy consumption, CPU frequency, etc.). Carbonaut provides the interface to the cloud native sustainability information, scaphandre provides the interface to the environmental data of one machine.
- **[Kube Green](https://github.com/kube-green/kube-green)** is a Kubernetes operator that looks at pod annotations to schedule and de-schedule them at specific times to improve energy utilization in terms of power mix. Kube Green does not measure or collect data, but uses a static approach to change how the cluster manages its resources. Carbonaut collects and integrates data, but does not change the system topology any further. In the future, the idea of Kube Green could be further developed with the data from Carbonaut.

{{< callout type="warning" >}}
  The current Carbonaut version is a POC (Proof of Concept). The project is not tested for large scale environments. There are several security and quality checks already in place but no security assessment was made and there are likely vulnerabilities if deployed. The system is not hardened. The data scheme is minimalistic which aligns with the nature of a POC but lacks therefore complexity which is required to cover all aspects of cloud native sustainability.
{{< /callout >}}

