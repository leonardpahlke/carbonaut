package resource

import "carbonaut.dev/pkg/provider/environment"

type (
	AccountName  string
	ProjectName  string
	ResourceName string

	AccountData map[ProjectName]ProjectData
	ProjectData map[ResourceName]*ResourceData
)

type ResourceData struct {
	DynamicData *DynamicData   `json:"dynamic_data" yaml:"dynamic_data"`
	StaticData  *StaticResData `json:"static_data"  yaml:"static_data"`
}

type DynamicData struct {
	ResData *DynamicResData             `json:"res_data" yaml:"res_data"`
	EnvData *environment.DynamicEnvData `json:"env_data" yaml:"env_data"`
}

type StaticData struct {
	ResData *StaticResData `json:"res_data" yaml:"res_data"`
}

// energy and utilization data
type DynamicResData struct {
	CPUFrequency          float64 `json:"cpu_frequency"          yaml:"cpu_frequency"`
	EnergyHostMicrojoules int     `json:"energy_host_mirojoules" yaml:"energy_host_mirojoules"`
	CPULoadPercentage     float64 `json:"cpu_load_percentage"    yaml:"cpu_load_percentage"`
}

// Data represents computer hardware data.
type StaticResData struct {
	ID       string    `json:"id"        yaml:"id"        default:"0131acc3-82d8-488b-a8e2-c4a00e897145"`
	User     string    `json:"user"      yaml:"user"      default:"root"`
	OS       *OS       `json:"os"        yaml:"os"`
	IPv4     string    `json:"ipv4"      yaml:"ipv4"      default:"145.40.93.80"`
	CPUs     []*CPU    `json:"cpus"      yaml:"cpus"`
	GPUs     []*GPU    `json:"gpus"      yaml:"gpus"`
	NICs     []*NIC    `json:"nics"      yaml:"nics"`
	Drives   []*DRIVE  `json:"drives"    yaml:"drives"`
	MemoryGB string    `json:"memory_gb" yaml:"memory_gb" default:"32GB"`
	Location *Location `json:"location"  yaml:"location"`
}

type CPU struct {
	Count        int    `json:"count"        yaml:"count"        default:"1"`
	Type         string `json:"type"         yaml:"type"         default:"Intel Xeon E-2278G 8-Core Processor @ 3.40GHz"`
	Cores        string `json:"cores"        yaml:"cores"        default:"8"`
	Threads      string `json:"threads"      yaml:"threads"      default:"16"`
	Speed        string `json:"speed"        yaml:"speed"        default:"3.40GHz"`
	Arch         string `json:"arch"         yaml:"arch"         default:"x86"`
	Model        string `json:"model"        yaml:"model"        default:"E-2278G"`
	Manufacturer string `json:"manufacturer" yaml:"manufacturer" default:"Intel"`
	Name         string `json:"name"         yaml:"name"         default:"Intel Xeon E-2278G Processor"`
}

type GPU struct {
	Count int    `json:"count" yaml:"count" default:"1"`
	Type  string `json:"type"  yaml:"type"  default:"Intel HD Graphics P630"`
}

type NIC struct {
	Count int    `json:"count" yaml:"count" default:"1"`
	Type  string `json:"type"  yaml:"type"  default:"10Gbps"`
}

type DRIVE struct {
	Count int    `json:"count" yaml:"count" default:"2"`
	Type  string `json:"type"  yaml:"type"  default:"SSD"`
	Size  string `json:"size"  yaml:"type"  default:"480GB"`
}

type OS struct {
	Version string `json:"version" yaml:"version" default:"12"`
	Distro  string `json:"distro"  yaml:"distro"  default:"debian"`
	Name    string `json:"name"    yaml:"name"    default:"Debian 12"`
}

type Location struct {
	City    string `json:"city"     yaml:"city"     default:"Frankfurt"`
	Country string `json:"country"  yaml:"country"  default:"DE"`
	Address string `json:"address"  yaml:"address"  default:"Kruppstrasse 121-127"`
	ZipCode string `json:"zip_code" yaml:"zip_code" default:"60388"`
	Code    string `json:"code"     yaml:"code"     default:"fr"`
}
