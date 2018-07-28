package performance

import (
	"time"
	"runtime"

	log "github.com/cihub/seelog"
	monitor "github.com/DarkMetrix/gofra/common/monitor/statsd"
)

func BeginMemoryPerformanceMonitorWithStatsd() {
	log.Infof("Begin Memory Performance Monitor with Statsd")

	ticker := time.NewTicker(time.Second * 10)

	for {
		var lastMemStats runtime.MemStats

		select {
		case t := <- ticker.C:
			log.Tracef("Ticker triggered! time:%v", t)

			memStats := &runtime.MemStats{}
			runtime.ReadMemStats(memStats)

			monitor.Gauge("/application/performance/memory,type=heap_objects", int64(memStats.HeapObjects))
			monitor.Gauge("/application/performance/memory,type=heap_alloc_bytes", int64(memStats.HeapAlloc))
			monitor.Gauge("/application/performance/memory,type=heap_sys_bytes", int64(memStats.HeapSys))
			monitor.Gauge("/application/performance/memory,type=heap_idle_bytes", int64(memStats.HeapIdle))
			monitor.Gauge("/application/performance/memory,type=heap_in_use_bytes", int64(memStats.HeapInuse))
			monitor.Gauge("/application/performance/memory,type=heap_released_bytes", int64(memStats.HeapReleased))

			if memStats.NumGC < lastMemStats.NumGC {
				monitor.Gauge("/application/performance/gc,type=gc_counts", int64(0))
			} else {
				monitor.Gauge("/application/performance/gc,type=gc_counts", int64(memStats.NumGC - lastMemStats.NumGC))
			}

			monitor.Gauge("/application/performance/gc,type=gc_pause", int64(memStats.PauseNs[(memStats.NumGC + 255) % 256]))

			lastMemStats = *memStats
		}
	}
}

func BeginMemoryPerformanceMonitorWithLog() {
	log.Infof("Begin Memory Performance Monitor with Log")

	ticker := time.NewTicker(time.Second * 10)

	for {
		var lastMemStats runtime.MemStats

		select {
		case t := <- ticker.C:
			log.Tracef("Ticker triggered! time:%v", t)

			memStats := &runtime.MemStats{}
			runtime.ReadMemStats(memStats)

			log.Infof(`heap_objects:%v\r\n heap_alloc_bytes:%v\r\n heap_sys_bytes:%v\r\n heap_idle_bytes:%v\r\n heap_in_user_bytes:%v\r\n heap_released_bytes:%v`,
				memStats.HeapObjects, memStats.HeapAlloc, memStats.HeapSys, memStats.HeapIdle, memStats.HeapInuse, memStats.HeapReleased)

			if memStats.NumGC < lastMemStats.NumGC {
				log.Infof("gc_counts:%v gc_pause:%v", 0, memStats.PauseNs[(memStats.NumGC + 255) % 256])
			} else {
				log.Infof("gc_counts:%v gc_pause:%v", memStats.NumGC - lastMemStats.NumGC, memStats.PauseNs[(memStats.NumGC + 255) % 256])
			}

			lastMemStats = *memStats
		}
	}
}
