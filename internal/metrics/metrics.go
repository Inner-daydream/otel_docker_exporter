package metrics

type ContainerMetricsProvider interface {
	GetContainersStatus() ([]ContainerStatus, error)
}

type ContainerStatus struct {
	ContainerID      string
	Name             string
	Health           int64
	State            int64
	RestartCount     int64
	MemoryUsagePer   float64 // in percentage of the host memory
	MemoryUsageBytes int64
	TotalMemory      int64
	CpuUsage         float64 // in percentage of the available cpu
	Image            string
	Uptime           int64 // in seconds
	AdditionalLabels map[string]string
}
