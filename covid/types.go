package covid

import (
	"github.com/Jeffail/gabs/v2"
	"time"
)

type Snapshot struct {
	Italy *gabs.Container
	Regions *gabs.Container
	Provinces *gabs.Container

	Ref RefInfo
}

type RefInfo struct {
	Updated time.Time
	Hash string
	Permalink string
}
