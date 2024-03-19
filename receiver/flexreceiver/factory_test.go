package flexreceiver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/receiver/receivertest"
)

func Test_DefaultConfigHasEmptyCommands(t *testing.T) {
	cfg := createDefaultConfig()
	rCfg, ok := cfg.(*Config)
	require.True(t, ok)

	require.Equal(t, len(rCfg.Commands), 0)
}

func Test_ReceiverCreationFailsWhenNoConsumerIsGiven(t *testing.T) {
	cfg := createDefaultConfig().(*Config)

	// Fails without consumer.
	_, err := createLogsReceiver(
		context.Background(),
		receivertest.NewNopCreateSettings(),
		cfg,
		nil,
	)
	require.Error(t, err)
}

func Test_ReceiverCreationAndStarSucceeds(t *testing.T) {
	cfg := createDefaultConfig().(*Config)

	// Succeeds with consumer.
	r, err := createLogsReceiver(
		context.Background(),
		receivertest.NewNopCreateSettings(),
		cfg,
		consumertest.NewNop(),
	)
	require.NoError(t, err)

	// Succeeds when starts.
	err = r.Start(context.Background(), componenttest.NewNopHost())
	require.NoError(t, err)
}
