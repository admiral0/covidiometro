package main

import (
	"github.com/markbates/pkger"
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	basedir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	log.Info().Str("basedir", basedir).Msg("loading in directory")

	ds, err := InitializeRepositories(basedir)
	if err != nil {
		log.Err(err).Msg("could not load repos")
		panic(err)
	}
	router := chi.NewRouter()
	router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		f, err := pkger.Open("/index.html")
		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte("could not open index.html"))
		}
		defer f.Close()
		_, err = io.Copy(writer, f)
		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte("could not read index.html"))
		}
	})
	router.Mount("/assets", http.FileServer(pkger.Dir("/assets")))
	router.Mount("/api", covidiometroApi(ds))

	http.ListenAndServe(":9000", router)
}
