package domain

import "time"

// ConfigSetting is a stored config setting
type ConfigSetting struct {
	ID       string    `json:"id"`
	Key      string    `json:"key"`
	Name     string    `json:"name"`
	Data     string    `json:"data"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

// SnapCache is a snap to cache
type SnapCache struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Arch     string    `json:"arch"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

// SnapDownload defines the details of a snap download file
type SnapDownload struct {
	Name         string `json:"name"`
	Arch         string `json:"arch"`
	Revision     int    `json:"revision"`
	URL          string `json:"url,omitempty"`
	Assertion    string `json:"assertion"`
	Sha3_384     string `json:"sha3-384,omitempty"`
	AssertionKey string `json:"key,omitempty"`
	Size         int64  `json:"size"`
	Filename     string `json:"filename"`
}
