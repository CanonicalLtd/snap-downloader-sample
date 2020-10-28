package web

import (
	"encoding/json"
	"io"
	"net/http"
)

// WatchLastRun retrieves the last time the watch daemon ran
func (srv Web) WatchLastRun(w http.ResponseWriter, r *http.Request) {
	lastRun := srv.Watch.LastRun()

	formatRecordResponse(lastRun, w)
}

// WatchInterval fetches the watch daemon interval
func (srv Web) WatchInterval(w http.ResponseWriter, r *http.Request) {
	interval := srv.Watch.GetInterval()

	formatRecordResponse(interval, w)
}

// WatchSetInterval fetches the watch daemon interval
func (srv Web) WatchSetInterval(w http.ResponseWriter, r *http.Request) {
	interval := srv.decodeWatchInterval(w, r)
	if interval == 0 {
		return
	}

	if err := srv.Watch.SetInterval(interval); err != nil {
		formatStandardResponse("interval", err.Error(), w)
		return
	}

	formatRecordResponse(interval, w)
}

func (srv Web) decodeWatchInterval(w http.ResponseWriter, r *http.Request) int {
	// Decode the JSON body
	var settings map[string]int
	err := json.NewDecoder(r.Body).Decode(&settings)
	switch {
	// Check we have some data
	case err == io.EOF:
		formatStandardResponse("data", "No request data supplied.", w)
		return 0
		// Check for parsing errors
	case err != nil:
		formatStandardResponse("decode-json", err.Error(), w)
		return 0
	}
	if settings["value"] <= 0 {
		formatStandardResponse("invalid", "the watch interval must be at least 5 seconds", w)
		return 0
	}

	return settings["value"]
}
