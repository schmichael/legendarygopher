package lg

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
)

type Decoder interface {
	Decode(v interface{}) error
}

func New(d Decoder) (*World, error) {
	w := &World{}
	if err := d.Decode(w); err != nil {
		return nil, err
	}
	w.init()
	return w, nil
}

type World struct {
	XMLName            xml.Name             `xml:"df_world" json:"-"`
	Regions            []*Region            `xml:"regions>region" json:"regions"`
	UndergroundRegions []*UndergroundRegion `xml:"underground_regions>underground_region" json:"underground_regions"`
	Sites              []*Site              `xml:"sites>site" json:"sites"`
	siteidx            map[int]*Site
	Artifacts          []*Artifact `xml:"artifacts>artifact" json:"artifacts"`
	Figures            []*Figure   `xml:"historical_figures>historical_figure" json:"historical_figures"`
	figidx             map[int]*Figure

	// useless?
	EntityPopulations []*EntityPopulation `xml:"entity_populations>entity_population" json:"-"`

	Entities []*Entity `xml:"entities>entity" json:"entities"`
	Events   []*Event  `xml:"historical_events>historical_event" json:"historical_events"`
}

func (w *World) init() {
	w.siteidx = make(map[int]*Site, len(w.Sites))
	for _, s := range w.Sites {
		w.siteidx[s.ID] = s
	}

	w.figidx = make(map[int]*Figure, len(w.Figures))
	for _, f := range w.Figures {
		w.figidx[f.ID] = f
	}
}

func (w *World) Figure(id int) *Figure {
	return w.figidx[id]
}

func (w *World) Site(id int) *Site {
	return w.siteidx[id]
}

func (w *World) FigureEvents(id int) <-chan *Event {
	out := make(chan *Event, 100)
	go func() {
		defer close(out)
		for _, e := range w.Events {
			if e.FigureID == id || e.SlayerFigureID == id {
				out <- e
			}
		}
	}()
	return out
}

func (w *World) RenderEvent(e *Event) string {
	switch e.Type {
	case "destroyed site":
		return fmt.Sprintf("Site %d for Civ %d destroyed by Civ %d", e.SiteID, e.DefenderCivID, e.AttackerCivID)
	case "change hf state":
		if e.SiteID != -1 {
			return fmt.Sprintf("%s %s %s", w.Figure(e.FigureID), e.State, w.Site(e.SiteID))
		}
		return fmt.Sprintf("%s %s", w.Figure(e.FigureID), e.State)
	case "hf died":
		if slayer := w.Figure(e.SlayerFigureID); slayer != nil {
			return fmt.Sprintf("%s slayed by %s", w.Figure(e.FigureID), slayer)
		}
		return fmt.Sprintf("%s died", w.Figure(e.FigureID))
	default:
		return fmt.Sprint("Event %d in %d (unknown type %q)", e.ID, e.Year, e.Type)
	}
}

func (w *World) String() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("Regions\n")
	for _, r := range w.Regions {
		fmt.Fprintf(buf, "%-3d %-30s %s\n", r.ID, r.Name, r.Type)
	}
	for _, r := range w.UndergroundRegions {
		fmt.Fprintf(buf, "%-3d %-10s %d\n", r.ID, r.Type, r.Depth)
	}
	for _, s := range w.Sites {
		fmt.Fprintf(buf, "%-5d %-14s %-40s %-7s\n", s.ID, s.Type, s.Name, s.Coords)
	}
	for _, a := range w.Artifacts {
		fmt.Fprintf(buf, "%-5d %-40s %-30s\n", a.ID, a.Name, a.Item)
	}
	for _, f := range w.Figures {
		fmt.Fprintf(buf, "%-5d %-40s entities:%d sites:%d spheres:%s\n", f.ID, f.Name, len(f.Entities), len(f.Sites), strings.Join(f.Spheres, ","))
	}
	fmt.Fprintf(buf, "Events: %d\n", len(w.Events))
	return buf.String()
}

type Region struct {
	ID   int    `xml:"id" json:"id"`
	Name string `xml:"name" json:"name"`
	Type string `xml:"type" json:"type"`
}

type UndergroundRegion struct {
	ID    int    `xml:"id" json:"id"`
	Type  string `xml:"type" json:"type"`
	Depth int    `xml:"depth" json:"depth"`
}

type Site struct {
	ID     int    `xml:"id" json:"id"`
	Type   string `xml:"type" json:"type"`
	Name   string `xml:"name" json:"name"`
	Coords string `xml:"coords" json:"coords"`
	//TODO is this unused?!
	//Structures []*Structure `xml:"structures
}

type Artifact struct {
	ID   int    `xml:"id" json:"id"`
	Name string `xml:"name" json:"name"`
	Item string `xml:"item" json:"item"`
}

type Figure struct {
	ID         int           `xml:"id" json:"id"`
	Name       string        `xml:"name" json:"name"`
	Race       string        `xml:"race" json:"race"`
	Caste      string        `xml:"caste" json:"caste"`
	Appeared   int           `xml:"appeared" json:"appeared"`
	BirthYear  int           `xml:"birth_year" json:"birth_year"`
	DeathYear  int           `xml:"death_year" json:"death_year"`
	AssocTypes string        `xml:"associated_types" json:"associated_types"`
	Entities   []*EntityLink `xml:"entity_link" json:"entity_link"`
	Sites      []*SiteLink   `xml:"site_link" json:"site_link"`
	Spheres    []string      `xml:"sphere" json:"sphere"`
}

func (f *Figure) String() string { return f.Name }

type EntityLink struct {
	// Type may be "enemy" ...
	Type string `xml:"link_type" json:"link_type"`
	ID   int    `xml:"entity_id" json:"entity_id"`
}

type SiteLink struct {
	// Type may be "lair" ...
	Type string `xml:"link_type" json:"link_type"`
	ID   int    `xml:"site_id" json:"site_id"`
}

// EntityPopulation -- empty?!
type EntityPopulation struct {
	ID int `xml:"id"`
}

type Entity struct {
	ID   int    `xml:"id" json:"id"`
	Name string `xml:"name" json:"name"`
}

type Event struct {
	ID      int `xml:"id" json:"id"`
	Year    int `xml:"year" json:"year"`
	Seconds int `xml:"seconds72" json:"seconds,omitempty"`

	// Type values: "change hf state","hf died","destroyed_site"
	Type string `xml:"type" json:"type,omitempty"`

	// AttackerCivID is set when Type=destroyed_site
	AttackerCivID int `xml:"attacker_civ_id" json:"attacker_civ_id"`

	// DefenderCivID is set when Type=destroyed_site
	DefenderCivID int `xml:"defender_civ_id" json:"defender_civ_id"`

	FigureID       int `xml:"hfid" json:"hfid"`
	SlayerFigureID int `xml:"slayer_hfid" json:"slayer_hfid"`
	SlayerItemID   int `xml:"slayer_item_id", json:"slayer_item_id"`

	// State values: visiting,settled,wandering
	State string `xml:"state" json:"state,omitempty"`

	// SiteCivID is set when Type=destroyed_site
	//FIXME Doesn't match either Attacker or Defender Civ ID... not sure what's
	//      going on.
	SiteCivID int `xml:"site_civ_id" json:"site_civ_id"`

	SiteID         int    `xml:"site_id" json:"site_id"`
	SubregionID    int    `xml:"subregion_id" json:"subregion_id"`
	FeatureLayerID int    `xml:"feature_layer_id" json:"feature_layer_id"`
	Coords         string `xml:"coords" json:"coords,omitempty"`
	Cause          string `xml:"cause" json:"cause,omitempty"`
}
