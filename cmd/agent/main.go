package main

import (
	"fmt"
	"time"

	"github.com/rekh/_temp/goprocoursedevopsproject/internal/agent"
)

var pollCount agent.Counter = 0

func main() {
	m := agent.Metrics{}

	intervals := agent.GetConfig().Intervals
	pollInterval := intervals.Poll * time.Second
	reportInterval := intervals.Report * time.Second

	poll := time.NewTicker(pollInterval)
	report := time.NewTicker(reportInterval)

	for {
		select {
		case <-poll.C:
			m.Get()
			fmt.Println(m)
			poll = time.NewTicker(pollInterval)

		case <-report.C:
			err := m.Send()
			if err != nil {
				panic(err)
			}
			report = time.NewTicker(reportInterval)
		}
	}
}
