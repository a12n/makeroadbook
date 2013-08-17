package main

import (
	"encoding/xml"
	"io"
	"log"
	"math"
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

func rad(degree float64) float64 {
	return degree * math.Pi / 180
}

func (p1 Wpt) distTo(p2 Wpt) float64 {
	// TODO: Use Vincenty formula?
	lat1, lon1 := rad(p1.Lat), rad(p1.Lon)
	lat2, lon2 := rad(p2.Lat), rad(p2.Lon)
	diffLat, diffLon := lat2 - lat1, lon2 - lon1
	sinLat, sinLon := math.Sin(diffLat / 2), math.Sin(diffLon / 2)
	a := sinLat * sinLat + sinLon * sinLon * math.Cos(lat1) * math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1 - a))
	return 6371E3 * c
}

func readGpx(reader io.Reader) (*Gpx, error) {
	ans := new(Gpx)
	err := xml.NewDecoder(reader).Decode(ans)
	if err != nil {
		return nil, err
	}	
	return ans, nil
}

func main() {
	gpx, err := readGpx(os.Stdin)
	if err != nil {
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
