package flexreceiver

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/flexreceiver/runner"
)

type Config struct {
	Commands []runner.CommandConfig `mapstructure:"commands"`
}

// Validate checks if the receiver configuration is valid.
func (cfg *Config) Validate() error {
	for _, command := range cfg.Commands {

		var err error

		// Check name.
		err = runner.ValidateCommandName(command)
		if err != nil {
			return err
		}

		// Check type.
		err = runner.ValidateCommandType(command)
		if err != nil {
			return err
		}

		// Check interval.
		err = runner.ValidateCommandInterval(command)
		if err != nil {
			return err
		}

		// Check command.
		err = runner.ValidateCommandCommand(command)
		if err != nil {
			return err
		}
	}

	return nil
}
