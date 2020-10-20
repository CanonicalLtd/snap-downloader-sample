package datastore

import "github.com/slimjim777/snap-downloader/domain"

// Datastore interface for the database logic
type Datastore interface {
	SettingsPut(key, name, data string) (string, error)
	SettingsGet(key, name string) (domain.ConfigSetting, error)

	SnapsCreate(name, arch string) (string, error)
	SnapsDelete(id string) error
	SnapsList() ([]domain.SnapCache, error)
}
