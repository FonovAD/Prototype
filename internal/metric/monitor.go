package metric

type Monitor interface {
	IncRequestsTotal(string, int)
	IncRequestLatency(string, string, float64)
	IncErrorCount(string, string, int)
}

type TestMonitor struct{}

func NewTest() Monitor {
	return &TestMonitor{}
}

func (m *TestMonitor) IncRequestsTotal(method string, status int) {
}

func (m *TestMonitor) IncRequestLatency(method, path string, duration float64) {
}

func (m *TestMonitor) IncErrorCount(method, path string, status int) {
}
