package server

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"carbonaut.dev/pkg/connector"
)

type Config struct {
	Port int `default:"8088" json:"port" yaml:"port"`
}

type Server struct {
	Connector *connector.C
	Log       *slog.Logger
	ExitChan  chan int
}

type CacheItem struct {
	Data      []byte
	Timestamp time.Time
}

var (
	cache         *CacheItem    = nil
	cacheDuration time.Duration = 10 * time.Second
)

func New(c *connector.C, logger *slog.Logger, exitChan chan int) *Server {
	return &Server{
		Connector: c,
		Log:       logger,
		ExitChan:  exitChan,
	}
}

func (s Server) Listen(cfg *Config) {
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		ReadHeaderTimeout: 10 * time.Second,
	}

	http.HandleFunc("/metrics-json", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		if cache != nil && now.Sub(cache.Timestamp) < cacheDuration {
			s.Log.Info("serving from cache")
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write(cache.Data)
			if err != nil {
				s.Log.Error("could not write response", "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		data, err := s.Connector.Collect()
		if err != nil {
			s.Log.Error("could not collect data from connector", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		s.Log.Info("collected new data", "data", *data)
		jsonData, err := json.Marshal(data)
		if err != nil {
			s.Log.Error("could not serialize data to JSON", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		cache = &CacheItem{
			Data:      jsonData,
			Timestamp: now,
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(jsonData)
		if err != nil {
			s.Log.Error("could not write response", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprintln(w, "Server is stopping"); err != nil {
			s.Log.Error("Error writing response", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		s.ExitChan <- 1
	})

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}
