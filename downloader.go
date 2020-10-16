package main

import (
	"flag"
	"fmt"
	"github.com/slimjim777/snap-downloader/service/datastore/sqlite"
	"github.com/slimjim777/snap-downloader/service/store"
	"github.com/slimjim777/snap-downloader/service/web"
	"log"
	"os"
)

var storeClient *store.SnapStore

func main() {
	// the arguments are used to authenticate with the store and store the macaroon
	parseArgs()

	// set up the dependency chain
	db, _ := sqlite.NewDatabase()
	storeClient = store.NewStore(db)

	// Start the web service
	srv := web.NewWebService(storeClient)
	log.Fatal(srv.Start())
}

func parseArgs() {
	var (
		configureOnly bool
		email         string
		password      string
		otp           string
		storeID       string
		series        string
	)

	flag.BoolVar(&configureOnly, "configure", false, "Configure the application and exit")
	flag.StringVar(&email, "email", "", "Email to authenticate with the store")
	flag.StringVar(&password, "password", "", "Password to authenticate with the store")
	flag.StringVar(&otp, "otp", "", "One-time PIN to authenticate with the store")
	flag.StringVar(&storeID, "store", "", "Store ID of the brand store")
	flag.StringVar(&series, "series", "16", "The model assertion series (default 16")
	flag.Parse()

	if !configureOnly {
		// No changes if we're not configuring the app
		return
	}

	// login to the store and cache the macaroon
	if err := storeClient.LoginUser(email, password, otp, storeID, series); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
