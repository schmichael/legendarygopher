package main

import (
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/schmichael/legendarygopher/lg"
)

var (
	indext     = template.Must(template.New("index").Parse(string(MustAsset("assets/templates/index.html"))))
	artifactst = template.Must(template.New("artifacts").Parse(string(MustAsset("assets/templates/artifacts.html"))))
	entitiest  = template.Must(template.New("entities").Parse(string(MustAsset("assets/templates/entities.html"))))
	eventst    = template.Must(template.New("events").Parse(string(MustAsset("assets/templates/events.html"))))
)

type server struct {
	World *lg.DFWorld
}

//go:generate go-bindata assets/...
func runserver(bind string, w *lg.DFWorld) {
	s := &server{World: w}
	http.HandleFunc("/", s.templateHandler(indext))
	http.HandleFunc("/artifacts", s.templateHandler(artifactst))
	http.HandleFunc("/entities", s.templateHandler(entitiest))
	http.HandleFunc("/events", s.templateHandler(eventst))
	http.HandleFunc("/assets/", s.assetHandler)
	if err := http.ListenAndServe(bind, nil); err != nil {
		log.Fatal(err)
	}
}

func (s *server) templateHandler(t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := t.Execute(w, s); err != nil {
			log.Printf("error executing template %s: %v", t.Name(), err)
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
	}
}

func (s *server) assetHandler(w http.ResponseWriter, r *http.Request) {
	// drop leading "/"
	path := strings.TrimLeft(r.URL.Path, "/")

	a, err := Asset(path)
	if err != nil {
		log.Printf("error loading asset %q: %v", path, err)
		w.WriteHeader(404)
		return
	}
	switch {
	case strings.HasSuffix(path, ".css"):
		w.Header().Set("Content-Type", "text/css")
	default:
		// If we don't recognize the type, don't return it
		log.Printf("unrecognized file type %q", path)
		w.WriteHeader(404)
		return
	}
	w.Write(a)
}
