package main

import (
	"github.com/downflux/go-kd/kd"
)

// Using the minimum number of needed drivers, begin with each one at a node
// visually each one would start at a far cluster and do clustering
// thats a good question, code for finding clusters
// really need to figure out an evaluation function
func implementation4(loads []*BasicLoad)(results [][]int){
	tree := buildKDTree(loads)
	loadMap := ConstructLoadMap(loads)
	driverCount := calculateMinimumDriverGuess(loads)
	truckDriversPaths := make([][]int, driverCount) // []truckPath
	truckDriversMilage := make([]float64, driverCount)
	truckDriversDone := make([]bool, driverCount) // Keep track of which trucks have decided they are full 
	orderFunc := LoadsOrderFuncs[4] // Longest routes at end of list, then we can pop off the back and improve array performance
	orderFunc(loads)
	pathAdded := true
	// assign a single location to each vehicle in a loop
	for len(loads) > 0{
		// This is poor performance adding one driver at a time once our starter drivers are 
		if !pathAdded {
			// If we failed to add any paths last time, we add a new truck to our fleet
			truckDriversPaths = append(truckDriversPaths, truckPath{})
			truckDriversMilage = append(truckDriversMilage, 0)
			truckDriversDone = append(truckDriversDone, false)
		}
		pathAdded = false

		for driverIndex, driverPath := range truckDriversPaths {
			if truckDriversDone[driverIndex] {
				continue
			}
			pathAdded = true
			// Add a first place for us to go
			if len(driverPath) == 0 {
				targetLoad := loads[len(loads)-1]
				loads = loads[:len(loads)-1]
				tree.Remove(*targetLoad.Pickup, targetLoad.Equal)

				truckDriversPaths[driverIndex] = []int{targetLoad.LoadNumber}
				truckDriversMilage[driverIndex] += StartingDistance(targetLoad)
				continue
			}

			// Grab where the truck is now
			thisTrucksCurrentSpot := driverPath[len(driverPath)-1]
			startingLocation := loadMap[thisTrucksCurrentSpot]

			// Each of these searches can probably be done in it's own thread, but they are completed so fast I don't think there would be a speed improvement
			possibleLoads := kd.KNN(tree, *startingLocation.DropOff, 1, func(bl *BasicLoad) bool {
				// If we have visited this node, we skip over it
				// can instead remove from tree, therefor improving performance as nodes re removed
				// So here where we decide multiple costs, like it would be better if we did some advanced smoozy and explored each different option
				possibleCost := truckDriversMilage[driverIndex]
				possibleCost += DistanceBetweenTwoPoints(startingLocation.DropOff, bl.Pickup) + bl.Distance + bl.DistanceToHome
				if possibleCost > MAX_DISTANCE {
					return false
				}
				return true
			})

			if len(possibleLoads) == 0 {
				truckDriversDone[driverIndex] = true
				continue
			}

			nextLoad := possibleLoads[0]
			tree.Remove(*nextLoad.Pickup, nextLoad.Equal)
			loads = RemoveFromArray(loads, nextLoad.LoadNumber)
			truckDriversPaths[driverIndex] = append(truckDriversPaths[driverIndex], nextLoad.LoadNumber)
			truckDriversMilage[driverIndex] = truckDriversMilage[driverIndex] + DistanceBetweenTwoPoints(startingLocation.DropOff, nextLoad.Pickup) + nextLoad.Distance
		}
	}

	return truckDriversPaths
}