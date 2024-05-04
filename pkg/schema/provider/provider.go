package provider

type T string

const (
	// Providers of this type discover the energy consumption of a host
	EnergyProvider T = "energy"
	// Providers of this type discover the it resources of a referenced it infrastructure
	ITResourceProvider T = "it_resource"
	// Providers of this type discover the geolocation of a host
	GeolocationProvider T = "geolocation"
	// Providers of this type discover the carbon intensity of a grid based on the geolocation of a host
	EmissionProvider T = "emission"
	// Providers of this type discover the energy mix of a grid based on the geolocation of a host
	EnergyMixProvider T = "energy_mix"
)
