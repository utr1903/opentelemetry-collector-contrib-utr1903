package flexreceiver

import (
	"path/filepath"
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/flexreceiver/runner"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/confmaptest"
)

func Test_CommandNameIsNotDefined(t *testing.T) {

	// Load config
	cm := loadConfig(t)

	// Create default config.
	cfg := NewFactory().CreateDefaultConfig()

	// Unmarshall config.
	unmarshallConfig(t, "no_name", cm, &cfg)

	require.Equal(t,
		component.ValidateConfig(cfg).Error(),
		runner.RUNNER_COMMAND_NAME_IS_NOT_DEFINED,
	)
}

func Test_CommandTypeIsNotDefined(t *testing.T) {

	// Load config
	cm := loadConfig(t)

	// Create default config.
	cfg := NewFactory().CreateDefaultConfig()

	// Unmarshall config.
	unmarshallConfig(t, "no_type", cm, &cfg)

	require.Equal(t,
		component.ValidateConfig(cfg).Error(),
		runner.RUNNER_COMMAND_TYPE_IS_NOT_DEFINED,
	)
}

func Test_CommandTypeIsNotSupported(t *testing.T) {

	// Load config
	cm := loadConfig(t)

	// Create default config.
	cfg := NewFactory().CreateDefaultConfig()

	// Unmarshall config.
	unmarshallConfig(t, "invalid_type", cm, &cfg)

	require.Equal(t,
		component.ValidateConfig(cfg).Error(),
		runner.RUNNER_COMMAND_TYPE_IS_NOT_SUPPORTED,
	)
}

func Test_CommandIntervalIsNotDefined(t *testing.T) {

	// Load config
	cm := loadConfig(t)

	// Create default config.
	cfg := NewFactory().CreateDefaultConfig()

	// Unmarshall config.
	unmarshallConfig(t, "no_interval", cm, &cfg)

	require.Equal(t,
		component.ValidateConfig(cfg).Error(),
		runner.RUNNER_COMMAND_INTERVAL_IS_NOT_DEFINED,
	)
}

func Test_CommandIntervalCouldNotBeParsed(t *testing.T) {

	// Load config
	cm := loadConfig(t)

	// Create default config.
	cfg := NewFactory().CreateDefaultConfig()

	// Unmarshall config.
	unmarshallConfig(t, "invalid_interval", cm, &cfg)

	require.Equal(t,
		component.ValidateConfig(cfg).Error(),
		runner.RUNNER_COMMAND_INTERVAL_COULD_NOT_BE_PARSED,
	)
}

func Test_CommandIntervalShouldBeAboveThreshold(t *testing.T) {

	// Load config
	cm := loadConfig(t)

	// Create default config.
	cfg := NewFactory().CreateDefaultConfig()

	// Unmarshall config.
	unmarshallConfig(t, "invalid_interval_threshold", cm, &cfg)

	require.Equal(t,
		component.ValidateConfig(cfg).Error(),
		runner.RUNNER_COMMAND_INTERVAL_SHOULD_BE_ABOVE_THRESHOLD,
	)
}

func Test_CommandCommandIsNotDefined(t *testing.T) {

	// Load config
	cm := loadConfig(t)

	// Create default config.
	cfg := NewFactory().CreateDefaultConfig()

	// Unmarshall config.
	unmarshallConfig(t, "no_command", cm, &cfg)

	require.Equal(t,
		component.ValidateConfig(cfg).Error(),
		runner.RUNNER_COMMAND_COMMAND_IS_NOT_DEFINED,
	)
}

// Loads config.
func loadConfig(t *testing.T) *confmap.Conf {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)
	return cm
}

// Unmarshalls config.
func unmarshallConfig(t *testing.T, cfgName string, cm *confmap.Conf, cfg *component.Config) {
	id := component.NewIDWithName(typeStr, cfgName)
	sub, err := cm.Sub(id.String())
	require.NoError(t, err)
	require.NoError(t, component.UnmarshalConfig(sub, cfg))
}
