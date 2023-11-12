package schema

type Geolocation struct {
	Country string `json:"country"`
	// Country code (e.g. "US")
	CountryCode string `json:"countryCode"`
	// Region code (e.g. "CA")
	Region string `json:"region"`
	// Region name (e.g. "California")
	RegionName string `json:"regionName"`
	// City (e.g. "Mountain View")
	City string `json:"city"`
	// Zip code
	Zip string `json:"zip"`
	// Latitude
	Lat float64 `json:"lat"`
	// Longitude
	Lon float64 `json:"lon"`
	// IP address
	IP string `json:"ip"`
}

type Energy struct {
	// Amount of energy consumed
	Amount float64 `json:"amount"`
	// Unit of the amount of energy consumed
	Unit Unit `json:"unit"`
	// Name of the IT resource
	Name string `json:"name"`
}

type Emission struct {
	CarbonIntensity float64 `json:"carbonIntensity"`
}

type CarbonBreakdown struct {
	Zone            string `json:"zone"`
	CarbonIntensity int64  `json:"carbonIntensity"`
	Datetime        string `json:"datetime"`
}

type Unit string

const (
	NANOWATT  Unit = "nanowatt"
	MICROWATT Unit = "microwatt"
	MILLIWATT Unit = "milliwatt"
	KILOWATT  Unit = "kilowatt"
	MEGAWATT  Unit = "megawatt"
	GIGAWATT  Unit = "gigawatt"
	TERAWATT  Unit = "terawatt"

	// ZEPTOJOULE Unit = "zeptojoule"
	// NANOJOULE  Unit = "nanojoule"
	// MICROJOULE Unit = "microjoule"
	// KILOJOULE  Unit = "kilojoule"
	// MEGAJOULE  Unit = "megajoule"
	// GIGAJOULE  Unit = "gigajoule"
	// TERAJOULE  Unit = "terajoule"
)

func (e Energy) ConvertToKilowatt() float64 {
	// Get the conversion factor from the unit
	conversionFactor := 1.0
	switch e.Unit {
	case NANOWATT:
		conversionFactor = 1.0 / 1000000000000
	case MICROWATT:
		conversionFactor = 1.0 / 1000000000
	case MILLIWATT:
		conversionFactor = 1.0 / 1000000
	case KILOWATT:
		return e.Amount
	case MEGAWATT:
		conversionFactor = 1000
	case GIGAWATT:
		conversionFactor = 1000000
	case TERAWATT:
		conversionFactor = 1000000000
	default:
		return 0.0
	}

	// Convert the energy to kilowatts
	return e.Amount * conversionFactor
}

func (e Energy) ConvertToMilliwatt() float64 {
	conversionFactor := 1.0
	switch e.Unit {
	case NANOWATT:
		conversionFactor = 1.0 / 1000
	case MICROWATT:
		return e.Amount
	case MILLIWATT:
		conversionFactor = 1000
	case KILOWATT:
		conversionFactor = 1000000000
	case MEGAWATT:
		conversionFactor = 1000000000000
	case GIGAWATT:
		conversionFactor = 1000000000000000
	case TERAWATT:
		conversionFactor = 1000000000000000000
	default:
		return 0.0
	}

	// Convert the energy to kilowatts
	return e.Amount * conversionFactor
}

type ITResource struct {
	// The name of the IT resource
	Name string `json:"name"`
}
