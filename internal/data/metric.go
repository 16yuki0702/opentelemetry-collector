// Copyright 2020 OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package data

import (
	"github.com/golang/protobuf/proto"
	otlpmetrics "github.com/open-telemetry/opentelemetry-proto/gen/go/metrics/v1"
)

// This file defines in-memory data structures to represent metrics.
// For the proto representation see https://github.com/open-telemetry/opentelemetry-proto/blob/master/opentelemetry/proto/metrics/v1/metrics.proto

// MetricData is the top-level struct that is propagated through the metrics pipeline.
// This is the newer version of consumerdata.MetricsData, but uses more efficient
// in-memory representation.
//
// This is a reference type (like builtin map).
//
// Must use NewMetricData functions to create new instances.
// Important: zero-initialized instance is not valid for use.
type MetricData struct {
	orig *[]*otlpmetrics.ResourceMetrics
}

// MetricDataFromOtlp creates the internal MetricData representation from the OTLP.
func MetricDataFromOtlp(orig []*otlpmetrics.ResourceMetrics) MetricData {
	return MetricData{&orig}
}

// MetricDataToOtlp converts the internal MetricData to the OTLP.
func MetricDataToOtlp(md MetricData) []*otlpmetrics.ResourceMetrics {
	return *md.orig
}

// NewMetricData creates a new MetricData.
func NewMetricData() MetricData {
	orig := []*otlpmetrics.ResourceMetrics(nil)
	return MetricData{&orig}
}

// Clone returns a copy of MetricData.
func (md MetricData) Clone() MetricData {
	otlp := MetricDataToOtlp(md)
	resourceMetricsClones := make([]*otlpmetrics.ResourceMetrics, 0, len(otlp))
	for _, resourceMetrics := range otlp {
		resourceMetricsClones = append(resourceMetricsClones,
			proto.Clone(resourceMetrics).(*otlpmetrics.ResourceMetrics))
	}
	return MetricDataFromOtlp(resourceMetricsClones)
}

func (md MetricData) ResourceMetrics() ResourceMetricsSlice {
	return newResourceMetricsSlice(md.orig)
}

func (md MetricData) SetResourceMetrics(v ResourceMetricsSlice) {
	*md.orig = *v.orig
}

// MetricCount calculates the total number of metrics.
func (md MetricData) MetricCount() int {
	metricCount := 0
	rms := md.ResourceMetrics()
	for i := 0; i < rms.Len(); i++ {
		ils := rms.Get(i).InstrumentationLibraryMetrics()
		for j := 0; j < ils.Len(); j++ {
			metricCount += ils.Get(j).Metrics().Len()
		}
	}
	return metricCount
}

type MetricType otlpmetrics.MetricDescriptor_Type

const (
	MetricTypeUnspecified         MetricType = MetricType(otlpmetrics.MetricDescriptor_UNSPECIFIED)
	MetricTypeGaugeInt64          MetricType = MetricType(otlpmetrics.MetricDescriptor_GAUGE_INT64)
	MetricTypeGaugeDouble         MetricType = MetricType(otlpmetrics.MetricDescriptor_GAUGE_DOUBLE)
	MetricTypeGaugeHistogram      MetricType = MetricType(otlpmetrics.MetricDescriptor_GAUGE_HISTOGRAM)
	MetricTypeCounterInt64        MetricType = MetricType(otlpmetrics.MetricDescriptor_COUNTER_INT64)
	MetricTypeCounterDouble       MetricType = MetricType(otlpmetrics.MetricDescriptor_COUNTER_DOUBLE)
	MetricTypeCumulativeHistogram MetricType = MetricType(otlpmetrics.MetricDescriptor_CUMULATIVE_HISTOGRAM)
	MetricTypeSummary             MetricType = MetricType(otlpmetrics.MetricDescriptor_SUMMARY)
)
