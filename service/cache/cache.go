package cache

import (
	"crypto"
	"fmt"
	"github.com/slimjim777/snap-downloader/domain"
	"github.com/slimjim777/snap-downloader/service/datastore"
	"github.com/snapcore/snapd/asserts"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
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
	ListDownloads() ([]domain.SnapDownload, error)
	DownloadPath(name, filename string) string

	Delta(name, arch string, fromRevision, toRevision int) (string, error)
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
	log.Printf("Download snap %s (%s)\n", download.Name, download.Arch)
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

// ListDownloads returns the snap downloads
func (c *Cache) ListDownloads() (files []domain.SnapDownload, err error) {
	if _, err := os.Stat(c.baseDir); err != nil {
		return []domain.SnapDownload{}, nil
	}

	// get the snaps in the directory
	err = filepath.Walk(c.baseDir, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(p) != ".snap" {
			return nil
		}

		download := domain.SnapDownload{
			Size:      info.Size(),
			Filename:  info.Name(),
			Assertion: strings.TrimSuffix(info.Name(), ".snap") + ".assert",
		}
		if err := snapFromFilename(&download); err != nil {
			return err
		}

		files = append(files, download)
		return nil
	})
	return files, nil
}

// DownloadPath returns the download path of the file
func (c *Cache) DownloadPath(name, filename string) string {
	return path.Join(c.baseDir, name, filename)
}

func snapFromFilename(download *domain.SnapDownload) error {
	name := strings.TrimSuffix(download.Filename, ".snap")

	parts := strings.SplitN(name, "_", 3)
	if len(parts) != 3 {
		return fmt.Errorf("filename is not in the expected format: %s", download.Filename)
	}
	revision, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("filename is not in the expected format: %s", download.Filename)
	}

	download.Name = parts[0]
	download.Revision = revision
	download.Arch = parts[2]
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
