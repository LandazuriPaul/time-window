package validator

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v3"
	"time"
)

var cronParser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

type TimeWindow struct {
	CronExpression string        `yaml:"cronExpression"`
	Duration       time.Duration `yaml:"duration"`
	Name           string        `yaml:"name"`
	Timezone       string        `yaml:"timezone"`
	CronSchedule   cron.Schedule
}

func NewTimeWindows(config string) ([]*TimeWindow, error) {
	var tws []*TimeWindow
	err := yaml.Unmarshal([]byte(config), &tws)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling the YAML: %w", err)
	}

	for _, tw := range tws {
		cronErr := tw.parseCronExpression()
		if cronErr != nil {
			return nil, fmt.Errorf("parsing the CRON expression of time window '%s': %w", tw.Name, cronErr)
		}
	}

	return tws, nil
}

func (tw *TimeWindow) parseCronExpression() error {
	schedule, err := cronParser.Parse(tw.CronExpression)
	if err != nil {
		return err
	}
	tw.CronSchedule = schedule
	return nil
}

func (tw *TimeWindow) isTimeIn(timestamp time.Time) bool {
	// FIXME: Timezone!
	if tw.CronSchedule.Next(timestamp.Add(-tw.Duration)).Before(timestamp) {
		return true
	}
	return false
}
