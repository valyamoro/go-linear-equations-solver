package main

import "math"

func (eq Equation) Satisfies(x, y float64) bool {
	v := eq.A*x + eq.B*y
	const eps = 1e-9
	return math.Abs(v-eq.C) <= eps
}

func lineIntersection(eq1, eq2 Equation) *Point {
	denom := eq1.A*eq2.B - eq2.A*eq1.B
	if math.Abs(denom) < 1e-12 {
		return nil
	}
	x := (eq1.C*eq2.B - eq2.C*eq1.B) / denom
	y := (eq1.A*eq2.C - eq2.A*eq1.C) / denom
	return &Point{X: x, Y: y}
}

func solveSystem(eqs []Equation) ([]Point, bool) {
	var points []Point
	for i := 0; i < len(eqs); i++ {
		for j := i + 1; j < len(eqs); j++ {
			if p := lineIntersection(eqs[i], eqs[j]); p != nil {
				valid := true
				for _, eq := range eqs {
					if !eq.Satisfies(p.X, p.Y) {
						valid = false
						break
					}
				}
				if valid {
					points = append(points, *p)
				}
			}
		}
	}
	points = dedupPts(points, 1e-6)
	return points, len(points) > 0
}

func boundarySegment(eq Equation, xmin, xmax, ymin, ymax float64) (xs, ys []float64) {
	const eps = 1e-9
	clamp := func(v, lo, hi float64) float64 {
		if v < lo {
			return lo
		}
		if v > hi {
			return hi
		}
		return v
	}
	var pts []Point
	add := func(x, y float64) {
		if x >= xmin-eps && x <= xmax+eps && y >= ymin-eps && y <= ymax+eps {
			pts = append(pts, Point{clamp(x, xmin, xmax), clamp(y, ymin, ymax)})
		}
	}
	if math.Abs(eq.B) > eps {
		add(xmin, (eq.C-eq.A*xmin)/eq.B)
		add(xmax, (eq.C-eq.A*xmax)/eq.B)
	}
	if math.Abs(eq.A) > eps {
		add((eq.C-eq.B*ymin)/eq.A, ymin)
		add((eq.C-eq.B*ymax)/eq.A, ymax)
	}
	pts = dedupPts(pts, 1e-6)
	if len(pts) < 2 {
		return nil, nil
	}
	p0, p1 := pts[0], pts[len(pts)-1]
	return []float64{p0.X, p1.X}, []float64{p0.Y, p1.Y}
}

func dedupPts(pts []Point, eps float64) []Point {
	var out []Point
	for _, p := range pts {
		dup := false
		for _, q := range out {
			if math.Hypot(p.X-q.X, p.Y-q.Y) < eps {
				dup = true
				break
			}
		}
		if !dup {
			out = append(out, p)
		}
	}
	return out
}

func r3(v float64) float64 { return math.Round(v*1000) / 1000 }