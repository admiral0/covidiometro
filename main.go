package main

import (
	"covidiometro/datastore"
	"covidiometro/util"
	"github.com/markbates/pkger"
	"goji.io"
	"goji.io/pat"
	"io"
	"log"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request){
	f, _ := pkger.Open("/assets/index.html")
	_, err := io.Copy(w, f)
	if err != nil {
		log.Panicln(err)
	}
}

func about(w http.ResponseWriter, r *http.Request){
	f, _ := pkger.Open("/assets/about.html")
	_, err := io.Copy(w, f)
	if err != nil {
		log.Panicln(err)
	}
}

func main() {
	util.SetupLogging()
	log.Println("Initializing in-memory repository")
	consumer := make(chan datastore.Dati, 1)
	go datastore.Updater(consumer)

	mux := goji.NewMux()

	// Handlers
	mux.HandleFunc(pat.Get("/"), index)
	mux.HandleFunc(pat.Get("/about"), about)
	mux.HandleFunc(pat.Get("/raw/italia.json"), func(writer http.ResponseWriter, request *http.Request) {
		dati := <- consumer
		raw := []byte(dati.ItaliaRaw)
		_, err := writer.Write(raw)
		util.ErrFatal(err)
	})

	// Static files handling
	dir := http.FileServer(pkger.Dir("/assets"))
	mux.Handle(pat.Get("/static/*"), dir)

	util.ErrFatal(http.ListenAndServe(":8000", mux))
}
