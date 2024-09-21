package monitor

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
)

type MonitorResult struct {
	Browser    string
	Version    string
	Platform   string // can be empty
	StableDate time.Time
}

func (r *MonitorResult) UpdateString(prevVersion *version.Version) string {
	currVersion, err := version.NewVersion(r.Version)
	if err != nil {
		return ""
	}

	var majorUpdate bool = false
	if prevVersion != nil && currVersion.GreaterThan(prevVersion) {
		majorUpdate = currVersion.Segments()[0] > prevVersion.Segments()[0]
	}

	var ret strings.Builder
	ret.WriteString(fmt.Sprintf("Browser `%s` got an update\n", r.Browser))

	ret.WriteString("Version: ")
	if prevVersion != nil {
		ret.WriteString(fmt.Sprintf("`%s` -> ", prevVersion.Original()))
	}
	ret.WriteString(fmt.Sprintf("`%s`\n", currVersion.Original()))

	ret.WriteString(fmt.Sprintf("Major: `%v`\n", majorUpdate))

	if r.Platform != "" {
		ret.WriteString(fmt.Sprintf("Platform: `%s`\n", r.Platform))
	}

	if !r.StableDate.IsZero() {
		ret.WriteString(fmt.Sprintf("Stable Date: `%s`\n", r.StableDate.Format("2006-01-02")))
	}

	return strings.TrimSpace(ret.String())
}

type MonitorFunc func() (*MonitorResult, error)

// An array of monitor functions, to make our life easier
var MonitorFuncs = []MonitorFunc{
	ChromeMonitor,
	FirefoxMonitor,
	ChromiumMonitor,
	EdgeMonitor,
	BraveMonitor,
}
