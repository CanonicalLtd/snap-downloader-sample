package store

import (
	"encoding/json"
	"fmt"
	"github.com/slimjim777/snap-downloader/service/datastore"
	"github.com/snapcore/snapd/asserts"
	"log"
	"net/http"
)

// Service is the interface for the store
type Service interface {
	LoginUser(email, password, otp, storeID, series string) error
	SnapInfo(name, arch string) (*ResponseSnapInfo, error)
	Macaroon() (map[string]string, error)
	GetSnapStream(snapURL string) (*http.Response, error)
	Assertion(assertType, key string) (asserts.Assertion, error)
}

// SnapStore interacts with a brand store
type SnapStore struct {
	Datastore datastore.Datastore
	headers   map[string]string
}

// NewStore creates a new store client
func NewStore(ds datastore.Datastore) *SnapStore {
	// check if we have cached headers (with the store macaroon)
	cfg, err := ds.SettingsGet("store", "headers")
	if err != nil {
		return &SnapStore{
			Datastore: ds,
		}
	}
	headers, err := readHeaders([]byte(cfg.Data))
	if err != nil {
		return &SnapStore{
			Datastore: ds,
		}
	}

	// use the cached headers, so no login needed
	return &SnapStore{
		Datastore: ds,
		headers:   headers,
	}
}

// LoginUser login to the store and request needed ACLs
func (sto *SnapStore) LoginUser(email, password, otp, storeID, series string) error {
	macaroon, discharge, err := LoginUser(email, password, otp, []string{"package_access"})
	if err != nil {
		return err
	}

	authHeader, err := AuthorizationHeader(macaroon, discharge)
	if err != nil {
		return err
	}

	// set the headers to access the brand store
	sto.headers = map[string]string{
		"Snap-Device-Store":   storeID,
		"Snap-Device-Series":  series,
		"Snap-Device-Channel": "stable",
		"Authorization":       authHeader,
		"Content-Type":        "application/json",
		"Accept":              "application/json",
	}

	// cache the headers in the database
	data, err := json.Marshal(sto.headers)
	if err != nil {
		return err
	}
	_, err = sto.Datastore.SettingsPut("store", "headers", string(data))
	return err
}

// Macaroon returns the stored macaroon from the data store
func (sto *SnapStore) Macaroon() (map[string]string, error) {
	cfg, err := sto.Datastore.SettingsGet("store", "headers")
	if err != nil {
		return nil, err
	}

	headers, err := readHeaders([]byte(cfg.Data))
	if err != nil {
		return nil, err
	}

	// remove the actual macaroon from the response
	delete(headers, "Authorization")
	delete(headers, "Content-Type")
	delete(headers, "Accept")

	headers["Created"] = cfg.Created.String()
	headers["Modified"] = cfg.Modified.String()
	return headers, nil
}

// SnapInfo lists the snaps in a brand store
func (sto SnapStore) SnapInfo(name, arch string) (*ResponseSnapInfo, error) {
	headers := map[string]string{"Snap-Device-Architecture": arch}
	for k, v := range sto.headers {
		headers[k] = v
	}

	u := fmt.Sprintf("%s/snaps/info/%s", apiBaseURL, name)
	resp, err := submitGETRequest(u, headers)
	if err != nil {
		log.Printf("Error fetching snap info: %v", err)
		return nil, err
	}

	var response ResponseSnapInfo

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetSnapStream the snap file stream
func (sto SnapStore) GetSnapStream(snapURL string) (*http.Response, error) {
	return submitGETRequest(snapURL, sto.headers)
}

// Assertion retrieves an assertion from the store
func (sto SnapStore) Assertion(assertType, key string) (asserts.Assertion, error) {
	log.Printf("Download %s assertion\n", assertType)
	headers := map[string]string{"Accept": "application/x.ubuntu.assertion"}
	for k, v := range sto.headers {
		if k != "Accept" {
			headers[k] = v
		}
	}

	u := fmt.Sprintf("%s/assertions/%s/%s", apiBaseURL, assertType, key)
	resp, err := submitGETRequest(u, headers)
	if err != nil {
		log.Printf("Error fetching snap info: %v", err)
		return nil, err
	}

	defer resp.Body.Close()
	dec := asserts.NewDecoder(resp.Body)
	return dec.Decode()
}
