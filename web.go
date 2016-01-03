package main

import (
	"encoding/json"
	"fmt"
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
	figurest   = template.Must(template.New("figures").Parse(string(MustAsset("assets/templates/figures.html"))))
	figuret    = template.Must(template.New("figure").Parse(string(MustAsset("assets/templates/figure.html"))))
)

type server struct {
	World *lg.World
}

//go:generate go-bindata assets/...
func runserver(bind string, w *lg.World) {
	s := &server{World: w}

	// Serverside rendered html
	http.HandleFunc("/", wrap(s.listHandler(indext)))
	http.HandleFunc("/artifacts", wrap(s.listHandler(artifactst)))
	http.HandleFunc("/entities", wrap(s.listHandler(entitiest)))
	http.HandleFunc("/events", wrap(s.listHandler(eventst)))
	http.HandleFunc("/figures", wrap(s.listHandler(figurest)))
	http.HandleFunc("/figures/", wrap(s.figureHandler))
	http.HandleFunc("/assets/", wrap(s.assetHandler))

	// API
	http.HandleFunc("/api/world", wrap(s.jsonify(w)))
	http.HandleFunc("/api/artifacts", wrap(s.jsonify(w.Artifacts)))
	http.HandleFunc("/api/entities", wrap(s.jsonify(w.Entities)))
	http.HandleFunc("/api/events", wrap(s.jsonify(w.Events)))
	http.HandleFunc("/api/figures", wrap(s.jsonify(w.Figures)))
	http.HandleFunc("/api/sites", wrap(s.jsonify(w.Sites)))
	http.HandleFunc("/api/regions", wrap(s.jsonify(w.Regions)))
	http.HandleFunc("/api/undergroundregions", wrap(s.jsonify(w.UndergroundRegions)))

	if err := http.ListenAndServe(bind, nil); err != nil {
		log.Fatal(err)
	}
}

func wrap(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer log.Printf(r.URL.Path)
		f(w, r)
	}
}

func (s *server) listHandler(t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := t.Execute(w, s); err != nil {
			log.Printf("error executing template %s: %v", t.Name(), err)
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
	}
}

func (s *server) figureHandler(w http.ResponseWriter, r *http.Request) {
	id := 0
	if _, err := fmt.Sscanf(r.URL.Path, "/figures/%d", &id); err != nil {
		log.Printf("error getting figure id from %q: %v", r.URL.Path, err)
		w.WriteHeader(404)
		w.Write([]byte(err.Error()))
		return
	}
	fig := s.World.Figure(id)
	if fig == nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "not found: figure %d", id)
		return
	}
	context := struct {
		Figure *lg.Figure
		World  *lg.World
	}{fig, s.World}
	if err := figuret.Execute(w, context); err != nil {
		log.Printf("error executing template %s: %v", figuret.Name(), err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
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

func (s *server) jsonify(collection interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if err := json.NewEncoder(w).Encode(collection); err != nil {
			w.WriteHeader(500)
			log.Printf("error encoding json for %q %T: %v", r.URL.Path, collection, err)
			return
		}
	}
}
