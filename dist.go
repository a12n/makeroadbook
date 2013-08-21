package main

import m "math"

const (
	ArcDegree float64 = m.Pi / 180
	EarthRadius = 6371E3
	Radian = 180 / m.Pi
)

// Aviation Formulary V1.46, by Ed Williams:
// http://s.best.vwh.net/avform.htm

func (p1 Wpt) DistTo(p2 Wpt) float64 {
	// TODO: Use Vincenty formula?
	lat1, lon1 := ArcDegree * p1.Lat, ArcDegree * p1.Lon
	lat2, lon2 := ArcDegree * p2.Lat, ArcDegree * p2.Lon
	a := m.Sin((lat1 - lat2) / 2)
	b := m.Sin((lon1 - lon2) / 2)
	d := 2 * m.Asin(m.Sqrt(a * a + m.Cos(lat1) * m.Cos(lat2) * b * b))
	return EarthRadius * d
}

func (p1 Wpt) CourseTo(p2 Wpt) float64 {
	lat1, lon1 := ArcDegree * p1.Lat, ArcDegree * p1.Lon
	lat2, lon2 := ArcDegree * p2.Lat, ArcDegree * p2.Lon
	if m.Cos(lat1) < 1E-9 {
		if lat1 > 0 {
			// Starting from N pole
			return 180
		} else {
			// Starting from S pole
			return 0
		}
	}
	y := m.Sin(lon1 - lon2) * m.Cos(lat2)
	x := m.Cos(lat1) * m.Sin(lat2) - m.Sin(lat1) * m.Cos(lat2) * m.Cos(lon1 - lon2)
	return Radian * m.Mod(m.Atan2(y, x), 2 * m.Pi)
}

// Positive XTD means right of course, negative means left.
func (p Wpt) TrackDist(pA, pB Wpt) (float64, float64) {
	dP := pA.DistTo(p) / EarthRadius
	cP := pA.CourseTo(p) * ArcDegree
	cB := pA.CourseTo(pB) * ArcDegree

	xt := m.Asin(m.Sin(dP) * m.Sin(cP - cB))

	sd := m.Sin(dP)
	sx := m.Sin(xt)
	at := m.Asin(m.Sqrt(sd * sd - sx * sx) / m.Cos(xt))

	return EarthRadius * xt, EarthRadius * at
}

func (s *TrkSeg) DistAlong() []float64 {
	ans := make([]float64, len(s.TrkPt))
	for i := 1; i < len(s.TrkPt); i++ {
		ans[i] = ans[i - 1] + s.TrkPt[i - 1].DistTo(*s.TrkPt[i])
	}
	return ans
}
