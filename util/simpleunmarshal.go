package util

import (
	"encoding/json"
)

func Js(v interface{}) []byte {
	r, _ := json.Marshal(v)
	return r
}

func JsRegioni(v *[]map[string]interface{}, regione int) []byte {
	nuovaMappa := make([]map[string]interface{}, 0)

	for _, mappa := range *v {
		codice := int(mappa["codice_regione"].(float64))
		if regione == codice {
			nuovaMappa = append(nuovaMappa, mappa)
		}
	}
	r, _ := json.Marshal(nuovaMappa)
	return r
}