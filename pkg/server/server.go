package server

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"carbonaut.dev/pkg/connector"
	"carbonaut.dev/pkg/util/cache"
)

type Config struct {
	Port int `default:"8088" json:"port" yaml:"port"`
}

type Server struct {
	Connector *connector.C
	ExitChan  chan int
	cache     *cache.Cache
}

var (
	cacheDuration time.Duration = 30 * time.Second
	// the cache just contains one key value pair
	cacheCollectDataKey string = "collect-data"
	cacheStaticDataKey  string = "static-data"
)

func New(c *connector.C, exitChan chan int) *Server {
	// Create a cache with an expiration time of 60 seconds, and which
	// purges expired items every 5 minutes
	newCache := cache.New(cacheDuration, 5*time.Minute)
	return &Server{
		Connector: c,
		ExitChan:  exitChan,
		cache:     newCache,
	}
}

func (s Server) Listen(cfg *Config) {
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		ReadHeaderTimeout: 10 * time.Second,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(fmt.Sprintf("Carbonaut Server is running. Metrics can be collected on path http://0.0.0.0%s/metrics-json", server.Addr))); err != nil {
			slog.Error("could not write response", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/static-data", func(w http.ResponseWriter, r *http.Request) {
		data := s.Connector.GetStaticData()
		cachedStaticData, found := s.cache.Get(cacheStaticDataKey)
		if found {
			cachedData, ok := cachedStaticData.([]byte)
			if !ok {
				http.Error(w, "Internal Server Error, cached value is not of type []byte", http.StatusInternalServerError)
				return
			}
			slog.Debug("serving from cache")
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write(cachedData)
			if err != nil {
				slog.Error("could not write response", "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		slog.Info("loaded static data from state", "data", *data)
		jsonData, err := json.Marshal(data)
		if err != nil {
			slog.Error("could not serialize data to JSON", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		s.cache.Set(cacheStaticDataKey, jsonData, cache.DefaultExpiration)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(jsonData)
		if err != nil {
			slog.Error("could not write response", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/metrics-json", func(w http.ResponseWriter, r *http.Request) {
		cachedProjectResources, found := s.cache.Get(cacheCollectDataKey)
		if found {
			cachedProviderData, ok := cachedProjectResources.([]byte)
			if !ok {
				http.Error(w, "Internal Server Error, cached value is not of type []byte", http.StatusInternalServerError)
				return
			}
			slog.Debug("serving from cache")
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write(cachedProviderData)
			if err != nil {
				slog.Error("could not write response", "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		data, err := s.Connector.Collect()
		if err != nil {
			slog.Error("could not collect data from connector", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		slog.Info("collected new data", "data", *data)
		jsonData, err := json.Marshal(data)
		if err != nil {
			slog.Error("could not serialize data to JSON", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		s.cache.Set(cacheCollectDataKey, jsonData, cache.DefaultExpiration)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(jsonData)
		if err != nil {
			slog.Error("could not write response", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprintln(w, "Server is stopping"); err != nil {
			slog.Error("Error writing response", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		s.ExitChan <- 1
	})

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}
