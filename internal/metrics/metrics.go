package metrics

type ContainerMetricsProvider interface {
	GetContainersStatus() ([]ContainerStatus, error)
}

type ContainerStatus struct {
	ContainerID  string
	Name         string
	Health       int64
	MemoryUsage  float64
	CpuUsage     float64
	Image        string
	Uptime       string
	RestartCount int64
	State        int64
}
