package lg

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
)

func New(r io.Reader) (*DFWorld, error) {
	w := &DFWorld{}
	if err := xml.NewDecoder(r).Decode(w); err != nil {
		return nil, err
	}
	return w, nil
}

type DFWorld struct {
	XMLName            xml.Name             `xml:"df_world"`
	Regions            []*Region            `xml:"regions>region"`
	UndergroundRegions []*UndergroundRegion `xml:"underground_regions>underground_region"`
	Sites              []*Site              `xml:"sites>site"`
	Artifacts          []*Artifact          `xml:"artifacts>artifact"`
	Figures            []*Figure            `xml:"historical_figures>historical_figure"`
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
		fmt.Fprintf(buf, "%-5d %-40s %-5d\n", f.ID, f.Name, len(f.Entities))
	}
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
	Appeared   int           `xml:"appeared"`
	BirthYear  int           `xml:"birth_year"`
	DeathYear  int           `xml:"death_year"`
	AssocTypes string        `xml:"associated_types"`
	Entities   []*EntityLink `xml:"entity_link"`
}

type EntityLink struct {
	Type string `xml:"link_type"`
	ID   int    `xml:"entity_id"`
}
