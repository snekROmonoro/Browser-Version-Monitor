package monitor

import "fmt"

// Microsoft Edge
// https://edgeupdates.microsoft.com/api/products

func EdgeMonitor() (*MonitorResult, error) {
	type ProductRelease struct {
		ReleaseID          int        `json:"ReleaseId"`
		Platform           string     `json:"Platform"`
		Architecture       string     `json:"Architecture"`
		CVEs               []any      `json:"CVEs"`
		ProductVersion     string     `json:"ProductVersion"`
		Artifacts          []any      `json:"Artifacts"`
		PublishedTime      CustomTime `json:"PublishedTime"`
		ExpectedExpiryDate CustomTime `json:"ExpectedExpiryDate"`
	}

	type ResponseJSON []struct {
		Product  string           `json:"Product"`
		Releases []ProductRelease `json:"Releases"`
	}

	resp, err := httpRequestDo(&httpRequest{
		Method: "GET",
		URL:    "https://edgeupdates.microsoft.com/api/products",
	})

	if err != nil {
		return nil, err
	}

	var body ResponseJSON
	if err := resp.DecodeJSON(&body); err != nil {
		return nil, err
	}

	for i := 0; i < len(body); i++ {
		if body[i].Product == "Stable" {
			for j := 0; j < len(body[i].Releases); j++ {
				if body[i].Releases[j].Platform == "Windows" {
					return &MonitorResult{
						Browser:    "Microsoft Edge",
						Version:    body[i].Releases[j].ProductVersion,
						Platform:   "Windows",
						StableDate: body[i].Releases[j].PublishedTime.Time, // it's always going to be same day
					}, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("stable channel not found for edge")
}
