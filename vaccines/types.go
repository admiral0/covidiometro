package vaccines

import (
	"github.com/Jeffail/gabs/v2"
	"github.com/admiral0/covidiometro"
)

type Snapshot struct {
	Regions *gabs.Container

	Ref covidiometro.RefInfo
}