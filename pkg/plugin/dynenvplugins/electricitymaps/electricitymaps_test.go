package electricitymaps_test

import (
	"log/slog"
	"os"
	"testing"

	"carbonaut.dev/pkg/plugin/dynenvplugins/electricitymaps"
	"carbonaut.dev/pkg/provider/account/project/resource"
	"carbonaut.dev/pkg/provider/types/dynenv"
	"github.com/creasty/defaults"
)

var (
	envVarKeyMock = "ELECTRICITY_MAP_ENV_KEY_TESTING"
	envVarKey     = "ELECTRICITY_MAP_AUTH_TOKEN"
	cfg           = dynenv.Config{
		Plugin:       &electricitymaps.PluginName,
		AccessKeyEnv: &envVarKeyMock,
	}
)

func TestQueryZones(t *testing.T) {
	os.Setenv(envVarKeyMock, "some value")
	p, err := electricitymaps.New(&cfg)
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 2; i++ {
		z, err := p.QueryZones()
		if err != nil {
			t.Error(err)
		}
		if len(*z) == 0 {
			t.Error("no zones parsed", z)
		}
	}
}

func TestGetDynamicEnvironmentData(t *testing.T) {
	if os.Getenv(envVarKey) == "" {
		t.Logf("skip test since %s is not set. This test requires an API key", envVarKey)
		t.SkipNow()
	}
	os.Setenv(envVarKeyMock, os.Getenv(envVarKey))
	p, err := electricitymaps.New(&cfg)
	if err != nil {
		t.Error(err)
	}

	locationData := resource.Location{}
	if err := defaults.Set(&locationData); err != nil {
		t.Error("internal testing error, failed to set location data defaults", "error", err)
	}
	d, err := p.GetDynamicEnvironmentData(&locationData)
	if err != nil {
		t.Error(err)
	}
	if d.RenewablePercentage+d.FossilFuelsPercentage == 0 {
		t.Error("RenewablePercentage and FossilFuelsPercentage is summed up 0 which kinda impossible")
	}
	slog.Info("GetDynamicEnvironmentData results", "energy breakdown", d)
}
