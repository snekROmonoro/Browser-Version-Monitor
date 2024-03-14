package monitor

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	urlpkg "net/url"
)

type httpRequest struct {
	Method    string
	URL       string
	Headers   map[string]string
	URLParams map[string]string
	Body      interface{}
}

type httpResponse struct {
	// Status code in string format (ex: 200 OK)
	Status string

	// Status code
	StatusCode int

	// The body of the response
	Body []byte
}

func (r *httpResponse) DecodeJSON(v interface{}) error {
	return json.Unmarshal(r.Body, v)
}

func httpRequestDo(req *httpRequest) (*httpResponse, error) {
	var bodyReader io.Reader = nil
	if req.Body != nil {
		b, _ := json.Marshal(req.Body)
		bodyReader = bytes.NewReader(b)
	}

	parsedUrl, err := urlpkg.Parse(req.URL)
	if err != nil {
		return nil, err
	}

	// append url params to the existing ones
	values := parsedUrl.Query()
	for key, val := range req.URLParams {
		values.Add(key, val)
	}

	parsedUrl.RawQuery = values.Encode()

	httpReq, err := http.NewRequest(req.Method, parsedUrl.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	if bodyReader != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	return &httpResponse{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Body:       body,
	}, nil
}
