package main

type Point struct {
	x float64
	y float64
}

type Points []Point

func myFunc(a, b int) int {
	return a + b
}

func (points Points) Len() int {
	return cap(points)
}

func (points Points) Less(i, j int) bool {
	a := points[i]
	b := points[j]
	return a.x < b.x || (a.x == b.x && a.y < b.y)
}

func (points Points) Swap(i, j int) {
	points[i], points[j] = points[j], points[i]
}
