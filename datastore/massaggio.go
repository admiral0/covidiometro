package datastore

import (
	"covidiometro/util"
	"encoding/json"
	"time"
)

func ScaricaDati() *Dati{
	tree, err := Clone()
	util.ErrFatal(err)
	italiaFile, err := tree.File(DatiItalia)
	util.ErrFatal(err)
	italia, err := italiaFile.Contents()
	util.ErrFatal(err)

	dati := NuoviDati()
	util.ErrFatal(json.Unmarshal([]byte(italia), &dati.Italia))
	dati.ItaliaRaw = italia
	now := time.Now()
	dati.lastUpdate = now
	dati.ttl = now.Add(update)

	return dati
}