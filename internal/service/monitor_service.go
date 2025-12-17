package service

import (
	"runtime"
	"time"

	"go-admin/internal/model"
)

// MonitorService handles system monitoring
type MonitorService struct {
	startTime time.Time
}

// NewMonitorService creates a new monitor service
func NewMonitorService() *MonitorService {
	return &MonitorService{
		startTime: time.Now(),
	}
}

// GetSystemInfo retrieves system information
func (s *MonitorService) GetSystemInfo() *model.SystemInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return &model.SystemInfo{
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
		GoVersion:  runtime.Version(),
		AppVersion: "1.0.0", // This should be configurable
		Uptime:     time.Since(s.startTime).String(),
	}
}

// GetSystemMetrics retrieves current system metrics
func (s *MonitorService) GetSystemMetrics() *model.SystemMetrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Calculate memory usage percentage
	totalMem := m.Sys
	usedMem := m.Alloc
	memUsage := float64(usedMem) / float64(totalMem) * 100

	// For demonstration purposes, we'll use dummy values for CPU and disk usage
	// In a real application, you would use system calls to get actual metrics
	cpuUsage := 25.5
	diskUsage := 45.2
	networkIn := 1024.0
	networkOut := 2048.0

	return &model.SystemMetrics{
		Timestamp:       time.Now(),
		CPUUsage:        cpuUsage,
		MemoryUsage:     memUsage,
		DiskUsage:       diskUsage,
		NetworkInbound:  networkIn,
		NetworkOutbound: networkOut,
		RequestCount:    1000, // This should come from actual request counting
		ErrorCount:      5,    // This should come from actual error counting
		CreatedAt:       time.Now(),
	}
}

// GetRecentMetrics retrieves recent system metrics for charting
func (s *MonitorService) GetRecentMetrics(count int) []model.SystemMetrics {
	// For demonstration purposes, we'll generate dummy data
	// In a real application, you would retrieve this from a database
	metrics := make([]model.SystemMetrics, count)
	now := time.Now()

	for i := 0; i < count; i++ {
		timestamp := now.Add(-time.Duration(count-i) * time.Minute)

		// Generate dummy values with some variation
		cpuUsage := 20.0 + float64(i%10)
		memUsage := 30.0 + float64((i*2)%20)
		diskUsage := 40.0 + float64(i%15)

		metrics[i] = model.SystemMetrics{
			ID:              uint(i + 1),
			Timestamp:       timestamp,
			CPUUsage:        cpuUsage,
			MemoryUsage:     memUsage,
			DiskUsage:       diskUsage,
			NetworkInbound:  1000.0 + float64(i*100),
			NetworkOutbound: 2000.0 + float64(i*150),
			RequestCount:    int64(100 + i*10),
			ErrorCount:      int64(i % 5),
			CreatedAt:       timestamp,
		}
	}

	return metrics
}
