package mockcloudplugin

import (
	"carbonaut.dev/pkg/provider/plugin"
	"carbonaut.dev/pkg/provider/resource"
	"carbonaut.dev/pkg/provider/types/staticres"
)

var PluginName plugin.Kind = "mockcloudplugin"

type p struct {
	cfg *staticres.Config
}

func New(cfg *staticres.Config) (p, error) {
	return p{
		cfg: cfg,
	}, nil
}

func (p) GetName() *plugin.Kind {
	return &PluginName
}

func (p) DiscoverProjectIdentifiers() ([]*resource.ProjectName, error) {
	prjA := resource.ProjectName("project-a")
	prjB := resource.ProjectName("project-b")
	data := make([]*resource.ProjectName, 0)
	data = append(data, &prjA, &prjB)
	return data, nil
}

func (p) DiscoverStaticResourceIdentifiers(pName *resource.ProjectName) ([]*resource.ResourceName, error) {
	resA := resource.ResourceName("resource-a")
	resB := resource.ResourceName("resource-b")
	resC := resource.ResourceName("resource-c")
	data := make([]*resource.ResourceName, 0)
	data = append(data, &resA, &resB, &resC)
	return data, nil
}

func (p) GetStaticResourceData(pName *resource.ProjectName, rName *resource.ResourceName) (*resource.StaticResData, error) {
	return &resource.StaticResData{
		ID:   "0131acc3-82d8-488b-a8e2-c4a00e897145",
		User: "root",
		OS: &resource.OS{
			Version: "12",
			Distro:  "debian",
			Name:    "Debian 12",
		},
		IPv4: "145.40.93.88",
		CPUs: []*resource.CPU{
			{
				Count:        1,
				Type:         "Intel Xeon E-2278G 8-Core Processor @ 3.40GHz",
				Cores:        "8",
				Threads:      "16",
				Speed:        "3.40GHz",
				Arch:         "x86",
				Model:        "E-2278G",
				Manufacturer: "Intel",
				Name:         "Intel Xeon E-2278G Processor",
			},
		},
		GPUs: []*resource.GPU{
			{
				Count: 1,
				Type:  "Intel HD Graphics P630",
			},
		},
		NICs: []*resource.NIC{
			{
				Count: 1,
				Type:  "10Gbps",
			},
		},
		Drives: []*resource.DRIVE{
			{
				Count: 2,
				Type:  "SSD",
				Size:  "480GB",
			},
		},
		MemoryGB: "32GB",
		Location: &resource.Location{
			City:    "Frankfurt",
			Country: "DE",
			Address: "Kruppstrasse 121-127",
			ZipCode: "60388",
			Code:    "fr",
		},
	}, nil
}
