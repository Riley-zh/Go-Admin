package model

import (
	"time"
)

// SystemMetrics represents system metrics data
type SystemMetrics struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Timestamp       time.Time `json:"timestamp"`
	CPUUsage        float64   `json:"cpu_usage"`
	MemoryUsage     float64   `json:"memory_usage"`
	DiskUsage       float64   `json:"disk_usage"`
	NetworkInbound  float64   `json:"network_inbound"`
	NetworkOutbound float64   `json:"network_outbound"`
	RequestCount    int64     `json:"request_count"`
	ErrorCount      int64     `json:"error_count"`
	CreatedAt       time.Time `json:"created_at"`
}

// SystemInfo represents system information
type SystemInfo struct {
	OS         string `json:"os"`
	Arch       string `json:"arch"`
	GoVersion  string `json:"go_version"`
	AppVersion string `json:"app_version"`
	Uptime     string `json:"uptime"`
}
