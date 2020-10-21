package cache

import (
	"fmt"
	"github.com/slimjim777/snap-downloader/domain"
	"github.com/slimjim777/snap-downloader/service/datastore"
	"github.com/slimjim777/snap-downloader/service/store"
	"github.com/snapcore/snapd/asserts"
	"log"
	"time"
)

const (
	//tickInterval = 300
	tickInterval = 5
)

// WatchService interface for watching snaps
type WatchService interface {
	Watch()
}

// Watch implements a build service
type Watch struct {
	data  datastore.Datastore
	store store.Service
	cache Service
}

// NewWatchService creates a new watch service
func NewWatchService(ds datastore.Datastore, store store.Service, cache Service) *Watch {
	return &Watch{
		data:  ds,
		store: store,
		cache: cache,
	}
}

// Watch service to watch for snap updates
func (srv *Watch) Watch() {
	// on an interval...
	ticker := time.NewTicker(time.Second * tickInterval)
	for range ticker.C {
		// get the snap list
		records, err := srv.data.SnapsList()
		if err != nil {
			log.Println("Error fetching snap list:", err)
			break
		}

		for _, r := range records {
			log.Printf("Check snap: %s (%s)", r.Name, r.Arch)
			download, updateFound, err := srv.checkForUpdates(r)
			if err != nil {
				log.Println("Error checking snap:", err)
				// check the next repo
				continue
			}
			if !updateFound {
				// no update so check the next repo
				continue
			}

			if err := srv.downloadSnap(download); err != nil {
				log.Println("Error checking snap:", err)
			}
		}
	}
	ticker.Stop()
}

func (srv *Watch) checkForUpdates(r domain.SnapCache) (*domain.SnapDownload, bool, error) {
	// get info on the snap
	info, err := srv.store.SnapInfo(r.Name, r.Arch)
	if err != nil {
		return nil, false, err
	}

	// check we have a stable release for the snap
	download := getStableRelease(info, r.Arch)
	if download.URL == "" {
		return nil, false, fmt.Errorf("no stable download found")
	}

	// check if we have this snap revision downloaded
	found := srv.cache.CheckDownloadForSnap(r.Name, download.Filename)

	return &download, !found, nil
}

func (srv *Watch) downloadSnap(download *domain.SnapDownload) error {
	// get the snap download stream
	resp, err := srv.store.GetSnapStream(download.URL)
	if err != nil {
		return err
	}

	// download the snap to the cache using the URL
	if err := srv.cache.DownloadSnap(resp, download); err != nil {
		return err
	}

	// get the snap assertions
	assertRev, err := srv.store.Assertion("snap-revision", download.AssertionKey)
	if err != nil {
		return err
	}
	assertDecl, err := srv.store.Assertion("snap-declaration", fmt.Sprintf("16/%s", assertRev.HeaderString("snap-id")))
	if err != nil {
		return err
	}
	assertAcct, err := srv.store.Assertion("account", assertRev.HeaderString("developer-id"))
	if err != nil {
		return err
	}
	assertAcctKey, err := srv.store.Assertion("account-key", assertRev.HeaderString("sign-key-sha3-384"))
	if err != nil {
		return err
	}

	// store the assertion
	if err := srv.cache.SnapAssertion([]asserts.Assertion{assertAcct, assertAcctKey, assertDecl, assertRev}, download); err != nil {
		return err
	}

	return nil
}

func getStableRelease(info *store.ResponseSnapInfo, arch string) domain.SnapDownload {
	download := domain.SnapDownload{
		Name: info.Name,
		Arch: arch,
	}

	for _, m := range info.ChannelMap {
		if m.Channel.Risk != "stable" {
			continue
		}
		download.URL = m.Download.URL
		download.Size = m.Download.Size
		download.Sha3_384 = m.Download.Sha3_384
		download.Revision = m.Revision
		download.Filename = fmt.Sprintf("%s_%d_%s.snap", info.Name, m.Revision, arch)
	}
	return download
}
