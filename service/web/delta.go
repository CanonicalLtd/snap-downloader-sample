package web

import (
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
)

// DeltaGet fetches a snap delta from the cache
func (srv Web) DeltaGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// validate the parameters
	fromRevision, err := strconv.Atoi(vars["fromRevision"])
	if err != nil {
		formatStandardResponse("delta", "the `from` revision is not a valid integer", w)
		return
	}
	toRevision, err := strconv.Atoi(vars["toRevision"])
	if err != nil {
		formatStandardResponse("delta", "the `to` revision is not a valid integer", w)
		return
	}

	// fetch or generate the delta
	delta, err := srv.Cache.Delta(vars["name"], vars["arch"], fromRevision, toRevision)
	if err != nil {
		formatStandardResponse("delta", err.Error(), w)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+path.Base(delta))
	w.Header().Set("Content-Type", StreamHeader)

	// return the delta file stream
	download, err := os.Open(delta)
	if err != nil {
		formatStandardResponse("download", err.Error(), w)
		return
	}
	defer download.Close()

	io.Copy(w, download)
}
