package main

import (
	"encoding/xml"
	"log"
	"os"
)

type Wpt struct {
	XMLName xml.Name `xml:"wpt"`
	Lat float64 `xml:"lat,attr"`
	Lon float64 `xml:"lon,attr"`
	Name string `xml:"name"`
	Desc string `xml:"desc"`
	Sym string `xml:"sym"`
	Type string `xml:"type"`
}

type TrkPt struct {
	XMLName xml.Name `xml:"trkpt"`
	Lat float64 `xml:"lat,attr"`
	Lon float64 `xml:"lon,attr"`
	Ele float64 `xml:"ele"`
}

type TrkSeg struct {
	XMLName xml.Name `xml:"trkseg"`
	TrkPt []TrkPt `xml:"trkpt"`
}

type Trk struct {
	XMLName xml.Name `xml:"trk"`
	TrkSeg []TrkSeg `xml:"trkseg"`
}

type Gpx struct {
	XMLName xml.Name `xml:"gpx"`
	Wpt []Wpt `xml:"wpt"`
	Trk Trk `xml:"trk"`
}

func main() {
	gpx := new(Gpx)
	if err := xml.NewDecoder(os.Stdin).Decode(gpx); err != nil {
		log.Fatal(err)
	}
	if len(gpx.Trk.TrkSeg) != 1 {
		log.Fatal("There must be exactly one 'trkseg' element in the GPX file")
	}
	// TODO
	for _, w := range gpx.Wpt {
		log.Printf("Lat %f, Lon %f, Name '%s', Sym '%s', Type '%s'\n",
			w.Lat, w.Lon, w.Name, w.Sym, w.Type)
	}
	for _, p := range gpx.Trk.TrkSeg[0].TrkPt {
		log.Printf("Lat %f, Lon %f, Ele %f\n",
			p.Lat, p.Lon, p.Ele)
	}
}
