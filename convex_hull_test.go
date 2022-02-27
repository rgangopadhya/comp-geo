package main

import (
	"testing"
)

var Square = Points{makePoint(0, 0), makePoint(1, 0), makePoint(1, 1), makePoint(0, 1)}

var Triangle = Points{makePoint(0, 0), makePoint(1, 0), makePoint(1, 1)}

var Hexagon = Points{makePoint(0, 0), makePoint(0.5, 0), makePoint(0.75, 0.25), makePoint(1, 0.5), makePoint(1, 1)}

var Pentagon = Points{
	// these points form the boundary
	makePoint(-0.67, 2.01),
	makePoint(1.39, 1.55),
	makePoint(0.97, 0.63),
	makePoint(0.41, -0.37),
	makePoint(-1.45, -0.53),
	makePoint(-1.9, 1),
}

var PentagonWithContainedPoints = append(Pentagon, Points{
	// these points are within the interior
	makePoint(-0.97, 0.57),
	makePoint(-0.63, 1.43),
	makePoint(0.57, 1.29),
	makePoint(0.23, 0.29),
	makePoint(-0.57, -0.25),
}...)

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
		hull := convexHull(polygon)
		if !pointsEqual(hull, polygon) {
			t.Fatal("Expected polygon to be hull, but wasn't", i, hull, polygon)
		}
	}
}
