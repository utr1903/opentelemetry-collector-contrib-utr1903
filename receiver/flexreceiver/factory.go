package flexreceiver

import (
	"context"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/flexreceiver/runner"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

const (
	typeStr = "flex"
)

// NewFactory creates a factory for flexreceiver.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		typeStr,
		createDefaultConfig,
		receiver.WithLogs(createLogsReceiver, component.StabilityLevelAlpha))
}

func createDefaultConfig() component.Config {
	return &Config{
		Commands: []runner.CommandConfig{},
	}
}

func createLogsReceiver(_ context.Context, params receiver.CreateSettings, baseCfg component.Config, consumer consumer.Logs) (receiver.Logs, error) {
	if consumer == nil {
		return nil, component.ErrNilNextConsumer
	}

	logger := params.Logger
	cfg := baseCfg.(*Config)

	rcv := &flexReceiver{
		logger:       logger,
		nextConsumer: consumer,
		config:       cfg,
	}

	return rcv, nil
}
