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
