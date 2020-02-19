package performance

import (
	"runtime"
	"time"

	monitor "github.com/DarkMetrix/gofra/pkg/monitor/statsd"
	log "github.com/cihub/seelog"
)

func BeginGoroutinePerformanceMonitorWithStatsd() {
	log.Infof("Begin Goroutine Performance Monitor with Statsd")

	ticker := time.NewTicker(time.Second * 10)

	for {
		select {
		case t := <-ticker.C:
			log.Tracef("Ticker triggered! time:%v", t)

			monitor.Gauge("/application/performance/goroutine,type=number", runtime.NumGoroutine())
		}
	}
}

func BeginGoroutinePerformanceMonitorWithLog() {
	log.Infof("Begin Goroutine Performance Monitor with Log")

	ticker := time.NewTicker(time.Second * 10)

	for {
		select {
		case t := <-ticker.C:
			log.Tracef("Ticker triggered! time:%v", t)

			log.Infof("Current goroutine count:%v", runtime.NumGoroutine())
		}
	}
}

