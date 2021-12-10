package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/gofrs/uuid"
	"github.com/tetratelabs/run"
	"github.com/tetratelabs/run/pkg/signal"
	log "github.com/tetratelabs/telemetry-gokit-log"
	svclogs "github.com/tetratelabs/telemetry-gokit-log/group"
	metrics "github.com/tetratelabs/telemetry-opencensus"
	svcmetrics "github.com/tetratelabs/telemetry-opencensus/group"

	"github.com/basvanbeek/rungroup-telemetry-examples/myservice"
)

func main() {
	// app instance id
	instanceID := uuid.Must(uuid.NewV4())

	// create our scoped logger
	sm := log.NewManager(log.NewSyncLogfmt(os.Stdout))
	sm.SetDefaultOutputLevel(log.Debug)

	// initialize the logging scopes we want to have for this application
	var (
		logMetrics   = sm.RegisterScope("metrics", "metrics handling related messages")
		logRunGroup  = sm.RegisterScope("metrics", "metrics handling related messages")
		logMyService = sm.RegisterScope("myservice", "myservice related messages")
	)

	// create our metrics sink
	ms := metrics.New(logMetrics)

	// create our run group
	g := run.Group{Logger: logRunGroup}

	// initialize MyService
	myService := &myservice.MyService{
		Logger:     logMyService,
		Metrics:    ms,
		InstanceID: instanceID.String(),
	}

	// register the configuration objects and services we wish to run
	g.Register(
		&signal.Handler{},                   // signal handler
		svclogs.New(sm),                     // scoped logger configuration service
		svcmetrics.New(svcmetrics.Config{}), // metrics exporter service
		myService,                           // MyService
	)

	// run our configuration phase
	if err := g.RunConfig(); err != nil {
		if errors.Is(err, run.ErrBailEarlyRequest) {
			os.Exit(0)
		}
		fmt.Printf("unexpected exit: %+v\n", err)
		os.Exit(-1)
	}

	// output our registered loggers
	sm.PrintRegisteredScopes()

	// run our registered run group services
	if err := g.Run(); err != nil {
		fmt.Printf("unexpected exit: %+v\n", err)
		os.Exit(-1)
	}
}
