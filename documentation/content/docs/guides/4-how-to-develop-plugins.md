---
title: Develop Plugins
weight: 4
slug: how-to-develop plugins
---


{{% steps %}}

### Setup

1. Setup your dev environment. Follow this [guide](/docs/guides/how-to-setup-dev-environment/).
2. Read up the [concepts documentation](/docs/components/#provider--plugins) about provider and plugins.

### Identify which provider interface you need to implement

1. Choose the provider which is used to implement the data you like to integrate.
   1. `dynres` implements energy usage, cpu frequency etc.
   2. `staticres` implements cloud provider, hypervisors etc. that manage your infrastructure resources.
   3. `dynenv` implements energy grid APIs
2. Setup your starting point (one of these)
   1. `dynres`: Copy the folder `pkg/plugin/dynresplugins/mockenergy` and rename it to your plugins name
   2. `staticres`: Copy the folder `pkg/plugin/staticresplugins/mockcloudplugin` and rename it to your plugins name
   3. `dynenv`: Copy the folder `pkg/plugin/dynenvplugins/mockenergymix` and rename it to your plugins name

### Implement your Plugin

1. Implement the provider interface. Add unit, integration and optional end to end testing.
   1. Make sure to add the provider reference to the plugin file:
      1. For `dynres` add it to `pkg/plugin/dynresplugins/dynresplugins.go`
      2. For `staticres` add it to `pkg/plugin/staticresplugins/staticresplugins.go`
      3. For `dynenv` add it to `pkg/plugin/dynenvplugins/dynenvplugins.go`
2. If applicable add caching.

```go
type p struct {
    // ...
	cache     *cache.Cache
}

func New(cfg *Config) (p, error) {
	// Create a cache with an expiration time of 60 seconds, and which
	// purges expired items every 5 minutes
	c := cache.New(60*time.Second, 5*time.Minute)

    // ...
}

func reqData(id string) (*Data, error) {
    // ...
    if cachedData, found := p.cache.Get(id); found {
        return cachedData, nil
    }

    // ...

    if err := p.cache.Add(id, data, cache.DefaultExpiration); err != nil {
		return nil, errors.New("unable to add data to internal cache")
	}

    // ...
}
```

### Push changes

Push changes and open a PR to propose merging changes to Carbonaut. See [contribution guide](/docs/contributing/).


{{% /steps %}}
