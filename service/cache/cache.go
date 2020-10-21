package cache

import (
	"crypto"
	"fmt"
	"github.com/slimjim777/snap-downloader/domain"
	"github.com/slimjim777/snap-downloader/service/datastore"
	"github.com/snapcore/snapd/asserts"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

// Service is the abstract cache service
type Service interface {
	SnapAdd(name, arch string) (string, error)
	SnapList() ([]domain.SnapCache, error)
	SnapDelete(id string) error
	CheckDownloadForSnap(snap, filename string) bool
	DownloadSnap(resp *http.Response, download *domain.SnapDownload) error
	SnapAssertion(data []asserts.Assertion, download *domain.SnapDownload) error
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

// CheckDownloadForSnap checks if we have a specific snap download
func (c *Cache) CheckDownloadForSnap(snap, filename string) bool {
	filepath := path.Join(c.baseDir, snap, filename)
	_, err := os.Stat(filepath)

	return err == nil
}

// DownloadSnap performs the snap download stream
func (c *Cache) DownloadSnap(resp *http.Response, download *domain.SnapDownload) error {
	// create the download path
	if err := os.MkdirAll(path.Join(c.baseDir, download.Name), 0755); err != nil {
		return err
	}

	// download the file
	filepath := path.Join(c.baseDir, download.Name, download.Filename)
	if err := downloadFromURL(resp, filepath); err != nil {
		return err
	}

	// check the size and hash of the download
	if err := verifyDownloadSetHash(download, filepath); err != nil {
		return err
	}

	return nil
}

// SnapAssertion stores the assertions for a snap
func (c *Cache) SnapAssertion(assertions []asserts.Assertion, download *domain.SnapDownload) error {
	filename := strings.TrimSuffix(download.Filename, "snap") + "assert"
	filepath := path.Join(c.baseDir, download.Name, filename)

	w, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer w.Close()

	enc := asserts.NewEncoder(w)
	for _, a := range assertions {
		if err := enc.Encode(a); err != nil {
			return err
		}
	}

	return nil
}

func downloadFromURL(resp *http.Response, filepath string) error {
	// create the .part file
	filepathPart := filepath + ".part"
	out, err := os.Create(filepathPart)
	if err != nil {
		return err
	}
	defer out.Close()
	defer resp.Body.Close()

	// stream the file download
	if _, err := io.Copy(out, resp.Body); err != nil {
		return err
	}

	// rename the .part file to the real filename
	err = os.Rename(filepathPart, filepath)
	return err
}

func verifyDownloadSetHash(download *domain.SnapDownload, filepath string) error {
	// check the file size
	info, err := os.Stat(filepath)
	if err != nil {
		return err
	}
	if info.Size() != download.Size {
		return fmt.Errorf("unexpected file size, got %v, expected: %v", info.Size(), download.Size)
	}

	// check the hash against the downloaded file
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	h := crypto.SHA3_384.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	actualSha3 := fmt.Sprintf("%x", h.Sum(nil))
	if actualSha3 != download.Sha3_384 {
		return fmt.Errorf("downloaded file hash is not as expected")
	}

	// set the hash in the URL-encoded format
	urlHash, _ := asserts.EncodeDigest(crypto.SHA3_384, h.Sum(nil))
	download.AssertionKey = urlHash

	return nil
}
