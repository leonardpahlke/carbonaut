package electricitymaps

import (
	"time"

	"carbonaut.dev/pkg/provider/environment"
)

type (
	Zones map[Zone]ZoneDetails
	Zone  string
)

type ZoneDetails struct {
	CountryName string `json:"countryName"`
	ZoneName    string `json:"zoneName"`
	DisplayName string `json:"displayName"`
}

func (data *ZoneEnergyBreakdown) Convert() *environment.DynamicEnvData {
	total := data.PowerConsumptionBreakdown.Nuclear + data.PowerConsumptionBreakdown.Geothermal + data.PowerConsumptionBreakdown.Biomass + data.PowerConsumptionBreakdown.Coal + data.PowerConsumptionBreakdown.Wind + data.PowerConsumptionBreakdown.Solar +
		data.PowerConsumptionBreakdown.Hydro + data.PowerConsumptionBreakdown.Gas + data.PowerConsumptionBreakdown.Oil + data.PowerConsumptionBreakdown.Unknown + data.PowerConsumptionBreakdown.HydroDischarge + data.PowerConsumptionBreakdown.BatteryDischarge

	if total == 0 {
		return &environment.DynamicEnvData{}
	}

	// Calculate percentages
	p := environment.DynamicEnvData{
		SolarPercentage:        float64(data.PowerConsumptionBreakdown.Solar) / float64(total) * 100,
		WindPercentage:         float64(data.PowerConsumptionBreakdown.Wind) / float64(total) * 100,
		HydroPercentage:        float64(data.PowerConsumptionBreakdown.Hydro) / float64(total) * 100,
		GeothermalPercentage:   float64(data.PowerConsumptionBreakdown.Geothermal) / float64(total) * 100,
		BiomassPercentage:      float64(data.PowerConsumptionBreakdown.Biomass) / float64(total) * 100,
		NuclearPercentage:      float64(data.PowerConsumptionBreakdown.Nuclear) / float64(total) * 100,
		GasPercentage:          float64(data.PowerConsumptionBreakdown.Gas) / float64(total) * 100,
		OilPercentage:          float64(data.PowerConsumptionBreakdown.Oil) / float64(total) * 100,
		CoalPercentage:         float64(data.PowerConsumptionBreakdown.Coal) / float64(total) * 100,
		OtherSourcesPercentage: float64(data.PowerConsumptionBreakdown.Unknown+data.PowerConsumptionBreakdown.BatteryDischarge+data.PowerConsumptionBreakdown.HydroDischarge) / float64(total) * 100,
	}

	// Calculate summary categories
	fossilFuels := data.PowerConsumptionBreakdown.Coal + data.PowerConsumptionBreakdown.Gas + data.PowerConsumptionBreakdown.Oil + data.PowerConsumptionBreakdown.Nuclear
	renewable := data.PowerConsumptionBreakdown.Wind + data.PowerConsumptionBreakdown.Solar + data.PowerConsumptionBreakdown.Hydro + data.PowerConsumptionBreakdown.Geothermal + data.PowerConsumptionBreakdown.Biomass

	p.FossilFuelsPercentage = float64(fossilFuels) / float64(total) * 100
	p.RenewablePercentage = float64(renewable) / float64(total) * 100

	return &p
}

// numbers are in MW - https://static.electricitymaps.com/api/docs/index.html#live-power-breakdown
type ZoneEnergyBreakdown struct {
	Zone                      string    `json:"zone"`
	Datetime                  time.Time `json:"datetime"`
	UpdatedAt                 time.Time `json:"updatedAt"`
	CreatedAt                 time.Time `json:"createdAt"`
	PowerConsumptionBreakdown struct {
		Nuclear          int `json:"nuclear"`
		Geothermal       int `json:"geothermal"`
		Biomass          int `json:"biomass"`
		Coal             int `json:"coal"`
		Wind             int `json:"wind"`
		Solar            int `json:"solar"`
		Hydro            int `json:"hydro"`
		Gas              int `json:"gas"`
		Oil              int `json:"oil"`
		Unknown          int `json:"unknown"`
		HydroDischarge   int `json:"hydro discharge"`
		BatteryDischarge int `json:"battery discharge"`
	} `json:"powerConsumptionBreakdown"`
	PowerProductionBreakdown struct {
		Nuclear          int         `json:"nuclear"`
		Geothermal       int         `json:"geothermal"`
		Biomass          int         `json:"biomass"`
		Coal             int         `json:"coal"`
		Wind             int         `json:"wind"`
		Solar            int         `json:"solar"`
		Hydro            int         `json:"hydro"`
		Gas              int         `json:"gas"`
		Oil              int         `json:"oil"`
		Unknown          interface{} `json:"unknown"`
		HydroDischarge   int         `json:"hydro discharge"`
		BatteryDischarge int         `json:"battery discharge"`
	} `json:"powerProductionBreakdown"`
	// PowerImportBreakdown struct {
	// 	BE int `json:"BE"`
	// 	ES int `json:"ES"`
	// 	GB int `json:"GB"`
	// } `json:"powerImportBreakdown"`
	// PowerExportBreakdown struct {
	// 	BE int `json:"BE"`
	// 	ES int `json:"ES"`
	// 	GB int `json:"GB"`
	// } `json:"powerExportBreakdown"`
	FossilFreePercentage  int    `json:"fossilFreePercentage"`
	RenewablePercentage   int    `json:"renewablePercentage"`
	PowerConsumptionTotal int    `json:"powerConsumptionTotal"`
	PowerProductionTotal  int    `json:"powerProductionTotal"`
	PowerImportTotal      int    `json:"powerImportTotal"`
	PowerExportTotal      int    `json:"powerExportTotal"`
	IsEstimated           bool   `json:"isEstimated"`
	EstimationMethod      string `json:"estimationMethod"`
}
