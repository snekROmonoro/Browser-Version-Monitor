package monitor

// https://product-details.mozilla.org/1.0/firefox_versions.json

func FirefoxMonitor() (*MonitorResult, error) {
	type ResponseJSON struct {
		FirefoxAurora                     string `json:"FIREFOX_AURORA"`
		FirefoxDevedition                 string `json:"FIREFOX_DEVEDITION"`
		FirefoxEsr                        string `json:"FIREFOX_ESR"`
		FirefoxEsrNext                    string `json:"FIREFOX_ESR_NEXT"`
		FirefoxNightly                    string `json:"FIREFOX_NIGHTLY"`
		LastMergeDate                     string `json:"LAST_MERGE_DATE"`
		LastReleaseDate                   string `json:"LAST_RELEASE_DATE"`
		LastSoftfreezeDate                string `json:"LAST_SOFTFREEZE_DATE"`
		LastStringfreezeDate              string `json:"LAST_STRINGFREEZE_DATE"`
		LatestFirefoxDevelVersion         string `json:"LATEST_FIREFOX_DEVEL_VERSION"`
		LatestFirefoxOlderVersion         string `json:"LATEST_FIREFOX_OLDER_VERSION"`
		LatestFirefoxReleasedDevelVersion string `json:"LATEST_FIREFOX_RELEASED_DEVEL_VERSION"`
		LatestFirefoxVersion              string `json:"LATEST_FIREFOX_VERSION"`
		NextMergeDate                     string `json:"NEXT_MERGE_DATE"`
		NextReleaseDate                   string `json:"NEXT_RELEASE_DATE"`
		NextSoftfreezeDate                string `json:"NEXT_SOFTFREEZE_DATE"`
		NextStringfreezeDate              string `json:"NEXT_STRINGFREEZE_DATE"`
	}

	resp, err := httpRequestDo(&httpRequest{
		Method: "GET",
		URL:    "https://product-details.mozilla.org/1.0/firefox_versions.json",
	})

	if err != nil {
		return nil, err
	}

	var body ResponseJSON
	if err := resp.DecodeJSON(&body); err != nil {
		return nil, err
	}

	return &MonitorResult{
		Browser:  "Firefox",
		Version:  body.LatestFirefoxVersion,
		Platform: "",
	}, nil
}
