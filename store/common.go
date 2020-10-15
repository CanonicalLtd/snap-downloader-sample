package store

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

const (
	apiBaseURL     = "https://api.snapcraft.io/"
	configFilename = "config.json"
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

func cacheHeaders(headers map[string]string) error {
	data, err := json.Marshal(headers)
	if err != nil {
		return err
	}

	// cache the headers in SNAP_DATA
	return ioutil.WriteFile(getConfigPath(), data, 0600)
}

func getConfigPath() string {
	if os.Getenv("SNAP_DATA") == "" {
		return configFilename
	}
	return path.Join(os.Getenv("SNAP_DATA"), configFilename)
}

func readHeaders() (map[string]string, error) {
	var headers map[string]string
	data, err := ioutil.ReadFile(getConfigPath())
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, headers); err != nil {
		return nil, err
	}

	return headers, nil
}
