package store

import (
	"github.com/snapcore/snapd/jsonutil/safejson"
	"github.com/snapcore/snapd/snap"
	"time"
)

// ResponseSnapInfo is the response from a snap info call
type ResponseSnapInfo struct {
	ChannelMap []*storeInfoChannelSnap `json:"channel-map"`
	Snap       storeSnap               `json:"snap"`
	Name       string                  `json:"name"`
	SnapID     string                  `json:"snap-id"`
}

// storeInfoChannelSnap is the snap-in-a-channel of which the channel map is made
type storeInfoChannelSnap struct {
	storeSnap
	Channel storeInfoChannel `json:"channel"`
}

// storeSnap holds the information sent as JSON by the store for a snap.
type storeSnap struct {
	Architectures []string           `json:"architectures"`
	Base          string             `json:"base"`
	Confinement   string             `json:"confinement"`
	Contact       string             `json:"contact"`
	CreatedAt     string             `json:"created-at"` // revision timestamp
	Description   safejson.Paragraph `json:"description"`
	Download      storeSnapDownload  `json:"download"`
	Epoch         snap.Epoch         `json:"epoch"`
	License       string             `json:"license"`
	Name          string             `json:"name"`
	Prices        map[string]string  `json:"prices"` // currency->price,  free: {"USD": "0"}
	Private       bool               `json:"private"`
	Publisher     snap.StoreAccount  `json:"publisher"`
	Revision      int                `json:"revision"` // store revisions are ints starting at 1
	SnapID        string             `json:"snap-id"`
	SnapYAML      string             `json:"snap-yaml"` // optional
	Summary       safejson.String    `json:"summary"`
	Title         safejson.String    `json:"title"`
	Type          snap.Type          `json:"type"`
	Version       string             `json:"version"`
	Website       string             `json:"website"`
	StoreURL      string             `json:"store-url"`

	// TODO: not yet defined: channel map

	// media
	Media []storeSnapMedia `json:"media"`

	CommonIDs []string `json:"common-ids"`
}

// storeInfoChannel is the channel description included in info results
type storeInfoChannel struct {
	Architecture string    `json:"architecture"`
	Name         string    `json:"name"`
	Risk         string    `json:"risk"`
	Track        string    `json:"track"`
	ReleasedAt   time.Time `json:"released-at"`
}

type storeSnapDownload struct {
	Sha3_384 string           `json:"sha3-384"`
	Size     int64            `json:"size"`
	URL      string           `json:"url"`
	Deltas   []storeSnapDelta `json:"deltas"`
}

type storeSnapDelta struct {
	Format   string `json:"format"`
	Sha3_384 string `json:"sha3-384"`
	Size     int64  `json:"size"`
	Source   int    `json:"source"`
	Target   int    `json:"target"`
	URL      string `json:"url"`
}

type storeSnapMedia struct {
	Type   string `json:"type"` // icon/screenshot
	URL    string `json:"url"`
	Width  int64  `json:"width"`
	Height int64  `json:"height"`
}
