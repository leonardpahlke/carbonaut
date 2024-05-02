package data

// computer hardware data
type StaticResourceData struct {
	IP         string
	CPUCores   int
	MEMORY_MB  int
	Arch       string
	Storage_GB int
}

// energy and utilization data
type DynamicResourceData struct {
	CPU_Frequency         float32
	Energy_Host_Milliwatt int
}

// location data
type StaticEnvironmentData struct {
	Region  string
	Country string
}

// energy mix data
type DynamicEnvironmentData struct{}
