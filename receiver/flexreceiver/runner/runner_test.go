package runner

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.uber.org/zap"
)

func Test_RunnersFailToExecuteCommands(t *testing.T) {

	// Create a valid config.
	cfg := createValidConfig()

	rnrTable := NewRunnerTable(zap.NewNop(), consumertest.NewNop(),
		WithTableName(cfg[0].Name),
		WithTableInterval(cfg[0].Interval),
		WithTableCommand(cfg[0].Command),
		WithTableKeys(cfg[0].Keys),
		WithTableSeparator(cfg[0].Separator),
		WithTableExecute(
			func(command string) ([]byte, error) {
				return nil, errors.New(RUNNER_COMMAND_EXECUTION_IS_FAILED)
			},
		),
	)

	logsTable, err := rnrTable.run()
	require.Error(t, err)
	require.Equal(t, err.Error(), RUNNER_COMMAND_EXECUTION_IS_FAILED)
	require.Nil(t, logsTable)

	rnrJson := NewRunnerJson(zap.NewNop(), consumertest.NewNop(),
		WithJsonName(cfg[0].Name),
		WithJsonInterval(cfg[0].Interval),
		WithJsonCommand(cfg[0].Command),
		WithJsonExecute(
			func(command string) ([]byte, error) {
				return nil, errors.New(RUNNER_COMMAND_EXECUTION_IS_FAILED)
			},
		),
	)

	logsJson, err := rnrJson.run()
	require.Error(t, err)
	require.Equal(t, err.Error(), RUNNER_COMMAND_EXECUTION_IS_FAILED)
	require.Nil(t, logsJson)
}

func Test_RunnersSucceedToExecuteCommands(t *testing.T) {

	// Create a valid config.
	cfg := createValidConfig()

	// Check table.
	resTable := mockCommandTableResponse()
	rnrTable := NewRunnerTable(zap.NewNop(), consumertest.NewNop(),
		WithTableName(cfg[0].Name),
		WithTableInterval(cfg[0].Interval),
		WithTableCommand(cfg[0].Command),
		WithTableKeys(cfg[0].Keys),
		WithTableSeparator(cfg[0].Separator),
		WithTableExecute(
			func(command string) ([]byte, error) {
				out := ""
				for row := range resTable {
					for col := range resTable[row] {
						if col != len(resTable[row])-1 {
							out = out + resTable[row][col] + " "
						} else {
							out = out + resTable[row][col]
						}
					}
					out = out + "\n"
				}
				return []byte(out), nil
			},
		),
	)

	// Check error.
	logsTable, err := rnrTable.run()
	require.NoError(t, err)
	require.Nil(t, err)

	// Check resource attributes.
	logsTableResourceAttr, _ := logsTable.ResourceLogs().At(0).Resource().Attributes().Get("command.name")
	require.Equal(t, cfg[0].Command, logsTableResourceAttr.AsString())

	// Check attributes.
	keys := strings.Split(cfg[0].Keys, ",")
	for row := range resTable {
		for col := range resTable[row] {
			exp := resTable[row][col]
			act, _ := logsTable.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(row).Attributes().Get(keys[col])
			require.Equal(t, exp, act.AsString())
		}
	}

	// Check json.
	resJson := mockCommandJsonResponse()
	rnrJson := NewRunnerJson(zap.NewNop(), consumertest.NewNop(),
		WithJsonName(cfg[0].Name),
		WithJsonInterval(cfg[0].Interval),
		WithJsonCommand(cfg[0].Command),
		WithJsonExecute(
			func(command string) ([]byte, error) {
				return json.Marshal(resJson)
			},
		),
	)

	// Check error.
	logsJson, err := rnrJson.run()
	require.NoError(t, err)
	require.Nil(t, err)

	// Check resource attributes.
	logsJsonResourceAttr, _ := logsJson.ResourceLogs().At(0).Resource().Attributes().Get("command.name")
	require.Equal(t, cfg[0].Command, logsJsonResourceAttr.AsString())

	// Check attributes.
	for key, val := range resJson {
		act, _ := logsJson.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Attributes().Get(key)
		require.Equal(t, fmt.Sprintf("%v", val), act.AsString())
	}
}

func createValidConfig() []CommandConfig {
	cfg := []CommandConfig{
		{
			Name:      "command1",
			Type:      RUNNER_TYPE_TABLE,
			Interval:  "5s",
			Command:   "command1",
			Keys:      "permissions,owner,group,size,file",
			Separator: " ",
		},
		{
			Name:     "command2",
			Type:     RUNNER_TYPE_JSON,
			Interval: "5s",
			Command:  "command2",
		},
	}
	return cfg
}

func mockCommandTableResponse() [][]string {
	return [][]string{
		{"-rw-r--r--", "user1", "group1", "579", "file1.txt"},
		{"-rw-r--rw-", "user2", "group2", "107", "file2.txt"},
	}
}

func mockCommandJsonResponse() map[string]interface{} {
	return map[string]interface{}{
		"key1": "val1",
		"key2": 0.2,
		"key3": 3,
	}
}
