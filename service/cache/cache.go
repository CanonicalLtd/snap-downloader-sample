package cache

import (
	"github.com/slimjim777/snap-downloader/domain"
	"github.com/slimjim777/snap-downloader/service/datastore"
)

// Service is the abstract cache service
type Service interface {
	SnapAdd(name, arch string) (string, error)
	SnapList() ([]domain.SnapCache, error)
	SnapDelete(id string) error
}

// Cache service for store snaps
type Cache struct {
	baseDir string
	data    datastore.Datastore
}

// NewService returns a snap cache service
func NewService(downloadPath string, ds datastore.Datastore) *Cache {
	return &Cache{
		baseDir: downloadPath,
		data:    ds,
	}
}

// SnapAdd add a snap to the watch list
func (c *Cache) SnapAdd(name, arch string) (string, error) {
	return c.data.SnapsCreate(name, arch)
}

// SnapList lists the snaps in the watch list
func (c *Cache) SnapList() ([]domain.SnapCache, error) {
	return c.data.SnapsList()
}

// SnapDelete removes a snap from the watch list
func (c *Cache) SnapDelete(id string) error {
	return c.data.SnapsDelete(id)
}
