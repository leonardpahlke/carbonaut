package main

import (
	"net/http"
	"os"
)

// This package serves as a fake Scaphandre server for testing purposes.

func main() {
	// read the local file containing the fake Scaphandre metrics
	// and serve them on the /metrics endpoint
	b, err := os.ReadFile("./scaph.txt")
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Write(b)
	})

	http.ListenAndServe(":8080", nil)
}
