package main

import (
	"sort"
)

func convexHull(points Points) Points {
	sort.Sort(points)
	upperHull := points[:3]
	for i := 3; i < len(points); i++ {
		lastThree := upperHull[len(upperHull)-3:]
		a, b, c := lastThree[0], lastThree[1], lastThree[2]
		for len(upperHull) >= 2 && !makesRightTurn(a, b, c) {

		}
	}
	return upperHull
}
