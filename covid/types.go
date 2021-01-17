package covid

import (
	"github.com/Jeffail/gabs/v2"
	"github.com/admiral0/covidiometro"
)

type Snapshot struct {
	Italy *gabs.Container
	Regions *gabs.Container
	Provinces *gabs.Container

	Ref covidiometro.RefInfo
}

