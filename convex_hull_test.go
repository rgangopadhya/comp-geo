package main

import (
	"fmt"
	"testing"
)

var Square = Points{makePoint(0, 0), makePoint(1, 0), makePoint(1, 1), makePoint(0, 1)}

var Triangle = Points{makePoint(0, 0), makePoint(1, 0), makePoint(1, 1)}

var Pentagon = Points{
	makePoint(-0.62, 1.56),
	makePoint(0.34, 1.54),
	makePoint(0.62, 1),
	makePoint(0.36, 0.32),
	makePoint(-1, 1),
}

var Hexagon = Points{
	// these points form the boundary
	makePoint(-0.67, 2.01),
	makePoint(1.39, 1.55),
	makePoint(0.97, 0.63),
	makePoint(0.41, -0.37),
	makePoint(-1.45, -0.53),
	makePoint(-1.9, 1),
}

var InteriorToHexagonPoints = Points{
	// these points are within the interior
	makePoint(-0.97, 0.57),
	makePoint(-0.63, 1.43),
	makePoint(0.57, 1.29),
	makePoint(0.23, 0.29),
	makePoint(-0.57, -0.25),
}

func pointsEqual(pointsA, pointsB Points) bool {
	if len(pointsA) != len(pointsB) {
		return false
	}
	makeSet := func(points Points) map[Point]bool {
		result := make(map[Point]bool, len(points))
		for _, point := range points {
			result[point] = true
		}
		return result
	}
	setA := makeSet(pointsA)
	for _, point := range pointsB {
		if _, found := setA[point]; !found {
			return false
		}
	}
	return true
}

func TestConvexHull(t *testing.T) {
	polygons := []Points{Square, Triangle, Pentagon, Hexagon}
	for i, polygon := range polygons {
		fmt.Println("About to test", i, polygon)
		hull := convexHull(polygon)
		if !pointsEqual(hull, polygon) {
			t.Fatal("Expected polygon to be hull, but wasn't", i, hull, polygon)
		}
	}

	var hexagonCopy = make(Points, len(Hexagon))
	copy(hexagonCopy, Hexagon)
	var hexagonWithContainedPoints = make(Points, len(Hexagon)+len(InteriorToHexagonPoints))
	hexagonWithContainedPoints = append(hexagonCopy, InteriorToHexagonPoints...)
	hull := convexHull(hexagonWithContainedPoints)
	if !pointsEqual(hull, Hexagon) {
		t.Fatal("Expected hexagonWithContainedPoints hull to be hexagon")
	}
}
