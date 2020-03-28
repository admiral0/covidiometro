package datastore

import (
	"covidiometro/util"
	"encoding/json"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"time"
)

func ScaricaDati() *Dati{
	tree, err := Clone()
	util.ErrFatal(err)

	dati := NuoviDati()
	massaggiaDati(tree, dati)

	return dati
}

func massaggiaDati(tree *object.Tree, d *Dati){
	//FIXME: Refactor di questa merda
	italiaFile, err := tree.File(DatiItalia)
	util.ErrFatal(err)
	italia, err := italiaFile.Contents()
	util.ErrFatal(err)

	regioniFile, err := tree.File(DatiRegioni)
	util.ErrFatal(err)
	regioni, err := regioniFile.Contents()
	util.ErrFatal(err)

	provinceFile, err := tree.File(DatiProvince)
	util.ErrFatal(err)
	province, err := provinceFile.Contents()
	util.ErrFatal(err)

	util.ErrFatal(json.Unmarshal([]byte(italia), &d.Italia))
	util.ErrFatal(json.Unmarshal([]byte(regioni), &d.Regioni))
	util.ErrFatal(json.Unmarshal([]byte(province), &d.Province))

	d.ItaliaRaw = italia //DEPRECATED

	d.AndroidBundle = BuildBundleV1(d)

	now := time.Now()

	d.lastUpdate = now
	d.ttl = now.Add(update)
}