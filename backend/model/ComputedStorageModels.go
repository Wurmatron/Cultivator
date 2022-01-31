package model

type ComputedNodeAvg struct {
	Name       string  `json:"name"`
	PlotCount  int64   `json:"plotCount"`
	FreeMemory int64   `json:"freeMemory"`
	CpuUsage   int64   `json:"cpuUsage"`
	DbSize     int64   `json:"dbSize"`
	Uptime     float64 `json:"uptime"`
}
