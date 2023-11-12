package provierutils

type Provider string

const (
	// Providers of this type discover the energy consumption of a host
	EnergyProvider Provider = "energy"
	// Providers of this type discover the it resources of a referenced it infrastructure
	ITResourceProvider Provider = "it_resource"
	// Providers of this type discover the geolocation of a host
	GeolocationProvider Provider = "geolocation"
	// Providers of this type discover the carbon intensity of a grid based on the geolocation of a host
	EmissionProvider Provider = "emission"
	// Providers of this type discover the energy mix of a grid based on the geolocation of a host
	EnergyMixProvider Provider = "energy_mix"
)
