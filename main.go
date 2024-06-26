package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/kd"
)

func main() {

	// Whatever our soltuion, it is loaded in new line by new line
	// so can either output as a single print or multiple printlns

	// In this problem, it does not matter if a truck goes to the farthest point first, or the closest, as long as it get there
	// Want to create a timeout from the get go, so we at least always return one solution we get

	// Looking at KD tree, but we want to store end nodes to start nodes.
	// actually, no, we just need to store start nodes, we use the end node to find the next closest
	// so k-d tree of start nodes, Search(endNode)
	//https://en.wikipedia.org/wiki/K-d_tree
	res, err := BasicRun(os.Args[1], implementation1)
	if err != nil {
		panic(err)
	}

	for _, r := range res {
		str := ""
		for _, v := range r {
			str = fmt.Sprintf("%s,%d", str, v)
		}
		str = fmt.Sprintf("[%s]", strings.TrimPrefix(str, ","))
		fmt.Println(str)
	}
}

type Results [][]int
type Implementation func(tree *kd.KD[*BasicLoad], loads []*BasicLoad) (res Result)

// Have this return the results
func BasicRun(filepath string, implementation Implementation) (results Results, err error) {
	loads, err := loadLoadFile(filepath)
	if err != nil {
		return
	}
	// minNumberOfDrivers := calculateMinimumDriverGuess(loads)
	// _ = minNumberOfDrivers // Not sure how this number will be used yet
	// tree := buildKDTree(loads)
	// res := implementation(tree, loads)
	results = implementation1LoadSwapPerformance(loads)
	return
}

func buildKDTree(loads []*BasicLoad) (tree *kd.KD[*BasicLoad]) {
	tree = kd.New[*BasicLoad](kd.O[*BasicLoad]{
		Data: loads,
		K:    2,
		N:    8, // Played around this number, at this scale seems to have little effect
	})

	// Check if we need to do this
	tree.Balance()
	return
}

func calculateMinimumDriverGuess(loads []*BasicLoad) (min int) {
	totalDistance := float64(0)
	for _, l := range loads {
		totalDistance += l.Distance
	}

	min = int(math.Ceil(totalDistance / float64(MAX_DISTANCE)))
	return
}

func ConstructLoadMap(loads []*BasicLoad) (loadMap map[int]*BasicLoad) {
	loadMap = map[int]*BasicLoad{}
	for _, l := range loads {
		loadMap[l.LoadNumber] = l
	}
	return
}
func calculateMinimumCost(results [][][]int, loadMap map[int]*BasicLoad) (solution [][]int, minCost float64) {
	minCost = math.MaxFloat64
	for _, res := range results {
		c := calculateCost(res, loadMap)
		if c < minCost {
			minCost = c
			solution = res
		}
	}
	return
}

func calculateCost(result [][]int, loadMap map[int]*BasicLoad) (cost float64) {
	cost += 500 * float64(len(result))

	for _, truckPath := range result {
		vectors := []vector.V{}
		vectors = append(vectors, *HOME_VECTOR)
		for _, stop := range truckPath {
			load := loadMap[stop]
			vectors = append(vectors, *load.Pickup, *load.DropOff)
		}
		vectors = append(vectors, *HOME_VECTOR)
		for x := 0; x < len(vectors) - 1; x ++ {
			cost += DistanceBetweenTwoPoints(&vectors[x], &vectors[x+1])
		}
	}
	return cost
}
