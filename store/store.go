package store

import (
	"encoding/json"
	"fmt"
	"log"
)

// SnapStore interacts with a brand store
type SnapStore struct {
	headers map[string]string
}

// NewStore creates a new store client
func NewStore() *SnapStore {
	// check if we have cached headers
	headers, err := readHeaders()
	if err != nil {
		return &SnapStore{}
	}

	// use the cached headers, so no login needed
	return &SnapStore{
		headers: headers,
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

	return cacheHeaders(sto.headers)
}

// SnapInfo lists the snaps in a brand store
func (sto SnapStore) SnapInfo(name string) (*ResponseSnapInfo, error) {
	u := fmt.Sprintf("%sv2/snaps/info/%s", apiBaseURL, name) //info
	resp, err := submitGETRequest(u, sto.headers)
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
