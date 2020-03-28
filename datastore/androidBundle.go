package datastore

import (
	"covidiometro/util"
	"encoding/json"
)

func BuildBundleV1(dati *Dati) []byte {
	bundle := make(map[string]interface{})

	bundle["italia"] = dati.Italia
	bundle["regioni"] = dati.Regioni
	bundle["province"] = dati.Province


	bundleString, err := json.Marshal(bundle)
	util.ErrFatal(err)

	return bundleString
}