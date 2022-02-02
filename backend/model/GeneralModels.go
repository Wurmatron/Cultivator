package model

type Node struct {
	Name                  string `json:"name"`
	LastSync              int64  `json:"lastSync"`
	Version               string `json:"version"`
	NodeStatus            string `json:"nodeStatus"`
	PlotCount             int64  `json:"plotCount"`
	PlotSize              string `json:"plotSize"`
	EstimatedTimeToWin    string `json:"estimatedTimeToWin"`
	EstimatedNetworkSpace string `json:"estimatedNetworkSpace"`
	FreeMemory            int64  `json:"freeMemory"`
	TotalMemory           int64  `json:"totalMemory"`
	CpuUsage              int64  `json:"cpuUsage"`
	DbSize                int64  `json:"dbSize"`
	Wallet                string `json:"wallet"`
}

type Harvester struct {
	ID               string  `json:"id"`
	Blockchain       string  `json:"blockchain"`
	LastSync         int64   `json:"lastSync"`
	Version          string  `json:"version"`
	ConnectionStatus string  `json:"connectionStatus"`
	PlotCount        int64   `json:"plotCount"`
	PlotSize         string  `json:"plotSize"`
	Drives           []Drive `json:"drives"`
	PlotHarvestTime  float64 `json:"plotHarvestTime"`
}

type Drive struct {
	Serial     string     `json:"serial"`
	Mount      string     `json:"mount"`
	TotalSpace int64      `json:"totalSpace"`
	FreeSpace  int64      `json:"freeSpace"`
	Smart      SmartDrive `json:"smart"`
	Plots      []string   `json:"plots"`
	PlotCount  int64      `json:"plotCount"`
}

type SmartDrive struct {
	Model            string  `json:"model"`
	Serial           string  `json:"serial"`
	MaxTemp          float64 `json:"maxTemp"`
	MinTemp          float64 `json:"minTemp"`
	CurrentTemp      float64 `json:"currentTemp"`
	PowerOnHours     int64   `json:"powerOnHours"`
	CycleCount       int64   `json:"cycleCount"`
	RecommendMinTemp float64 `json:"recommendMinTemp"`
	RecommendMaxTemp float64 `json:"recommendMaxTemp"`
}

type HistoryPlots struct {
	Last1h    int64 `json:"last1h"`
	Last6h    int64 `json:"last6h"`
	Last12h   int64 `json:"last12h"`
	Last24h   int64 `json:"lsat24h"`
	Last7d    int64 `json:"last7d"`
	LastMonth int64 `json:"lastMonth"`
}

type HistoryPlotting struct {
	AvgTimePerPlot string `json:"avgTimePerPlot"`
	PlotsToday     int64  `json:"plotsToday"`
	PlotsYesterday int64  `json:"plotsYesterday"`
	PlotsWeek      int64  `json:"plotsWeek"`
	PlotsMonth     int64  `json:"plotsMonth"`
}

type HistoryNode struct {
	Name        string  `json:"name"`
	Uptime      string  `json:"uptime"`
	MemoryUsage float64 `json:"memoryUsage"`
	CpuUsage    float64 `json:"cpuUsage"`
}

type HistoryDrives struct {
	AvgTemp        string `json:"avgTemp"`
	FreeSpace      string `json:"freeSpace"`
	AvgUtilization string `json:"avgUtilization"`
}

type HistoryPower struct {
	Avg1h  float64 `json:"avg1h"`
	Avg3h  float64 `json:"avg3h"`
	Avg6h  float64 `json:"avg6h"`
	Avg12h float64 `json:"avg12h"`
	Avg24h float64 `json:"avg24h"`
	Avg1w  float64 `json:"avg1w"`
}

type History struct {
	Plots    HistoryPlots    `json:"plots"`
	Plotting HistoryPlotting `json:"plotting"`
	Node     HistoryNode     `json:"node"`
	Drives   HistoryDrives   `json:"drives"`
	Power    HistoryPower    `json:"power"`
}

type Power struct {
	MaxPowerDraw          int64   `json:"maxPowerDraw"`
	CurrentLoad           int64   `json:"currentLoad"`
	InputVoltage          float64 `json:"inputVoltage"`
	OutputVoltage         float64 `json:"outputVoltage"`
	InputVoltageNominal   float64 `json:"inputVoltageNominal"`
	Wattage               float64 `json:"wattage"`
	Model                 string  `json:"model"`
	Serial                string  `json:"serial"`
	BatteryVoltage        float64 `json:"batteryVoltage"`
	BatteryVoltageNominal float64 `json:"BatteryVoltageNominal"`
	BatteryType           string  `json:"batteryType"`
	BatteryCharge         int64   `json:"batteryCharge"`
	BatteryRuntime        int64   `json:"batteryRuntime"`
	Generated             int64   `json:"generated"`
}
