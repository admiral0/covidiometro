package vaccines

import (
	"github.com/Jeffail/gabs/v2"
	"time"
)

type Snapshot struct {
	Regions *gabs.Container

	Ref RefInfo
}

type RefInfo struct {
	Updated time.Time
	Hash string
	Permalink string
}