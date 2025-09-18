package metrics

import (
	"runtime"
	"time"
)

// StartSystemMetricsCollection starts collecting custom system metrics periodically
func StartSystemMetricsCollection() {
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			MemoryUsage.Set(float64(m.Alloc))
		}
	}()
}
