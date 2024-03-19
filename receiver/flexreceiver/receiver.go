package flexreceiver

import (
	"context"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/flexreceiver/runner"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"
)

type flexReceiver struct {
	host         component.Host
	cancel       context.CancelFunc
	logger       *zap.Logger
	nextConsumer consumer.Logs
	config       *Config
}

// Starts the receiver.
func (r *flexReceiver) Start(ctx context.Context, host component.Host) error {
	r.host = host
	ctx, r.cancel = context.WithCancel(ctx)

	r.logger.Info("Starting receiver...")

	for _, command := range r.config.Commands {

		var rnr runner.Runner
		if command.Type == runner.RUNNER_TYPE_TABLE {
			rnr = r.instantiateRunnerTable(command)

		} else if command.Type == runner.RUNNER_TYPE_JSON {
			rnr = r.instantiateRunnerJson(command)
		}

		go rnr.Run(ctx)
	}
	return nil
}

// Instantiates a runner for table.
func (r *flexReceiver) instantiateRunnerTable(command runner.CommandConfig) runner.Runner {
	return runner.NewRunnerTable(r.logger, r.nextConsumer,
		runner.WithTableName(command.Name),
		runner.WithTableInterval(command.Interval),
		runner.WithTableCommand(command.Command),
		runner.WithTableKeys(command.Keys),
		runner.WithTableSeparator(command.Separator),
	)
}

// Instantiates a runner for json.
func (r *flexReceiver) instantiateRunnerJson(command runner.CommandConfig) runner.Runner {
	return runner.NewRunnerTable(r.logger, r.nextConsumer,
		runner.WithTableName(command.Name),
		runner.WithTableInterval(command.Interval),
		runner.WithTableCommand(command.Command),
	)
}

// Shuts down the receiver.
func (r *flexReceiver) Shutdown(ctx context.Context) error {
	r.cancel()
	return nil
}
