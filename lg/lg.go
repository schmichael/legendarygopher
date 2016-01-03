package lg

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

func New(r io.Reader) (*DFWorld, error) {
	w := &DFWorld{}
	if err := xml.NewDecoder(r).Decode(w); err != nil {
		return nil, err
	}
	w.init()
	return w, nil
}

type DFWorld struct {
	XMLName            xml.Name             `xml:"df_world"`
	Regions            []*Region            `xml:"regions>region"`
	UndergroundRegions []*UndergroundRegion `xml:"underground_regions>underground_region"`
	Sites              []*Site              `xml:"sites>site"`
	siteidx            map[int]*Site
	Artifacts          []*Artifact `xml:"artifacts>artifact"`
	Figures            []*Figure   `xml:"historical_figures>historical_figure"`
	figidx             map[int]*Figure

	// useless?
	EntityPopulations []*EntityPopulation `xml:"entity_populations>entity_population"`

	Entities []*Entity `xml:"entities>entity"`
	Events   []*Event  `xml:"historical_events>historical_event"`
}

func (w *DFWorld) init() {
	w.siteidx = make(map[int]*Site, len(w.Sites))
	for _, s := range w.Sites {
		w.siteidx[s.ID] = s
	}

	w.figidx = make(map[int]*Figure, len(w.Figures))
	for _, f := range w.Figures {
		w.figidx[f.ID] = f
	}
}

func (w *DFWorld) Figure(id int) *Figure {
	return w.figidx[id]
}

func (w *DFWorld) Site(id int) *Site {
	return w.siteidx[id]
}

func (w *DFWorld) FigureEvents(id int) <-chan *Event {
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

func (w *DFWorld) RenderEvent(e *Event) string {
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

func (w *DFWorld) String() string {
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
	ID   int    `xml:"id"`
	Name string `xml:"name"`
	Type string `xml:"type"`
}

type UndergroundRegion struct {
	ID    int    `xml:"id"`
	Type  string `xml:"type"`
	Depth int    `xml:"depth"`
}

type Site struct {
	ID     int    `xml:"id"`
	Type   string `xml:"type"`
	Name   string `xml:"name"`
	Coords string `xml:"coords"`
	//Structures []*Structure `xml:"structures // unused?!
}

type Artifact struct {
	ID   int    `xml:"id"`
	Name string `xml:"name"`
	Item string `xml:"item"`
}

type Figure struct {
	ID         int           `xml:"id"`
	Name       string        `xml:"name"`
	Race       string        `xml:"race"`
	Caste      string        `xml:"caste"`
	Appeared   int           `xml:"appeared"`
	BirthYear  int           `xml:"birth_year"`
	DeathYear  int           `xml:"death_year"`
	AssocTypes string        `xml:"associated_types"`
	Entities   []*EntityLink `xml:"entity_link"`
	Sites      []*SiteLink   `xml:"site_link"`
	Spheres    []string      `xml:"sphere"`
}

func (f *Figure) String() string { return f.Name }

type EntityLink struct {
	// Type may be "enemy" ...
	Type string `xml:"link_type"`
	ID   int    `xml:"entity_id"`
}

type SiteLink struct {
	// Type may be "lair" ...
	Type string `xml:"link_type"`
	ID   int    `xml:"site_id"`
}

// EntityPopulation -- empty?!
type EntityPopulation struct {
	ID int `xml:"id"`
}

type Entity struct {
	ID   int    `xml:"id"`
	Name string `xml:"name"`
}

type Event struct {
	ID             int    `xml:"id"`
	Year           int    `xml:"year"`
	Seconds        int    `xml:"seconds72"`
	Type           string `xml:"type"`            // "change hf state","hf died","destroyed_site"
	AttackerCivID  int    `xml:"attacker_civ_id"` // used for type=destroyed_site
	DefenderCivID  int    `xml:"defender_civ_id"` // used for type=destroyed_site
	FigureID       int    `xml:"hfid"`
	SlayerFigureID int    `xml:"slayer_hfid"`
	SlayerItemID   int    `xml:"slayer_item_id"`
	State          string `xml:"state"`       // visiting,settled,wandering
	SiteCivID      int    `xml:"site_civ_id"` // used for type=destroyed_site
	SiteID         int    `xml:"site_id"`
	SubregionID    int    `xml:"subregion_id"`
	FeatureLayerID int    `xml:"feater_layer_id"`
	Coords         string `xml:"coords"`
	Cause          string `xml:"cause"`
}
