package monitor

import (
	"fmt"
	"github.com/hashicorp/go-version"
	"regexp"
	"strings"
	"time"
)

// Credits: https://github.com/brave/brave-versions
// https://api.github.com/repos/brave/brave-browser/releases

func BraveMonitor() (*MonitorResult, error) {
	type ResponseJSON []struct {
		Name        string    `json:"name"`
		Draft       bool      `json:"draft"`
		PreRelease  bool      `json:"prerelease"`
		CreatedAt   time.Time `json:"created_at"`
		PublishedAt time.Time `json:"published_at"`
	}

	resp, err := httpRequestDo(&httpRequest{
		Method: "GET",
		URL:    "https://api.github.com/repos/brave/brave-browser/releases",
	})

	if err != nil {
		return nil, err
	}

	var body ResponseJSON
	if err := resp.DecodeJSON(&body); err != nil {
		return nil, err
	}

	var bestMatch *MonitorResult = nil
	var bestMatchVersion *version.Version = nil

	for _, release := range body {
		if release.Draft || release.PreRelease {
			continue
		}

		release.Name = strings.TrimSpace(release.Name)

		if !strings.HasPrefix(release.Name, "Release v") {
			continue
		}

		var re = regexp.MustCompile(`(?m)\(Chromium ([0-9.]+)\)`)
		var chromiumVersion = re.FindStringSubmatch(release.Name)

		if len(chromiumVersion) < 2 {
			continue
		}

		parsedVersion, err := version.NewVersion(chromiumVersion[1])
		if err != nil {
			continue
		}

		if bestMatch == nil || parsedVersion.GreaterThan(bestMatchVersion) {
			bestMatch = &MonitorResult{
				Browser: "Brave",
				Version: parsedVersion.Original(),
			}

			bestMatchVersion = parsedVersion
		}
	}

	if bestMatch != nil {
		return bestMatch, nil
	}

	return nil, fmt.Errorf("brave release data not found")
}
