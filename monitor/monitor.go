package monitor

type MonitorResult struct {
	Browser  string
	Version  string
	Platform string // can be empty
}

type MonitorFunc func() (*MonitorResult, error)

// An array of monitor functions, to make our life easier
var MonitorFuncs = []MonitorFunc{
	FirefoxMonitor,
	ChromiumMonitor,
}
