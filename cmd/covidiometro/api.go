package main

import (
	"github.com/go-chi/chi"
	"net/http"
)

func covidiometroApi(ds *DataSources) http.Handler {
	r := chi.NewRouter()
	return r
}