package metrics

type ContainerMetricsProvider interface {
	GetContainersStatus() ([]ContainerStatus, error)
}

type ContainerStatus struct {
	ContainerID  string
	Name         string
	Health       int64
	MemoryUsage  float64 // in percentage of the host memory
	CpuUsage     float64 // in percentage of the host CPU
	Image        string
	Uptime       int64 // in seconds
	RestartCount int64
	State        int64
}
