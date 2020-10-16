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
