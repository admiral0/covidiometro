package android

import (
	"covidiometro/datastore"
	"log"
	"net/http"

	"goji.io"
	"goji.io/pat"
)

func RegisterHandlers(mux *goji.Mux) {
	mux.HandleFunc(pat.Get("/api/android/v1/"), BundleV1)
}

func BundleV1(w http.ResponseWriter, r *http.Request) {
	d := datastore.Holder.Get()
	_, err := w.Write(d.AndroidBundle)
	if err != nil {
		log.Println(err)
	}
}
