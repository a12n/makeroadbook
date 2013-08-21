package main

import (
	"errors"
	"math"
	"sort"
	"time"
)

type Waypoint struct {
	Dist float64
	DistPrev float64
	Desc string
}

type Checkpoint struct {
	Waypoint
	Name string
	OpensAfter time.Duration `json:"-"`
	OpensAfterStr string `json:"OpensAfter"`
	ClosesAfter time.Duration `json:"-"`
	ClosesAfterStr string `json:"ClosesAfter"`
}

type RoadBook struct {
	Name string
	Checkpoints []Checkpoint
	Waypoints []Waypoint
}

type CheckpointArray []Checkpoint
func (s CheckpointArray) Len() int { return len(s) }
func (s CheckpointArray) Less(i, j int) bool { return s[i].Dist < s[j].Dist }
func (s CheckpointArray) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type WaypointArray []Waypoint
func (s WaypointArray) Len() int { return len(s) }
func (s WaypointArray) Less(i, j int) bool { return s[i].Dist < s[j].Dist }
func (s WaypointArray) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func SpeedRange(dist float64) (float64, float64) {
	tolerance := 2.5 / 100
	if dist < (200 + 200 * tolerance) {
		return 15, 34
	} else if dist < (400 + 400 * tolerance) {
		return 15, 32
	} else if dist < (600 + 600 * tolerance) {
		return 15, 30
	} else if dist < (1000 + 1000 * tolerance) {
		return 11.428, 28
	}
	return 13.333, 26
}

func CheckpointTimes(dist float64) (time.Duration, time.Duration) {
	if dist < 0.100 {
		return 0, 1 * time.Hour
	}
	min, max := SpeedRange(dist)
	start := time.Duration(dist / max * float64(time.Hour))
	end := time.Duration(dist / min * float64(time.Hour))
	start = (start / time.Second) * time.Second
	end = (end / time.Second) * time.Second
	return start, end
}

func NewRoadBook(gpx *Gpx) (*RoadBook, error) {
	if len(gpx.Trk) != 1 || len(gpx.Trk[0].TrkSeg) != 1 {
		return nil, errors.New("There must be exactly one track in the GPX file")
	}

	rbk := new(RoadBook)

	if gpx.Metadata != nil {
		if gpx.Metadata.Name != nil {
			rbk.Name = *gpx.Metadata.Name
		}
	}

	trk := gpx.Trk[0]
	trkseg := trk.TrkSeg[0]
	dist := trkseg.DistAlong()

	for _, gpxWpt := range gpx.Wpt {
		jMin := -1
		xMin, aMin := 0.0, 0.0
		for j := 0; j < len(trkseg.TrkPt) - 1; j++ {
			x, a := gpxWpt.TrackDist(*trkseg.TrkPt[j], *trkseg.TrkPt[j + 1])
			x = math.Abs(x)
			if jMin < 0 || x < xMin {
				jMin = j
				xMin = x
				aMin = a
			}
		}

		wpt := Waypoint{}
		wpt.Dist = (dist[jMin] + aMin) / 1000
		if gpxWpt.Desc != nil {
			wpt.Desc = *gpxWpt.Desc
		} else if gpxWpt.Name != nil {
			wpt.Desc = *gpxWpt.Name
		}
		rbk.Waypoints = append(rbk.Waypoints, wpt)

		if (gpxWpt.Type != nil && (*gpxWpt.Type == "Viewpoint" || *gpxWpt.Type == "Sightseeing")) {
			cpt := Checkpoint{}
			cpt.Waypoint = wpt
			if gpxWpt.Name != nil {
				cpt.Name = *gpxWpt.Name
			}
			rbk.Checkpoints = append(rbk.Checkpoints, cpt)
		}
	}

	sort.Sort(CheckpointArray(rbk.Checkpoints))
	sort.Sort(WaypointArray(rbk.Waypoints))

	for i := 0; i < len(rbk.Checkpoints); i++ {
		if i > 0 {
			rbk.Checkpoints[i].DistPrev =
				rbk.Checkpoints[i].Dist - rbk.Checkpoints[i - 1].Dist
		}
		rbk.Checkpoints[i].OpensAfter, rbk.Checkpoints[i].ClosesAfter =
			CheckpointTimes(rbk.Checkpoints[i].Dist)
		rbk.Checkpoints[i].OpensAfterStr = rbk.Checkpoints[i].OpensAfter.String()
		rbk.Checkpoints[i].ClosesAfterStr = rbk.Checkpoints[i].ClosesAfter.String()
	}

	for i := 0; i < len(rbk.Waypoints); i++ {
		if i > 0 {
			rbk.Waypoints[i].DistPrev =
				rbk.Waypoints[i].Dist - rbk.Waypoints[i - 1].Dist
		}
	}

	return rbk, nil
}
