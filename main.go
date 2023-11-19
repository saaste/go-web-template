package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", handleIndex)
	fileServer(r, "/static", http.Dir("ui/static"))

	err := http.ListenAndServe(fmt.Sprintf(":%s", "8000"), r)
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server closed")
	} else if err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {

	type TemplateData struct {
		Title string
	}

	data := TemplateData{
		Title: "Hello World!",
	}

	t, err := template.ParseFiles("ui/templates/index.html")
	if err != nil {
		log.Printf("Failed to parse template: %v\n", err)
		http.Error(w, "Internal Server Error", 500)
	}
	err = t.ExecuteTemplate(w, "index", data)
	if err != nil {
		log.Printf("Faild to execute template: %v", err)
		http.Error(w, "Internal Server Error", 500)
	}
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
