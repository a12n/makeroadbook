package main

import (
	"encoding/xml"
	"log"
	"os"
)

type Wpt struct {
	Lat float64 `xml:"lat,attr"`
	Lon float64 `xml:"lon,attr"`
	Ele *float64 `xml:"ele"`
	Name *string `xml:"name"`
	Desc *string `xml:"desc"`
	Sym *string `xml:"sym"`
	Type *string `xml:"type"`
}

type TrkSeg struct {
	TrkPt []Wpt `xml:"trkpt"`
}

type Trk struct {
	TrkSeg []TrkSeg `xml:"trkseg"`
}

type Gpx struct {
	Wpt []Wpt `xml:"wpt"`
	Trk []Trk `xml:"trk"`
}

func main() {
	gpx := new(Gpx)
	if err := xml.NewDecoder(os.Stdin).Decode(gpx); err != nil {
		log.Fatal(err)
	}
	if len(gpx.Trk) != 1 || len(gpx.Trk[0].TrkSeg) != 1 {
		log.Fatal("There must be exactly one 'trkseg' element in the GPX file")
	}
	// TODO
	for _, p := range gpx.Wpt {
		log.Printf("%+v\n", p)
	}
	for _, p := range gpx.Trk[0].TrkSeg[0].TrkPt {
		log.Printf("%+v\n", p)
	}
}
