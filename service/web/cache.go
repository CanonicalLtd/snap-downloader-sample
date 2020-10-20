package web

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type snapAddRequest struct {
	Name string `json:"name"`
	Arch string `json:"arch"`
}

// CacheSnapList lists snaps in the watch list
func (srv Web) CacheSnapList(w http.ResponseWriter, r *http.Request) {
	snaps, err := srv.Cache.SnapList()
	if err != nil {
		formatStandardResponse("cache-add", err.Error(), w)
		return
	}

	formatRecordsResponse(snaps, w)
}

// CacheSnapAdd adds a snap to the watch list
func (srv Web) CacheSnapAdd(w http.ResponseWriter, r *http.Request) {
	req := srv.decodeSnapAdd(w, r)
	if req == nil {
		return
	}

	if err := snapValidate(req); err != nil {
		formatStandardResponse("cache-add", err.Error(), w)
		return
	}

	id, err := srv.Cache.SnapAdd(req.Name, req.Arch)
	if err != nil {
		formatStandardResponse("cache-add", err.Error(), w)
		return
	}

	formatRecordResponse(id, w)
}

// CacheSnapDelete removes a snap from the watch list
func (srv Web) CacheSnapDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if err := srv.Cache.SnapDelete(vars["id"]); err != nil {
		formatStandardResponse("cache-delete", err.Error(), w)
		return
	}

	formatStandardResponse("", "", w)
}

func (srv Web) decodeSnapAdd(w http.ResponseWriter, r *http.Request) *snapAddRequest {
	// Decode the JSON body
	req := snapAddRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	switch {
	// Check we have some data
	case err == io.EOF:
		formatStandardResponse("data", "No request data supplied.", w)
		return nil
		// Check for parsing errors
	case err != nil:
		formatStandardResponse("decode-json", err.Error(), w)
		return nil
	}
	return &req
}

func snapValidate(req *snapAddRequest) error {
	if req.Name == "" || req.Arch == "" {
		return fmt.Errorf("name and architecture must be entered")
	}
	return nil
}
