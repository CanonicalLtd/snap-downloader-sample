package web

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/slimjim777/snap-downloader/service/cache"
	"github.com/slimjim777/snap-downloader/service/store"
	"net/http"
)

const (
	defaultPort    = "8888"
	defaultDocRoot = "./static"
)

// Web implements the web service
type Web struct {
	Store store.Service
	Cache cache.Service
}

// NewWebService starts a new web service
func NewWebService(snapStore store.Service, cacheSrv cache.Service) *Web {
	return &Web{
		Store: snapStore,
		Cache: cacheSrv,
	}
}

// Start the web service
func (srv Web) Start() error {
	listenOn := fmt.Sprintf(":%s", defaultPort)
	fmt.Printf("Starting service on port %s\n", listenOn)
	return http.ListenAndServe(listenOn, srv.Router())
}

// Router returns the application router
func (srv Web) Router() *mux.Router {
	// Start the web service router
	router := mux.NewRouter()

	// unauthenticated routes
	router.Handle("/v1/login", Middleware(http.HandlerFunc(srv.LoginUser))).Methods("POST")
	router.Handle("/v1/auth", Middleware(http.HandlerFunc(srv.Macaroon))).Methods("GET")

	// authenticated routes
	router.Handle("/v1/snaps", srv.MiddlewareWithAuth(http.HandlerFunc(srv.CacheSnapList))).Methods("GET")
	router.Handle("/v1/snaps", srv.MiddlewareWithAuth(http.HandlerFunc(srv.CacheSnapAdd))).Methods("POST")
	router.Handle("/v1/snaps/{id}", srv.MiddlewareWithAuth(http.HandlerFunc(srv.CacheSnapDelete))).Methods("DELETE")
	router.Handle("/v1/downloads", srv.MiddlewareWithAuth(http.HandlerFunc(srv.CacheDownloadList))).Methods("GET")
	router.Handle("/v1/downloads/{name}/{filename}", srv.MiddlewareWithAuth(http.HandlerFunc(srv.CacheDownloadFile))).Methods("GET")

	// serve the static path
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir(defaultDocRoot)))
	router.PathPrefix("/static/").Handler(fs)

	router.Handle("/", srv.MiddlewareWithAuth(http.HandlerFunc(srv.Index))).Methods("GET")
	router.Handle("/login", Middleware(http.HandlerFunc(srv.Index))).Methods("GET")

	return router
}
