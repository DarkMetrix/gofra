package performance

import (
	"time"
	"runtime"

	log "github.com/cihub/seelog"
	monitor "github.com/DarkMetrix/gofra/common/monitor/statsd"
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

