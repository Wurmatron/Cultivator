package model

type ComputedNodeAvg struct {
	Name       string  `json:"name"`
	PlotCount  int64   `json:"plotCount"`
	FreeMemory int64   `json:"freeMemory"`
	CpuUsage   int64   `json:"cpuUsage"`
	DbSize     int64   `json:"dbSize"`
	Uptime     float64 `json:"uptime"`
}

type ComputedHarvesterAvg struct {
	ID              string
	Blockchain      string
	PlotCount       int64
	Drives          []ComputedDriveAvg
	PlotHarvestTime float64
}

type ComputedDriveAvg struct {
	Serial     string
	TotalSpace int64
	FreeSpace  int64
	Smart      ComputedSmartAvg
	PlotCount  int64
}

type ComputedSmartAvg struct {
	MaxTemp      float64
	MinTemp      float64
	CurrentTemp  float64
	PowerOnHours int64
	CycleCount   int64
}

type ComputedPowerAvg struct {
	MaxPowerDraw   int64
	Wattage        float64
	CurrentLoad    int64
	InputVoltage   float64
	BatteryVoltage float64
}
