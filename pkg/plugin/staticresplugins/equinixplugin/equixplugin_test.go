package equinixplugin_test

import (
	"os"
	"testing"

	"carbonaut.dev/pkg/plugin/staticresplugins/equinixplugin"
	"carbonaut.dev/pkg/provider/types/staticres"
)

func TestDiscoverProjectIdentifiers(t *testing.T) {
	authKey := "METAL_AUTH_TOKEN"
	if os.Getenv(authKey) == "" {
		t.Log("auth key not found as environment variable skip equinix test", authKey)
		t.SkipNow()
	}

	cfg := staticres.Config{
		Plugin:       &equinixplugin.PluginName,
		AccessKeyEnv: &authKey,
	}
	eP, err := equinixplugin.New(&cfg)
	if err != nil {
		t.Error("could not create equinixplugin", err)
	}
	_, err = eP.DiscoverProjectIdentifiers()
	if err != nil {
		t.Error("could not DiscoverProjectIdentifiers", err)
	}
}

func TestDiscoverGetResourceData(t *testing.T) {
	authKey := "METAL_AUTH_TOKEN"
	if os.Getenv(authKey) == "" {
		t.Log("auth key not found as environment variable skip equinix test", authKey)
		t.SkipNow()
	}

	cfg := staticres.Config{
		Plugin:       &equinixplugin.PluginName,
		AccessKeyEnv: &authKey,
	}
	eP, err := equinixplugin.New(&cfg)
	if err != nil {
		t.Error("could not create equinixplugin", err)
	}
	discoveredProjects, err := eP.DiscoverProjectIdentifiers()
	if err != nil {
		t.Error("could not DiscoverProjectIdentifiers", err)
	}

	for i := range discoveredProjects {
		_, err := eP.PGetResourceData(discoveredProjects[i])
		if err != nil {
			t.Error("could not PGetResourceData", err)
		}
		// TODO: improve check of cache
		_, err = eP.PGetResourceData(discoveredProjects[i])
		if err != nil {
			t.Error("could not PGetResourceData", err)
		}
	}
}
