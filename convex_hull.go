package main

import (
	"sort"
)

func remove(points Points, index int) Points {
	// The worst case of remove is O(n) here
	// TODO: https://github.com/ross-oreto/go-tree
	// Implement this with an AVL BST (balanced binary search tree)
	// which supports efficient insertion, lookup, and deletion.
	// Building the BST would be O(nlogn) for the sort, but we've already
	// sorted the points. From there you just continually insert.
	return append(points[:index], points[index+1:]...)
}

func convexHull(points Points) Points {
	/*
		convexHull algorithm
		points: a set of points in the Plane
		output: a set of points which form the convexHull (the smallest perimeter convex polygon containing `points`)
			in clockwise order
	*/
	sort.Sort(points)

	isHullValid := func(hull Points) bool {
		if len(hull) <= 2 {
			// From what we know, this is good!
			return true
		}
		lastThree := hull[len(hull)-3:]
		a, b, c := lastThree[0], lastThree[1], lastThree[2]
		return makesRightTurn(a, b, c)
	}

	incrementHull := func(hull Points, point Point) Points {
		hull = append(hull, point)

		for !isHullValid(hull) {
			// delete the middle of the last three until we make a right turn
			// note that the total number of executions of this inner loop across
			// the entire set of points is no more than the set size (n).
			hull = remove(hull, len(hull)-2)
		}
		return hull
	}

	upperHull := Points{points[0], points[1]}
	for i := 2; i < len(points); i++ {
		upperHull = incrementHull(upperHull, points[i])
	}

	lowerHull := Points{points[len(points)-1], points[len(points)-2]}
	for i := len(points) - 3; i >= 0; i-- {
		lowerHull = incrementHull(lowerHull, points[i])
	}

	// Remove the first and last point of lowerHull to dedupe
	lowerHull = append(lowerHull[1 : len(lowerHull)-1])
	return append(upperHull, lowerHull...)
}
