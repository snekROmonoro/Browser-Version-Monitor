package monitor

import (
	"fmt"
)

// Google Chrome: https://chromestatus.com/api/v0/channels

func ChromeMonitor() (*MonitorResult, error) {
	type ResponseJSON map[string]struct {
		BranchPoint CustomTime `json:"branch_point"`
		// LateStableDate CustomTime `json:"late_stable_date"`
		StableDate CustomTime `json:"stable_date"`
		Version    int        `json:"version"`
	}

	resp, err := httpRequestDo(&httpRequest{
		Method: "GET",
		URL:    "https://chromestatus.com/api/v0/channels",
	})

	if err != nil {
		return nil, err
	}

	// google has anti-bot shit (or something idk), remove everything from r.Body until a '{' is found
	for i := 0; i < len(resp.Body); i++ {
		if resp.Body[i] == '{' {
			resp.Body = resp.Body[i:]
			break
		}
	}

	var body ResponseJSON
	if err := resp.DecodeJSON(&body); err != nil {
		return nil, err
	}

	stable, ok := body["stable"]
	if !ok {
		return nil, fmt.Errorf("stable channel not found")
	}

	return &MonitorResult{
		Browser:    "Google Chrome",
		Version:    fmt.Sprintf("%d", stable.Version),
		Platform:   "", // any
		StableDate: stable.StableDate.Time,
	}, nil
}
