package datastore

import (
	"covidiometro/util"
	"encoding/json"
	"time"

	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func ScaricaDati() *Dati {
	tree, err := Clone()
	util.ErrFatal(err)

	dati := NuoviDati()
	massaggiaDati(tree, dati)

	return dati
}

func massaggiaDati(tree *object.Tree, d *Dati) {

	carica := func(nome string, dest *[]map[string]interface{}) {
		file, err := tree.File(nome)
		util.ErrFatal(err)
		contenuti, err := file.Contents()
		util.ErrFatal(err)
		util.ErrFatal(json.Unmarshal([]byte(contenuti), &dest))
	}

	carica(DatiItalia, &d.Italia)
	carica(DatiRegioni, &d.Regioni)
	carica(DatiProvince, &d.Province)

	d.AndroidBundle = BuildBundleV1(d)

	stuzzica(d)

	popolaRegioni(d)

	now := time.Now()

	d.lastUpdate = now
	d.ttl = now.Add(update)
}

func stuzzica(dati *Dati) {
	campo_dt(dati.Italia, "tamponi", "tamponi_dt")
	campo_dt(dati.Italia, "deceduti", "deceduti_dt")
	campo_dt(dati.Italia, "dimessi_guariti", "dimessi_guariti_dt")
	campo_dt(dati.Italia, "terapia_intensiva", "terapia_intensiva_dt")
	campo_dt(dati.Italia, "ricoverati_con_sintomi", "ricoverati_con_sintomi_dt")

	campo_unfold_dt(dati.Regioni, "codice_regione","tamponi", "tamponi_dt")
	campo_unfold_dt(dati.Regioni, "codice_regione", "deceduti", "deceduti_dt")
	campo_unfold_dt(dati.Regioni, "codice_regione", "dimessi_guariti", "dimessi_guariti_dt")
	campo_unfold_dt(dati.Regioni, "codice_regione", "terapia_intensiva", "terapia_intensiva_dt")
	campo_unfold_dt(dati.Regioni, "codice_regione", "ricoverati_con_sintomi", "ricoverati_con_sintomi_dt")
}

func campo_dt(lista []map[string]interface{}, originale string, nuovo string){
	var precedente float64 = 0
	for _, mappa := range lista {
		mappa[nuovo] = mappa[originale].(float64) - precedente
		precedente = mappa[originale].(float64)
	}
}

func campo_unfold_dt(lista []map[string]interface{}, disambiguatore string, originale string, nuovo string){
	precedente := make(map[float64]float64)
	for _, mappa := range lista {
		codice := mappa[disambiguatore].(float64)
		p, ok := precedente[codice]
		if ! ok {
			p = 0
		}
		mappa[nuovo] = mappa[originale].(float64) - p
		precedente[codice] = mappa[originale].(float64)
	}
}

func popolaRegioni(dati *Dati) {
	for _, mappa := range dati.Regioni {
		regione := mappa["denominazione_regione"].(string)
		codice := int(mappa["codice_regione"].(float64))
		_, ok := dati.MappaRegioni[codice]
		if ! ok {
			dati.MappaRegioni[codice] = regione
		}
	}
}