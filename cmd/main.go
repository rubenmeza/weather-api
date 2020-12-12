package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if ok == false {
		log.Fatal(errors.New("enviroment variable PORT must be defined"))
		return
	}
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Mount("/weather", weatherResource{}.Routes())

	http.ListenAndServe(":"+port, r)
}
