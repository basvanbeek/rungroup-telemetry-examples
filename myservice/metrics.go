package myservice

import (
	"sync"

	"github.com/tetratelabs/telemetry"
)

var (
	mInit            sync.Once
	lvServiceID      telemetry.LabelValue
	lRequestMethod   telemetry.Label
	lRequestStatus   telemetry.Label
	mRequestCount    telemetry.Metric
	mRequestDuration telemetry.Metric
)

func SyncMetrics(m telemetry.MetricSink, serviceInstanceID string) {
	mInit.Do(func() {
		// initialize our label dimensions
		lService := m.NewLabel("serviceInstance")
		lRequestMethod = m.NewLabel("method")
		lRequestStatus = m.NewLabel("status")

		// use our service's instance id as a dimension
		lvServiceID = lService.Insert(serviceInstanceID)

		// initialize our metrics
		mRequestCount = m.NewSum(
			"requests_total",
			"total amount of requests",
			telemetry.WithLabels(lService, lRequestMethod, lRequestStatus),
		).With(lvServiceID)
		mRequestDuration = m.NewDistribution(
			"requests_duration",
			"latency in milliseconds per request",
			[]float64{0, 25, 50, 75, 100, 200, 400, 600, 800, 1000, 2000, 4000, 6000},
			telemetry.WithLabels(lService, lRequestMethod, lRequestStatus),
		).With(lvServiceID)
	})
}
