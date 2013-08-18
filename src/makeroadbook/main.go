package main

import (
	"container/heap"
	"flag"
	"fmt"
	"gpx"
	"log"
	"math"
	"os"
	"strings"
)

// TODO: Download input from HTTP URI?
// TODO: HTML templates for output?
// TODO: Space-partitioning for TrkPt?

type outEntry struct {
	dist float64
	desc string
}

type outQueue []*outEntry

func (s outQueue) Len() int {
	return len(s)
}

func (s outQueue) Less(i, j int) bool {
	return s[i].dist < s[j].dist
}

func (s outQueue) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *outQueue) Pop() interface{} {
	n := len(*s)
	ans := (*s)[n - 1]
	*s = (*s)[0 : (n - 1)]
	return ans
}

func (s *outQueue) Push(x interface{}) {
	*s = append(*s, x.(*outEntry))
}

func distAlong(s *gpx.TrkSeg) []float64 {
	ans := make([]float64, len(s.TrkPt))
	for i := 1; i < len(s.TrkPt); i++ {
		ans[i] = ans[i - 1] + s.TrkPt[i - 1].DistTo(*s.TrkPt[i])
	}
	return ans
}

func main() {
	flag.Parse()
	var err error
	var input *gpx.Gpx
	if flag.NArg() > 0 {
		input, err = gpx.ReadFile(flag.Arg(0))
	} else {
		input, err = gpx.Read(os.Stdin)
	}
	if err != nil {
		log.Fatal(err)
	}
	if len(input.Trk) != 1 || len(input.Trk[0].TrkSeg) != 1 {
		log.Fatal("There must be exactly one track in the GPX file")
	}

	trk := input.Trk[0]
	trkseg := trk.TrkSeg[0]

	dist := distAlong(trkseg)

	output := &outQueue{}
	heap.Init(output)

	for _, wpt := range input.Wpt {
		jMin := -1
		xMin, aMin := 0.0, 0.0
		for j := 0; j < len(trkseg.TrkPt) - 1; j++ {
			x, a := wpt.TrackDist(*trkseg.TrkPt[j], *trkseg.TrkPt[j + 1])
			x = math.Abs(x)
			if jMin < 0 || x < xMin {
				jMin = j
				xMin = x
				aMin = a
			}
		}
		entry := &outEntry{ dist: dist[jMin] + aMin, desc: "" }
		if wpt.Desc != nil {
			entry.desc = *wpt.Desc
		} else if wpt.Name != nil {
			entry.desc = *wpt.Name
		}
		heap.Push(output, entry)
	}

	for output.Len() > 0 {
		entry := heap.Pop(output).(*outEntry)
		fmt.Printf("%s  %s\n",
			strings.Replace(fmt.Sprintf("%05.1f", entry.dist / 1000), ".", ",", 1),
			entry.desc)
	}
}
