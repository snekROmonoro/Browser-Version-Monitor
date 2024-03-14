package monitor

import "fmt"

// Chromium: https://chromiumdash.appspot.com/fetch_releases?channel=Stable&platform=Windows&num=1
// still Chromium: https://versionhistory.googleapis.com/v1/chrome/platforms/win/channels/stable/versions

func ChromiumMonitor() (*MonitorResult, error) {
	type ResponseJSON []struct {
		Channel                    string `json:"channel"`
		ChromiumMainBranchPosition int    `json:"chromium_main_branch_position"`
		Hashes                     struct {
			Angle    string `json:"angle"`
			Chromium string `json:"chromium"`
			Dawn     string `json:"dawn"`
			Devtools string `json:"devtools"`
			Pdfium   string `json:"pdfium"`
			Skia     string `json:"skia"`
			V8       string `json:"v8"`
			Webrtc   string `json:"webrtc"`
		} `json:"hashes"`
		Milestone       int    `json:"milestone"`
		Platform        string `json:"platform"`
		PreviousVersion string `json:"previous_version"`
		Time            int64  `json:"time"`
		Version         string `json:"version"`
	}

	resp, err := httpRequestDo(&httpRequest{
		Method: "GET",
		URL:    "https://chromiumdash.appspot.com/fetch_releases?channel=Stable&platform=Windows&num=1",
	})

	if err != nil {
		return nil, err
	}

	var body ResponseJSON
	if err := resp.DecodeJSON(&body); err != nil {
		return nil, err
	}

	if len(body) != 1 {
		return nil, fmt.Errorf("expected 1 body, got %d", len(body))
	}

	return &MonitorResult{
		Browser:  "Chromium",
		Version:  body[0].Version,
		Platform: body[0].Platform,
	}, nil
}
