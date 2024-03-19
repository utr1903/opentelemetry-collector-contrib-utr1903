package flexreceiver

import (
	"context"
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/flexreceiver/runner"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/receiver/receivertest"
)

func Test_ReceiverStartAndShutdownIsSuccessful(t *testing.T) {

	cfg := createDefaultConfig().(*Config)
	cfg.Commands = createValidCommands()

	// Instantiate receiver.
	r, err := createLogsReceiver(
		context.Background(),
		receivertest.NewNopCreateSettings(),
		cfg,
		consumertest.NewNop(),
	)
	require.NoError(t, err)
	require.NotNil(t, r)

	// Start receiver
	err = r.Start(context.Background(), componenttest.NewNopHost())
	require.NoError(t, err)

	// Shutdown receiver
	err = r.Shutdown(context.Background())
	require.NoError(t, err)
}

func createValidCommands() []runner.CommandConfig {
	return []runner.CommandConfig{
		{
			Name:      "command1",
			Type:      "table",
			Interval:  "5s",
			Command:   "command1",
			Keys:      "key1,key2,key3",
			Separator: " ",
		},
		{
			Name:     "command2",
			Type:     "json",
			Interval: "5s",
			Command:  "command2",
		},
	}
}
