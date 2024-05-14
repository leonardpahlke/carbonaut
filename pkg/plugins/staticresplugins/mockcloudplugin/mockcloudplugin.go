package mockcloudplugin

import (
	"carbonaut.dev/pkg/schema/plugin"
	"carbonaut.dev/pkg/schema/provider/data/account/project"
	"carbonaut.dev/pkg/schema/provider/data/account/project/resource"
	"carbonaut.dev/pkg/schema/provider/types/staticres"
)

var PluginName plugin.Kind = "mockcloudplugin"

type p struct{}

func New() p {
	return p{}
}

func (p) GetName() *plugin.Kind {
	return &PluginName
}

func (p) DiscoverProjectIdentifiers(cfg *staticres.Config) ([]*project.Name, error) {
	prjA := project.Name("project-a")
	prjB := project.Name("project-b")
	data := make([]*project.Name, 0)
	data = append(data, &prjA, &prjB)
	return data, nil
}

func (p) DiscoverStaticResourceIdentifiers(cfg *staticres.Config, pName *project.Name) ([]*resource.Name, error) {
	resA := resource.Name("resource-a")
	resB := resource.Name("resource-b")
	resC := resource.Name("resource-c")
	data := make([]*resource.Name, 0)
	data = append(data, &resA, &resB, &resC)
	return data, nil
}

func (p) GetStaticResourceData(cfg *staticres.Config, pName *project.Name, rName *resource.Name) (*resource.StaticResData, error) {
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
				Cores:        8,
				Threads:      16,
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
