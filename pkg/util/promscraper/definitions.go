// prom2json
// 	large parts of the logic copied from this project https://github.com/prometheus/prom2json
// 	with minor edits to turn it into a library and bake it into carbonaut

package promscraper

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"strconv"

	"github.com/prometheus/common/expfmt"
	"google.golang.org/protobuf/encoding/protodelim"
	"google.golang.org/protobuf/proto"

	dto "github.com/prometheus/client_model/go"
)

const acceptHeader = `application/vnd.google.protobuf;proto=io.prometheus.client.MetricFamily;encoding=delimited;q=0.7,text/plain;version=0.0.4;q=0.3`

// Family mirrors the MetricFamily proto message.
type Family struct {
	Name    string        `json:"name"`
	Help    string        `json:"help"`
	Type    string        `json:"type"`
	Metrics []interface{} `json:"metrics,omitempty"` // Either metric or summary.
}

// Metric is for all "single value" metrics, i.e. Counter, Gauge, and Untyped.
type Metric struct {
	Labels      map[string]string `json:"labels,omitempty"`
	TimestampMs string            `json:"timestamp_ms,omitempty"`
	Value       string            `json:"value"`
}

// Summary mirrors the Summary proto message.
type Summary struct {
	Labels      map[string]string `json:"labels,omitempty"`
	TimestampMs string            `json:"timestamp_ms,omitempty"`
	Quantiles   map[string]string `json:"quantiles,omitempty"`
	Count       string            `json:"count"`
	Sum         string            `json:"sum"`
}

// Histogram mirrors the Histogram proto message.
type Histogram struct {
	Labels      map[string]string `json:"labels,omitempty"`
	TimestampMs string            `json:"timestamp_ms,omitempty"`
	Buckets     map[string]string `json:"buckets,omitempty"`
	Count       string            `json:"count"`
	Sum         string            `json:"sum"`
}

// newFamily consumes a MetricFamily and transforms it to the local Family type.
func newFamily(dtoMF *dto.MetricFamily) *Family {
	mf := &Family{
		Name:    dtoMF.GetName(),
		Help:    dtoMF.GetHelp(),
		Type:    dtoMF.GetType().String(),
		Metrics: make([]interface{}, len(dtoMF.GetMetric())),
	}
	for i, m := range dtoMF.GetMetric() {
		switch dtoMF.GetType() {
		case dto.MetricType_SUMMARY:
			mf.Metrics[i] = Summary{
				Labels:      makeLabels(m),
				TimestampMs: makeTimestamp(m),
				Quantiles:   makeQuantiles(m),
				Count:       strconv.FormatUint(m.GetSummary().GetSampleCount(), 10),
				Sum:         strconv.FormatFloat(m.GetSummary().GetSampleSum(), 'f', -1, 64),
			}
		case dto.MetricType_HISTOGRAM:
			mf.Metrics[i] = Histogram{
				Labels:      makeLabels(m),
				TimestampMs: makeTimestamp(m),
				Buckets:     makeBuckets(m),
				Count:       strconv.FormatUint(m.GetHistogram().GetSampleCount(), 10),
				Sum:         strconv.FormatFloat(m.GetHistogram().GetSampleSum(), 'f', -1, 64),
			}
		default:
			value := getValue(m)
			mf.Metrics[i] = Metric{
				Labels:      makeLabels(m),
				TimestampMs: makeTimestamp(m),
				Value:       strconv.FormatFloat(value, 'f', -1, 64),
			}
		}
	}
	return mf
}

func getValue(m *dto.Metric) float64 {
	switch {
	case m.GetGauge() != nil:
		return m.GetGauge().GetValue()
	case m.GetCounter() != nil:
		return m.GetCounter().GetValue()
	case m.GetUntyped() != nil:
		return m.GetUntyped().GetValue()
	default:
		return 0.
	}
}

func makeLabels(m *dto.Metric) map[string]string {
	result := map[string]string{}
	for _, lp := range m.GetLabel() {
		result[lp.GetName()] = lp.GetValue()
	}
	return result
}

func makeTimestamp(m *dto.Metric) string {
	if m.TimestampMs == nil {
		return ""
	}
	return strconv.FormatInt(m.GetTimestampMs(), 10)
}

func makeQuantiles(m *dto.Metric) map[string]string {
	result := map[string]string{}
	for _, q := range m.GetSummary().GetQuantile() {
		quantile := strconv.FormatFloat(q.GetQuantile(), 'f', -1, 64)
		value := strconv.FormatFloat(q.GetValue(), 'f', -1, 64)
		result[quantile] = value
	}
	return result
}

func makeBuckets(m *dto.Metric) map[string]string {
	result := map[string]string{}
	for _, b := range m.GetHistogram().GetBucket() {
		upperBound := strconv.FormatFloat(b.GetUpperBound(), 'f', -1, 64) // Assumes GetUpperBound returns a float64
		cumulativeCount := strconv.FormatUint(b.GetCumulativeCount(), 10) // Assumes GetCumulativeCount returns a uint64
		result[upperBound] = cumulativeCount
	}

	return result
}

// fetchMetricFamilies retrieves metrics from the provided URL, decodes them
// into MetricFamily proto messages, and sends them to the provided channel. It
// returns after all MetricFamilies have been sent. The provided transport
// may be nil (in which case the default Transport is used).
func fetchMetricFamilies(url string, ch chan<- *dto.MetricFamily, transport http.RoundTripper) error {
	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		close(ch)
		return fmt.Errorf("creating %s request for URL %q failed: %v", http.MethodGet, url, err)
	}
	req.Header.Add("Accept", acceptHeader)
	client := http.Client{Transport: transport}
	resp, err := client.Do(req)
	if err != nil {
		close(ch)
		return fmt.Errorf("executing %s request for URL %q failed: %v", http.MethodGet, url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		close(ch)
		return fmt.Errorf("%s request for URL %q returned HTTP status %s", http.MethodGet, url, resp.Status)
	}
	return parseResponse(resp, ch)
}

// parseResponse consumes an http.Response and pushes it to the MetricFamily
// channel. It returns when all MetricFamilies are parsed and put on the
// channel.
func parseResponse(resp *http.Response, ch chan<- *dto.MetricFamily) error {
	mediatype, params, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err == nil && mediatype == "application/vnd.google.protobuf" &&
		params["encoding"] == "delimited" &&
		params["proto"] == "io.prometheus.client.MetricFamily" {
		defer close(ch)
		for {
			mf := &dto.MetricFamily{}
			if _, err = ReadDelimited(resp.Body, mf); err != nil {
				if err == io.EOF {
					break
				}
				return fmt.Errorf("reading metric family protocol buffer failed: %v", err)
			}
			ch <- mf
		}
	} else {
		if err := parseReader(resp.Body, ch); err != nil {
			return err
		}
	}
	return nil
}

// parseReader consumes an io.Reader and pushes it to the MetricFamily
// channel. It returns when all MetricFamilies are parsed and put on the
// channel.
func parseReader(in io.Reader, ch chan<- *dto.MetricFamily) error {
	defer close(ch)
	// We could do further content-type checks here, but the
	// fallback for now will anyway be the text format
	// version 0.0.4, so just go for it and see if it works.
	var parser expfmt.TextParser
	metricFamilies, err := parser.TextToMetricFamilies(in)
	if err != nil {
		return fmt.Errorf("reading text format failed: %v", err)
	}
	for _, mf := range metricFamilies {
		ch <- mf
	}
	return nil
}

// AddLabel allows to add key/value labels to an already existing Family.
func (f *Family) AddLabel(key, val string) {
	for i, item := range f.Metrics {
		if m, ok := item.(Metric); ok {
			m.Labels[key] = val
			f.Metrics[i] = m
		}
	}
}

// ReadDelimited decodes a message from the provided length-delimited stream,
// where the length is encoded as 32-bit varint prefix to the message body.
// It returns the total number of bytes read and any applicable error.  This is
// roughly equivalent to the companion Java API's
// MessageLite#parseDelimitedFrom.  As per the reader contract, this function
// calls r.Read repeatedly as required until exactly one message including its
// prefix is read and decoded (or an error has occurred).  The function never
// reads more bytes from the stream than required.  The function never returns
// an error if a message has been read and decoded correctly, even if the end
// of the stream has been reached in doing so.  In that case, any subsequent
// calls return (0, io.EOF).
func ReadDelimited(r io.Reader, m proto.Message) (n int, err error) {
	cr := &countingReader{r: r}
	opts := protodelim.UnmarshalOptions{
		MaxSize: -1,
	}

	err = opts.UnmarshalFrom(cr, m)
	return cr.n, err
}

type countingReader struct {
	r io.Reader
	n int
}

// implements protodelim.Reader
func (r *countingReader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	if n > 0 {
		r.n += n
	}
	return n, err
}

// implements protodelim.Reader
func (r *countingReader) ReadByte() (byte, error) {
	var buf [1]byte
	for {
		n, err := r.Read(buf[:])
		if n == 0 && err == nil {
			// io.Reader states: Callers should treat a return of 0 and nil as
			// indicating that nothing happened; in particular it does not
			// indicate EOF.
			continue
		}
		return buf[0], err
	}
}
