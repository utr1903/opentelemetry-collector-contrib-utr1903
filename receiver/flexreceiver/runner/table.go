package runner

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

const RUNNER_TYPE_TABLE = "table"

type runnerTableOpts struct {
	runnerOpts

	Keys      string
	Separator string
}

type runnerTableFunc func(*runnerTableOpts)

type RunnerTable struct {
	logger       *zap.Logger
	nextConsumer consumer.Logs

	opts *runnerTableOpts
}

func NewRunnerTable(logger *zap.Logger, nextConsumer consumer.Logs, optFuncs ...runnerTableFunc) *RunnerTable {

	// Initialize opts.
	opts := &runnerTableOpts{
		runnerOpts: runnerOpts{
			Execute: defaultExecuteFunc,
		},
	}

	// Apply opts.
	for _, f := range optFuncs {
		f(opts)
	}

	// Instantiate runner.
	return &RunnerTable{
		logger:       logger,
		nextConsumer: nextConsumer,
		opts:         opts,
	}
}

// Set name.
func WithTableName(name string) runnerTableFunc {
	return func(opts *runnerTableOpts) {
		opts.Name = name
	}
}

// Set interval.
func WithTableInterval(interval string) runnerTableFunc {
	return func(opts *runnerTableOpts) {
		opts.Interval = interval
	}
}

// Set command.
func WithTableCommand(command string) runnerTableFunc {
	return func(opts *runnerTableOpts) {
		opts.Command = command
	}
}

// Set keys.
func WithTableKeys(keys string) runnerTableFunc {
	return func(opts *runnerTableOpts) {
		opts.Keys = keys
	}
}

// Set separator.
func WithTableSeparator(separator string) runnerTableFunc {
	return func(opts *runnerTableOpts) {
		opts.Separator = separator
	}
}

// Set execute.
func WithTableExecute(execute func(string) ([]byte, error)) runnerTableFunc {
	return func(opts *runnerTableOpts) {
		opts.Execute = execute
	}
}

func (r *RunnerTable) Run(ctx context.Context) {
	interval, _ := time.ParseDuration(r.opts.Interval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			logs, err := r.run()
			if err != nil {
				continue
			}
			r.nextConsumer.ConsumeLogs(ctx, *logs)
		case <-ctx.Done():
			return
		}
	}
}

func (r *RunnerTable) run() (*plog.Logs, error) {
	logs := plog.NewLogs()
	resourceLog := logs.ResourceLogs().AppendEmpty()
	logResource := resourceLog.Resource()

	// Add resource attributes
	attrs := logResource.Attributes()
	attrs.PutStr("command.name", r.opts.Name)

	scopeLogs := resourceLog.ScopeLogs().AppendEmpty()

	colKeys := strings.Split(r.opts.Keys, ",")
	out, err := r.opts.Execute(r.opts.Command)
	// out, err := exec.Command("bash", "-c", r.opts.Command).Output()
	if err != nil {
		r.logger.Error("Executing command is failed.",
			zap.String("name", r.opts.Name),
			zap.String("command", r.opts.Command),
			zap.String("error", err.Error()),
		)
		return nil, errors.New(RUNNER_COMMAND_EXECUTION_IS_FAILED)
	}

	lines := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	for _, line := range lines {
		logRecord := scopeLogs.LogRecords().AppendEmpty()
		logRecord.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))

		colVals := strings.Split(line, r.opts.Separator)

		// No specific message is required.
		logRecord.Body().SetStr("")

		// Add attributes
		for i := range colKeys {
			logRecord.Attributes().PutStr(colKeys[i], colVals[i])
		}
	}
	return &logs, nil
}
