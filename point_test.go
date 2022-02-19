package main

import (
	"fmt"
	"sort"
	"testing"
)

func TestPoints(t *testing.T) {
	var points Points
	var (
		a = Point{x: 0.001, y: 1.102}
		b = Point{x: 0, y: 0}
		c = Point{x: 10, y: -5.1}
		d = Point{x: -100, y: 1000}
		e = Point{x: -100, y: 0}
	)
	points = Points{a, b, c, d, e}
	if length := points.Len(); length != 5 {
		t.Fatalf("Invalid length %d, expected %d", length, 3)
	}
	prevFirst, prevSecond := points[0], points[1]
	points.Swap(0, 1)
	if newFirst, newSecond := points[0], points[1]; newFirst == prevFirst || newSecond == prevSecond {
		t.Fatalf("Didn't swap?")
	}
	fmt.Println("About to sort array", points)
	sort.Sort(points)
	fmt.Println("Got sorted array", points)
	if points[0].x != -100 || points[0].y != 0 {
		t.Fatalf("Didn't sort?")
	}
}
