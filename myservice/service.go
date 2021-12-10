package myservice

import (
	"math/rand"
	"time"

	"github.com/tetratelabs/run"
	"github.com/tetratelabs/telemetry"
)

var (
	_ run.Config    = (*MyService)(nil)
	_ run.PreRunner = (*MyService)(nil)
	_ run.Service   = (*MyService)(nil)
)

type MyService struct {
	Logger     telemetry.Logger
	Metrics    telemetry.MetricSink
	InstanceID string

	cErr chan error
}

func (m MyService) Name() string {
	return "my-service"
}

func (m *MyService) FlagSet() *run.FlagSet {
	fs := run.NewFlagSet("MyService Options")

	return fs
}

func (m *MyService) Validate() error {
	return nil
}

func (m *MyService) PreRun() error {
	m.cErr = make(chan error)
	SyncMetrics(m.Metrics, m.InstanceID)
	return nil
}

func (m *MyService) Serve() error {
	callNames := []string{"create", "read", "update", "delete"}
	callStatuses := []string{"200", "400", "500"}
	for {
		select {
		case err := <-m.cErr:
			return err
		default:
		}
		// fake a method call...
		callName := callNames[rand.Intn(len(callNames))]
		callDuration := float64(rand.Int31n(10000))
		callStatus := callNames[rand.Intn(len(callStatuses))]
		mRequestCount.With(
			lRequestMethod.Insert(callName),
			lRequestStatus.Insert(callStatus),
		).Record(1)
		mRequestDuration.With(
			lRequestMethod.Insert(callName),
			lRequestStatus.Insert(callStatus),
		).Record(callDuration)

		// try again shortly
		time.Sleep(100 * time.Millisecond)
	}
}

func (m *MyService) GracefulStop() {
	close(m.cErr)
}
