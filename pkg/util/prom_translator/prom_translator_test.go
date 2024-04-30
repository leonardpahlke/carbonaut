// prom2json
// 	large parts of the logic copied from this project https://github.com/prometheus/prom2json
// 	with minor edits to turn it into a library and bake it into carbonaut

package prom_translator

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"testing"
// 	"time"

// 	"github.com/prometheus/client_golang/prometheus/promhttp"

// 	"go.opentelemetry.io/otel/attribute"
// 	"go.opentelemetry.io/otel/exporters/prometheus"
// 	api "go.opentelemetry.io/otel/metric"
// 	"go.opentelemetry.io/otel/sdk/metric"
// )

// func TestCollect(t *testing.T) {
// 	// port, err := freeport.GetFreePort()
// 	// if err != nil {
// 	// 	t.Errorf("library error, %v", err)
// 	// 	t.Fail()
// 	// }
// 	port := 2224
// 	startup_delay_seconds := 1
// 	t.Logf("Start prometheus test server on port %d", port)
// 	go test_server_up(port)

// 	t.Logf("Wait %d seconds for the prometheus server to start up", startup_delay_seconds)
// 	time.Sleep(time.Duration(startup_delay_seconds) * time.Second)

// 	prom2JsonConfig := Prom2Json{
// 		URL: fmt.Sprintf("http://0.0.0.0:%d", port),
// 	}
// 	t.Logf("Run CollectJSON test")
// 	b, err := CollectJSON(prom2JsonConfig)
// 	if err != nil {
// 		t.Errorf("could not collect metrics err: %v", err)
// 		t.Fail()
// 	}
// 	str := string(b)
// 	if str == "" {
// 		t.Errorf("received an empty string as response")
// 		t.Fail()
// 	}
// 	t.Logf("received string %s: ", str)

// 	d := []Family{}
// 	json.Unmarshal(b, &d)
// 	if len(d) == 0 {
// 		t.Errorf("received an empty list of metrics as response")
// 		t.Fail()
// 	}
// }

// func test_server_up(port int) {
// 	meterName := fmt.Sprintf("carbonaut-test-prom2json-%d", time.Now().Nanosecond())
// 	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
// 	ctx := context.Background()

// 	exporter, err := prometheus.New()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	provider := metric.NewMeterProvider(metric.WithReader(exporter))
// 	meter := provider.Meter(meterName)

// 	fmt.Println("Start the prometheus HTTP server and pass the exporter Collector to it")
// 	go serveMetrics(port)

// 	opt := api.WithAttributes(
// 		attribute.Key("A").String("B"),
// 		attribute.Key("C").String("D"),
// 	)

// 	// This is the equivalent of prometheus.NewCounterVec
// 	counter, err := meter.Float64Counter("foo", api.WithDescription("a simple counter"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	counter.Add(ctx, 5, opt)

// 	gauge, err := meter.Float64ObservableGauge("bar", api.WithDescription("a fun little gauge"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	_, err = meter.RegisterCallback(func(_ context.Context, o api.Observer) error {
// 		n := -10. + rng.Float64()*(90.) // [-10, 100)
// 		o.ObserveFloat64(gauge, n, opt)
// 		return nil
// 	}, gauge)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// This is the equivalent of prometheus.NewHistogramVec
// 	histogram, err := meter.Float64Histogram(
// 		"baz",
// 		api.WithDescription("a histogram with custom buckets and rename"),
// 		api.WithExplicitBucketBoundaries(64, 128, 256, 512, 1024, 2048, 4096),
// 	)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	histogram.Record(ctx, 136, opt)
// 	histogram.Record(ctx, 64, opt)
// 	histogram.Record(ctx, 701, opt)
// 	histogram.Record(ctx, 830, opt)

// 	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
// 	<-ctx.Done()
// }

// func serveMetrics(port int) {
// 	log.Printf("serving metrics at localhost:%d/metrics", port)
// 	http.Handle("/metrics", promhttp.Handler())
// 	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil) //nolint:gosec // Ignoring G114:
// 	// Use of net/http serve function that has no support for setting timeouts.
// 	if err != nil {
// 		fmt.Printf("error serving http: %v", err)
// 		return
// 	}
// }
