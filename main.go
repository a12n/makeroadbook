package main

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"os"
)

func main() {
	var buf []byte
	var err error
	var gpx *Gpx = new(Gpx)
	var rbk *RoadBook
	if err = xml.NewDecoder(os.Stdin).Decode(gpx); err != nil {
		log.Fatal(err)
	}
	if rbk, err = NewRoadBook(gpx); err != nil {
		log.Fatal(err)
	}
	if buf, err = json.MarshalIndent(rbk, "", "  "); err != nil {
		log.Fatal(err)
	}
	os.Stdout.Write(buf)
}
