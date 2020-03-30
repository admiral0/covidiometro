package util

import (
	"encoding/json"
)

func Js(v interface{}) []byte {
	r, _ := json.Marshal(v)
	return r
}
