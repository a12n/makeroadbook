package main

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"os"
	"text/template"
)

func main() {
	var buf []byte
	var err error
	var gpx *Gpx = new(Gpx)
	var rbk *RoadBook
	var tpl *template.Template
	if err = xml.NewDecoder(os.Stdin).Decode(gpx); err != nil {
		log.Fatal(err)
	}
	if rbk, err = NewRoadBook(gpx); err != nil {
		log.Fatal(err)
	}
	if tpl, err = template.ParseFiles("template.txt"); err != nil {
		if buf, err = json.MarshalIndent(rbk, "", "  "); err != nil {
			log.Fatal(err)
		} else {
			os.Stdout.Write(buf)
		}
	} else {
		tpl.Execute(os.Stdout, rbk)
	}
}
