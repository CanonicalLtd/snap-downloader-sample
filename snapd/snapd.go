package snapd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/snapcore/snapd/overlord/auth"
	"io"
	"io/ioutil"
	"net"
	"net/http"
)

const (
	socketFile     = "/run/snapd.socket"
	urlAssertions  = "/v2/assertions"
	urlSnaps       = "/v2/snaps"
	urlLogin       = "/v2/login"
	urlDownload    = "/v2/download"
	typeAssertions = "application/x.ubuntu.assertion"
	typeJSON       = "application/json"
	baseURL        = "http://localhost"
)

// Client is the abstract client interface
type Client interface {
	Login(email, password, otp string) (*auth.UserState, error)
}

// Snapd service to access the snapd REST API
type Snapd struct {
	client *http.Client
}

// NewClient returns a snapd API client
func NewClient(downloadPath string) *Snapd {
	return &Snapd{
		client: &http.Client{
			Transport: &http.Transport{
				DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("unix", socketFile)
				},
			},
		},
	}
}

func (snap *Snapd) call(method, url, contentType string, body io.Reader) (*http.Response, error) {
	u := baseURL + url

	switch method {
	case "POST":
		return snap.client.Post(u, contentType, body)
	case "GET":
		return snap.client.Get(u)
	default:
		return nil, fmt.Errorf("unsupported method: %s", method)
	}
}

// Login authenticates with the snap store via snapd
func (snap *Snapd) Login(email, password, otp string) (*auth.UserState, error) {
	params := map[string]string{
		"email": email, "password": password, "otp": otp,
	}
	data, err := json.Marshal(&params)
	if err != nil {
		return nil, err
	}

	resp, err := snap.call("POST", urlLogin, typeJSON, bytes.NewReader(data))

	defer resp.Body.Close()

	var authUser map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&authUser)

	user := authUser["result"].(auth.UserState)

	return &user, err
}

// Download fetches a snap and its assertion
func (snap *Snapd) Download(name string) error {
	params := map[string]string{
		"snap-name": name,
	}
	data, err := json.Marshal(&params)
	if err != nil {
		return err
	}

	resp, err := snap.call("POST", urlDownload, typeJSON, bytes.NewReader(data))

	defer resp.Body.Close()
	bb, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bb))
	return err
}
