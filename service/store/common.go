package store

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const (
	apiBaseURL = "https://api.snapcraft.io/v2"
)

func submitPOSTRequest(url string, headers map[string]string, data []byte) (*http.Response, error) {
	r, _ := http.NewRequest("POST", url, bytes.NewReader(data))

	return requestDo(r, headers)
}

func submitGETRequest(url string, headers map[string]string) (*http.Response, error) {
	r, _ := http.NewRequest("GET", url, nil)

	return requestDo(r, headers)
}

func requestDo(r *http.Request, headers map[string]string) (*http.Response, error) {
	for k, v := range headers {
		r.Header.Set(k, v)
	}
	client := http.Client{}
	return client.Do(r)
}

func readHeaders(data []byte) (map[string]string, error) {
	var headers map[string]string

	if err := json.Unmarshal(data, &headers); err != nil {
		return nil, err
	}

	return headers, nil
}
