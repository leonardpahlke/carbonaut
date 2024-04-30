// prom2json
// 	large parts of the logic copied from this project https://github.com/prometheus/prom2json
// 	with minor edits to turn it into a library and bake it into carbonaut

package prom_translator

import (
	"math"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	dto "github.com/prometheus/client_model/go"
)

type testCase struct {
	name    string
	mFamily *dto.MetricFamily
	output  *Family
}

var tcs = []testCase{
	{
		name: "test counter",
		mFamily: &dto.MetricFamily{
			Name:   strPtr("counter1"),
			Help:   new(string),
			Type:   metricTypePtr(dto.MetricType_COUNTER),
			Metric: []*dto.Metric{{Label: []*dto.LabelPair{createLabelPair("tag1", "abc"), createLabelPair("tag2", "def")}, Counter: &dto.Counter{Value: floatPtr(1)}}, {Label: []*dto.LabelPair{createLabelPair("tag1", "foo"), createLabelPair("tag2", "bar")}, TimestampMs: intPtr(123456), Counter: &dto.Counter{Value: floatPtr(42)}}, {Label: []*dto.LabelPair{}, Counter: &dto.Counter{Value: floatPtr(2)}}, {Label: []*dto.LabelPair{createLabelPair("inf", "neg")}, Counter: &dto.Counter{Value: floatPtr(math.Inf(-1))}}, {Label: []*dto.LabelPair{createLabelPair("inf", "pos")}, Counter: &dto.Counter{Value: floatPtr(math.Inf(1))}}},
			Unit:   new(string),
		},
		output: &Family{
			Name: "counter1",
			Help: "",
			Type: "COUNTER",
			Metrics: []interface{}{
				Metric{
					Labels: map[string]string{
						"tag2": "def",
						"tag1": "abc",
					},
					Value: "1",
				},
				Metric{
					Labels: map[string]string{
						"tag2": "bar",
						"tag1": "foo",
					},
					TimestampMs: "123456",
					Value:       "42",
				},
				Metric{
					Labels: map[string]string{},
					Value:  "2",
				},
				Metric{
					Labels: map[string]string{
						"inf": "neg",
					},
					Value: "-Inf",
				},
				Metric{
					Labels: map[string]string{
						"inf": "pos",
					},
					Value: "+Inf",
				},
			},
		},
	},
	{
		name: "test summaries",
		mFamily: &dto.MetricFamily{
			Name: strPtr("summary1"),
			Type: metricTypePtr(dto.MetricType_SUMMARY),
			Metric: []*dto.Metric{
				{
					// Test summary with NaN
					Label: []*dto.LabelPair{
						createLabelPair("tag1", "abc"),
						createLabelPair("tag2", "def"),
					},
					Summary: &dto.Summary{
						SampleCount: uintPtr(1),
						SampleSum:   floatPtr(2),
						Quantile: []*dto.Quantile{
							createQuantile(0.5, 3),
							createQuantile(0.9, 4),
							createQuantile(0.99, math.NaN()),
						},
					},
				},
			},
		},
		output: &Family{
			Name: "summary1",
			Help: "",
			Type: "SUMMARY",
			Metrics: []interface{}{
				Summary{
					Labels: map[string]string{
						"tag1": "abc",
						"tag2": "def",
					},
					Quantiles: map[string]string{
						"0.5":  "3",
						"0.9":  "4",
						"0.99": "NaN",
					},
					Count: "1",
					Sum:   "2",
				},
			},
		},
	},
	{
		name: "test histograms",
		mFamily: &dto.MetricFamily{
			Name: strPtr("histogram1"),
			Type: metricTypePtr(dto.MetricType_HISTOGRAM),
			Metric: []*dto.Metric{
				{
					// Test summary with NaN
					Label: []*dto.LabelPair{
						createLabelPair("tag1", "abc"),
						createLabelPair("tag2", "def"),
					},
					Histogram: &dto.Histogram{
						SampleCount: uintPtr(1),
						SampleSum:   floatPtr(2),
						Bucket: []*dto.Bucket{
							createBucket(250000, 3),
							createBucket(500000, 4),
							createBucket(1e+06, 5),
						},
					},
				},
			},
		},
		output: &Family{
			Name: "histogram1",
			Help: "",
			Type: "HISTOGRAM",
			Metrics: []interface{}{
				Histogram{
					Labels: map[string]string{
						"tag1": "abc",
						"tag2": "def",
					},
					Buckets: map[string]string{
						"250000": "3",
						"500000": "4",
						"1e+06":  "5",
					},
					Count: "1",
					Sum:   "2",
				},
			},
		},
	},
}

func TestConvertToMetricFamily(t *testing.T) {
	for _, tc := range tcs {
		output := newFamily(tc.mFamily)
		if !reflect.DeepEqual(tc.output, output) {
			t.Errorf("test case %s: conversion to metricFamily format failed:\nexpected:\n%s\n\nactual:\n%s",
				tc.name, spew.Sdump(tc.output), spew.Sdump(output))
		}
	}
}

func strPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}

func metricTypePtr(mt dto.MetricType) *dto.MetricType {
	return &mt
}

func uintPtr(u uint64) *uint64 {
	return &u
}

func intPtr(i int64) *int64 {
	return &i
}

func createLabelPair(name string, value string) *dto.LabelPair {
	return &dto.LabelPair{
		Name:  &name,
		Value: &value,
	}
}

func createQuantile(q float64, v float64) *dto.Quantile {
	return &dto.Quantile{
		Quantile: &q,
		Value:    &v,
	}
}

func createBucket(bound float64, count uint64) *dto.Bucket {
	return &dto.Bucket{
		UpperBound:      &bound,
		CumulativeCount: &count,
	}
}
