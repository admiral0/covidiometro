package main

import (
	"covidiometro/android"
	"covidiometro/datastore"
	"covidiometro/template"
	"covidiometro/util"
	"log"
	"net/http"
	"strconv"

	"github.com/markbates/pkger"
	"goji.io"
	"goji.io/pat"
)

func index(w http.ResponseWriter, r *http.Request) {
	dati := datastore.Holder.Get()
	template.Index(dati, w)
}

func about(w http.ResponseWriter, r *http.Request) {
	template.About(w)
}

func regione(w http.ResponseWriter, r *http.Request) {
	regione, err := strconv.Atoi(pat.Param(r, "regione"))
	dati := datastore.Holder.Get()
	if err == nil {
		template.Regione(dati, regione, w)
	}else {
		w.WriteHeader(404)
	}

}

func main() {
	util.SetupLogging()

	log.Println("Initializing in-memory repository")
	datastore.Inizializzazione()
	datastore.ScaricaDati()
	go datastore.Updater()

	mux := goji.NewMux()
	// Handlers
	mux.HandleFunc(pat.Get("/"), index)
	mux.HandleFunc(pat.Get("/regione/:regione"), regione)
	mux.HandleFunc(pat.Get("/about"), about)

	android.RegisterHandlers(mux)
	// Static files handling
	dir := http.FileServer(pkger.Dir("/assets"))
	mux.Handle(pat.Get("/static/*"), dir)

	util.ErrFatal(http.ListenAndServe(":8000", mux))
}
