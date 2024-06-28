package main

const IMPLEMENTATION_5_DEPTH = 4
const IMPLEMENTATION_5_WIDTH = 3

// I want to take the starting list idea from implementation 4 and set a max depth of the search. At that point, you return up your options and determine
// which search so far has the best performance. Then you start from that point and continuing on
func implementation5(loads []*BasicLoad) {

}

// Do we want to evaluate the 'score' at the node level, or deeper?
func implementation5Recursion(depth int) (routes [][][]int) {
	if depth > IMPLEMENTATION_5_DEPTH {
		return routes
	}
	return
}

// func evaluateBestOption

type implementation5Parameter struct {
	TruckDriversPaths  [][]int
	TruckDriversMilage []float64
	TruckDriversDone   []bool
	LoadMap            map[int]*BasicLoad // Might not need this
	Loads              []*BasicLoad
}
