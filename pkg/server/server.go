package server

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"carbonaut.dev/pkg/connector"
)

type Config struct {
	Port int `json:"port" yaml:"port" default:"8088"`
}

type Server struct {
	Connector *connector.C
	Log       *slog.Logger
	ExitChan  chan int
}

func New(Connector *connector.C, Log *slog.Logger, ExitChan chan int) *Server {
	return &Server{
		Connector: Connector,
		Log:       Log,
		ExitChan:  ExitChan,
	}
}

func (s Server) Listen(cfg *Config) {
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Port),
	}

	http.HandleFunc("/metrics-json", func(w http.ResponseWriter, r *http.Request) {
		data, err := s.Connector.Collect()
		if err != nil {
			s.Log.Error("could not collect data from connector", "component", "server", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		s.Log.Info("collected data", "component", "server", "data", *data)

		jsonData, err := json.Marshal(data)
		if err != nil {
			s.Log.Error("could not serialize data to JSON", "component", "server", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Server is stopping")
		s.ExitChan <- 1
	})

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}
