package runner

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"time"
)

const (
	MIN_INTERVAL_SECONDS = 5

	RUNNER_COMMAND_NAME_IS_NOT_DEFINED                = "runner: command name is not defined"
	RUNNER_COMMAND_TYPE_IS_NOT_DEFINED                = "runner: command type is not defined"
	RUNNER_COMMAND_TYPE_IS_NOT_SUPPORTED              = "runner: command type is not supported"
	RUNNER_COMMAND_TYPE_TABLE_NEEDS_KEYS              = "runner: command type requires keys"
	RUNNER_COMMAND_TYPE_TABLE_NEEDS_SEPARATOR         = "runner: command type requires separator"
	RUNNER_COMMAND_INTERVAL_IS_NOT_DEFINED            = "runner: command interval is not defined"
	RUNNER_COMMAND_INTERVAL_COULD_NOT_BE_PARSED       = "runner: command interval could not be parsed"
	RUNNER_COMMAND_INTERVAL_SHOULD_BE_ABOVE_THRESHOLD = "runner: command interval has to be set to at least 5 seconds (5s)"
	RUNNER_COMMAND_COMMAND_IS_NOT_DEFINED             = "runner: command command is not defined"

	RUNNER_COMMAND_EXECUTION_IS_FAILED = "runner: command execution is failed"
)

func defaultExecuteFunc(command string) ([]byte, error) {
	var shell string
	var param string

	if runtime.GOOS == "windows" {
		shell = "cmd"
		param = "/C"
	} else {
		shell = "/bin/sh"
		param = "-c"
	}

	return exec.Command(shell, param, command).Output()
}

type CommandConfig struct {
	Name     string `mapstructure:"name"`
	Type     string `mapstructure:"type"` // table or json
	Interval string `mapstructure:"interval"`
	Command  string `mapstructure:"command"`

	// Relevant for type=table
	Keys      string `mapstructure:"keys"`
	Separator string `mapstructure:"separator"`
}

type runnerOpts struct {
	Name     string
	Interval string
	Command  string
	Execute  func(string) ([]byte, error)
}

type Runner interface {
	Run(ctx context.Context)
}

// Validates command name.
func ValidateCommandName(command CommandConfig) error {
	// Check if name is defined.
	if command.Name == "" {
		fmt.Println("Command is not given a name.")
		return fmt.Errorf(RUNNER_COMMAND_NAME_IS_NOT_DEFINED)
	}
	return nil
}

// Validates command type.
func ValidateCommandType(command CommandConfig) error {
	// Check if type is defined.
	if command.Type == "" {
		fmt.Println("Command is not given a type.")
		return fmt.Errorf(RUNNER_COMMAND_TYPE_IS_NOT_DEFINED)
	}

	// Check if type is supported.
	if command.Type != RUNNER_TYPE_TABLE && command.Type != RUNNER_TYPE_JSON {
		fmt.Println("Type of given command [" + command.Name + "] is not suported: " + command.Type + ". Supported values are: " + RUNNER_TYPE_TABLE + ", " + RUNNER_TYPE_JSON)
		return fmt.Errorf(RUNNER_COMMAND_TYPE_IS_NOT_SUPPORTED)
	}

	if command.Type == RUNNER_TYPE_TABLE {
		// Check if keys is defined.
		if command.Keys == "" {
			fmt.Println("Keys field is required for given command [" + command.Name + "] of type [" + command.Type + "]")
			return fmt.Errorf(RUNNER_COMMAND_TYPE_TABLE_NEEDS_KEYS)
		}

		// Check if separator is defined.
		if command.Keys == "" {
			fmt.Println("Separator field is required for given command [" + command.Name + "] of type [" + command.Type + "]")
			return fmt.Errorf(RUNNER_COMMAND_TYPE_TABLE_NEEDS_SEPARATOR)
		}
	}
	return nil
}

// Validates command interval.
func ValidateCommandInterval(command CommandConfig) error {
	// Check if interval is defined.
	if command.Interval == "" {
		fmt.Println("Command is not given a interval.")
		return fmt.Errorf(RUNNER_COMMAND_INTERVAL_IS_NOT_DEFINED)
	}

	// Parse interval.
	interval, err := time.ParseDuration(command.Interval)
	if err != nil {
		return fmt.Errorf(RUNNER_COMMAND_INTERVAL_COULD_NOT_BE_PARSED)
	}

	// Check if interval is above the threshold.
	if interval.Seconds() < MIN_INTERVAL_SECONDS {
		return fmt.Errorf(RUNNER_COMMAND_INTERVAL_SHOULD_BE_ABOVE_THRESHOLD)
	}
	return nil
}

// Validates command command.
func ValidateCommandCommand(command CommandConfig) error {
	// Check if command is defined.
	if command.Command == "" {
		fmt.Println("Command is not given a command.")
		return fmt.Errorf(RUNNER_COMMAND_COMMAND_IS_NOT_DEFINED)
	}
	return nil
}
