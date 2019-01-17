package plugins

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/megaease/easegateway/pkg/logger"
	"github.com/megaease/easegateway/pkg/pipelines"
	"github.com/megaease/easegateway/pkg/task"
)

type staticProbabilityLimiterConfig struct {
	PluginCommonConfig
	PassPr float32 `json:"pass_pr"`
}

func staticProbabilityLimiterConfigConstructor() Config {
	return new(staticProbabilityLimiterConfig)
}

func (c *staticProbabilityLimiterConfig) Prepare(pipelineNames []string) error {
	err := c.PluginCommonConfig.Prepare(pipelineNames)
	if err != nil {
		return err
	}

	if c.PassPr < 0 || c.PassPr > 1 {
		return fmt.Errorf("invalid passing probability %f", c.PassPr)
	}

	if c.PassPr == 0 {
		logger.Warnf("ZERO passing probability has been applied, no request could be processed!")
	}

	if c.PassPr == 1 {
		logger.Warnf("1.0 passing probability has been applied, no request could be limited!")
	}

	return nil
}

type staticProbabilityLimiter struct {
	conf *staticProbabilityLimiterConfig
}

func staticProbabilityLimiterConstructor(conf Config) (Plugin, PluginType, bool, error) {
	c, ok := conf.(*staticProbabilityLimiterConfig)
	if !ok {
		return nil, ProcessPlugin, false, fmt.Errorf(
			"config type want *staticProbabilityLimiterConfig got %T", conf)
	}

	return &staticProbabilityLimiter{
		conf: c,
	}, ProcessPlugin, false, nil
}

func (l *staticProbabilityLimiter) Prepare(ctx pipelines.PipelineContext) {
	// Nothing to do.
}

func (l *staticProbabilityLimiter) Run(ctx pipelines.PipelineContext, t task.Task) error {
	if rand.Float32() < 1.0-l.conf.PassPr {
		t.SetError(fmt.Errorf("service is unavailable caused by probability limit"), task.ResultFlowControl)
	}
	return nil
}

func (l *staticProbabilityLimiter) Name() string {
	return l.conf.PluginName()
}

func (l *staticProbabilityLimiter) CleanUp(ctx pipelines.PipelineContext) {
	// Nothing to do.
}

func (l *staticProbabilityLimiter) Close() {
	// Nothing to do.
}

////

func init() {
	rand.Seed(time.Now().UnixNano())
}