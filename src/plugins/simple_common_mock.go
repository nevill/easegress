package plugins

import (
	"fmt"
	"strings"

	"github.com/hexdecteam/easegateway-types/pipelines"
	"github.com/hexdecteam/easegateway-types/plugins"
	"github.com/hexdecteam/easegateway-types/task"

	"common"
)

type simpleCommonMockConfig struct {
	common.PluginCommonConfig
	PluginConcerned        string `json:"plugin_concerned"`
	TaskErrorCodeConcerned string `json:"task_error_code_concerned"`
	// TODO: Supports multiple key and value pairs
	MockTaskDataKey   string `json:"mock_task_data_key"`
	MockTaskDataValue string `json:"mock_task_data_value"`

	taskErrorCodeConcerned task.TaskResultCode
}

func simpleCommonMockConfigConstructor() plugins.Config {
	return &simpleCommonMockConfig{
		TaskErrorCodeConcerned: "ResultFlowControl",
	}
}

func (c *simpleCommonMockConfig) Prepare(pipelineNames []string) error {
	err := c.PluginCommonConfig.Prepare(pipelineNames)
	if err != nil {
		return err
	}

	ts := strings.TrimSpace
	c.PluginConcerned = ts(c.PluginConcerned)
	c.MockTaskDataKey = ts(c.MockTaskDataKey)
	c.MockTaskDataValue = ts(c.MockTaskDataValue)

	if len(c.PluginConcerned) == 0 {
		return fmt.Errorf("invalid plugin name")
	}

	if !task.ValidResultCodeName(c.TaskErrorCodeConcerned) {
		return fmt.Errorf("invalid task error code")
	} else {
		c.taskErrorCodeConcerned = task.ResultCodeValue(c.TaskErrorCodeConcerned)
	}

	if len(c.MockTaskDataKey) == 0 {
		return fmt.Errorf("invalid mock task data key")
	}

	return nil
}

////

type simpleCommonMock struct {
	conf *simpleCommonMockConfig
}

func simpleCommonMockConstructor(conf plugins.Config) (plugins.Plugin, plugins.PluginType, error) {
	c, ok := conf.(*simpleCommonMockConfig)
	if !ok {
		return nil, plugins.ProcessPlugin, fmt.Errorf("config type want *simpleCommonMockConfig got %T", conf)
	}

	m := &simpleCommonMock{
		conf: c,
	}

	return m, plugins.ProcessPlugin, nil
}

func (m *simpleCommonMock) Prepare(ctx pipelines.PipelineContext) {
	// Nothing to do.
}

func (m *simpleCommonMock) Run(ctx pipelines.PipelineContext, t task.Task) error {
	t.AddRecoveryFunc("mockBrokenTaskOutput",
		getTaskRecoveryFuncInSimpleCommonMock(m.conf.PluginConcerned, m.conf.taskErrorCodeConcerned,
			m.conf.MockTaskDataKey, m.conf.MockTaskDataValue))
	return nil
}

func (m *simpleCommonMock) Name() string {
	return m.conf.PluginName()
}

func (m *simpleCommonMock) CleanUp(ctx pipelines.PipelineContext) {
	// Nothing to do.
}

func (m *simpleCommonMock) Close() {
	// Nothing to do.
}

////

func getTaskRecoveryFuncInSimpleCommonMock(pluginConcerned string, taskErrorCodeConcerned task.TaskResultCode,
	mockTaskDataKey, mockTaskDataValue string) task.TaskRecovery {

	return func(t task.Task, errorPluginName string) bool {
		if errorPluginName != pluginConcerned || t.ResultCode() != taskErrorCodeConcerned {
			return false
		}

		t.WithValue(mockTaskDataKey, mockTaskDataValue)

		return true
	}
}
