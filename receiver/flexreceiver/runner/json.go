package runner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

const RUNNER_TYPE_JSON = "json"

const RUNNER_COMMAND_JSON_PARSING_IS_FAILED = "runner: command json parsing is failed"

type runnerJsonOpts struct {
	runnerOpts
}

type runnerJsonFunc func(*runnerJsonOpts)

type RunnerJson struct {
	logger       *zap.Logger
	nextConsumer consumer.Logs
	opts         *runnerJsonOpts
}

func NewRunnerJson(logger *zap.Logger, nextConsumer consumer.Logs, optFuncs ...runnerJsonFunc) *RunnerJson {

	// Initialize opts.
	opts := &runnerJsonOpts{
		runnerOpts: runnerOpts{
			Execute: defaultExecuteFunc,
		},
	}

	// Apply opts
	for _, f := range optFuncs {
		f(opts)
	}

	// Instantiate runner
	return &RunnerJson{
		logger:       logger,
		nextConsumer: nextConsumer,
		opts:         opts,
	}
}

// Set name
func WithJsonName(name string) runnerJsonFunc {
	return func(opts *runnerJsonOpts) {
		opts.Name = name
	}
}

// Set interval
func WithJsonInterval(interval string) runnerJsonFunc {
	return func(opts *runnerJsonOpts) {
		opts.Interval = interval
	}
}

// Set command
func WithJsonCommand(command string) runnerJsonFunc {
	return func(opts *runnerJsonOpts) {
		opts.Command = command
	}
}

// Set execute.
func WithJsonExecute(execute func(string) ([]byte, error)) runnerJsonFunc {
	return func(opts *runnerJsonOpts) {
		opts.Execute = execute
	}
}

func (r *RunnerJson) Run(ctx context.Context) {
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

func (r *RunnerJson) run() (*plog.Logs, error) {
	logs := plog.NewLogs()
	resourceLog := logs.ResourceLogs().AppendEmpty()
	logResource := resourceLog.Resource()

	// Add resource attributes
	attrs := logResource.Attributes()
	attrs.PutStr("command.name", r.opts.Name)

	scopeLogs := resourceLog.ScopeLogs().AppendEmpty()

	out, err := r.opts.Execute(r.opts.Command)
	if err != nil {
		r.logger.Error("Executing command is failed.",
			zap.String("name", r.opts.Name),
			zap.String("command", r.opts.Command),
			zap.String("error", err.Error()),
		)
		return nil, errors.New(RUNNER_COMMAND_EXECUTION_IS_FAILED)
	}

	outJson := map[string]interface{}{}
	err = json.Unmarshal(out, &outJson)
	if err != nil {
		r.logger.Error("Parsing JSON is failed.",
			zap.String("name", r.opts.Name),
			zap.String("command", r.opts.Command),
			zap.String("error", err.Error()),
		)
		return nil, errors.New(RUNNER_COMMAND_JSON_PARSING_IS_FAILED)
	}

	logRecord := scopeLogs.LogRecords().AppendEmpty()
	logRecord.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))

	// No specific message is required.
	logRecord.Body().SetStr("")

	// Add attributes
	for k, v := range outJson {
		logRecord.Attributes().PutStr(k, fmt.Sprintf("%v", v))
	}

	return &logs, nil
}
