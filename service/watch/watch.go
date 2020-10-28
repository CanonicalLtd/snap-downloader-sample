package watch

import (
	"fmt"
	"github.com/slimjim777/snap-downloader/domain"
	cache2 "github.com/slimjim777/snap-downloader/service/cache"
	"github.com/slimjim777/snap-downloader/service/datastore"
	"github.com/slimjim777/snap-downloader/service/store"
	"log"
	"time"
)

const (
	tickInterval = 300
	minInterval  = 5
)

// Service interface for watching snaps
type Service interface {
	GetInterval() int
	SetInterval(newInterval int) error
	LastRun() string
	Watch()
}

// Watch implements a build service
type Watch struct {
	data  datastore.Datastore
	store store.Service
	cache cache2.Service
}

// NewWatchService creates a new watch service
func NewWatchService(ds datastore.Datastore, store store.Service, cache cache2.Service) *Watch {
	return &Watch{
		data:  ds,
		store: store,
		cache: cache,
	}
}

// Watch service to watch for snap updates
func (srv *Watch) Watch() {
	// on an interval...
	interval := srv.watchInterval()
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ticker.C:
			// run the tasks for this cycle
			srv.runCycle()

			// check if we need to adjust the ticker interval
			currentInterval := srv.watchInterval()
			if interval != currentInterval {
				interval = currentInterval
				ticker.Stop()
				ticker = time.NewTicker(interval)
			}
		}
	}
	ticker.Stop()
}

func (srv *Watch) watchInterval() time.Duration {
	// get the setting from the data store
	interval := srv.GetInterval()
	return time.Duration(interval) * time.Second
}

func (srv *Watch) runCycle() {
	// set the last run
	_ = srv.setLastRun()

	// get the snap list
	records, err := srv.data.SnapsList()
	if err != nil {
		log.Println("Error fetching snap list:", err)
		return
	}

	// refresh the store headers, as a login may have happened
	_ = srv.store.GetHeaders()

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
	assertions, err := srv.store.SnapAssertions(download)
	if err != nil {
		return err
	}

	// store the assertion
	return srv.cache.SnapAssertion(assertions, download)
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
