package main

import (
	"maps"

	"github.com/downflux/go-geometry/nd/vector"
	"github.com/downflux/go-kd/kd"
)

const IMP2_TREE_EXP_SIZE = 3

type truckPath []int            // this is a path of truck a
type possibleTruckPaths [][]int // This is the paths truck a could take
type PathsOfSomeTrucks [][]int  // This is the paths of truck a, b, and c

// Prunning from a tree gives a 2 - 3 ms speedup, don't think the memory cost
// of making new tree clones will be worth it but idk?

// Implementation 1 will grab the furthest points from home to start from
// Features this is missing is actually iterating through and finding the best solution

// This no work at all, it just runs forever (i.e. > 10 minutes)
func Implementation2RecursionStartingPoint2(tree *kd.KD[*BasicLoad], loads []*BasicLoad, alreadyResults PathsOfSomeTrucks) (finalResults [][][]int) {
	results := make([][][]int, 0)

	clonedLoads := make([]*BasicLoad, len(loads))
	copy(clonedLoads, loads)
	// We get all the next possible paths we can take
	possibleTruckRoutes := Implementation2PossibleNewTruckRoutes(tree, clonedLoads, alreadyResults)
	// No more truck routes to be made
	if len(possibleTruckRoutes) == 0 {
		// fmt.Println("no routes found")
		return
	}

	if len(alreadyResults) == 0 {
		results = append(results, possibleTruckRoutes)
	} else {
		// The set of routes for truck c, want to add to the alreadyResults which is (a1), (b4)
		for _, atruckRoute := range possibleTruckRoutes {
			trucksRoute := make(PathsOfSomeTrucks, len(alreadyResults))
			// we have a, b
			// The already existing paths is not getting copied into the routing
			for aTrucksPathIndex := range alreadyResults {
				trucksRoute[aTrucksPathIndex] = make([]int, len(alreadyResults[aTrucksPathIndex]))
				copy(trucksRoute[aTrucksPathIndex], alreadyResults[aTrucksPathIndex])
			}
			// We now have one permutaion of trucks a, b, c's route
			trucksRoute = append(trucksRoute, atruckRoute)
			results = append(results, trucksRoute)
		}
	}

	for r := range results {
		someResults := Implementation2RecursionStartingPoint2(tree, clonedLoads, results[r])
		if len(someResults) > 0 {
			finalResults = append(finalResults, someResults...)
		}else{
			finalResults = append(finalResults, results[r])
		}
	}
	return finalResults
}

// Given routes from the set of trucks so far, generate a set of possible truck routes
func Implementation2PossibleNewTruckRoutes(tree *kd.KD[*BasicLoad], loads []*BasicLoad, alreadyMadeTruckRoutes [][]int) (results [][]int) {
	// Our starting point will be passed to implementation2TruckRecursion, we will then take each of those results
	// and change our remaining loads. Probably want to filter somehow on some of these, otherwise we will grow incredibly fast
	visitedLoads := make(map[int]struct{})
	for _, truckRoutes := range alreadyMadeTruckRoutes {
		for _, visitedNode := range truckRoutes {
			loads = RemoveFromArray(loads, visitedNode)
			visitedLoads[visitedNode] = struct{}{}
		}

	}
	if len(loads) == 0 {
		// fmt.Println("visitedloads 1: ", len(visitedLoads))
		// We exit
		return
	}
	if len(visitedLoads) > 200 {
		panic("visited loads is toooo long")
	}
	startingNode := loads[0]
	// fmt.Printf("starting node: %+v\n", startingNode)
	distanceSoFar := StartingDistance(startingNode)

	loads = RemoveFromArray(loads, startingNode.LoadNumber)
	visitedLoads[startingNode.LoadNumber] = struct{}{}
	if len(loads) == 0 {
		// fmt.Println("visitedloads 2: ", len(visitedLoads))
		// We exit
		return
	}

	results = implementation2TruckRecursion(tree, loads[0].DropOff, distanceSoFar, visitedLoads)
	if len(results) == 0 {
		results = [][]int{[]int{startingNode.LoadNumber}}
	} else {
		// We are a node on the trail
		// Adding which node we visited to the result
		thisNodeArray := []int{startingNode.LoadNumber}
		for x := range results {
			results[x] = append(thisNodeArray, results[x]...)
		}
	}
	return
}

// So we run this as if we are looking at the routes for a single truck
// At the end of the tree, where we have no more option nodes, we would return a single int, then that would be combined with the roots value
// currentMilage is how far we traveled to get to our current node
func implementation2TruckRecursion(tree *kd.KD[*BasicLoad], currentNode *vector.V, currentMilage float64, visitedLoads map[int]struct{}) (results possibleTruckPaths) {
	results = make([][]int, 0)
	possibleLoads := kd.KNN(tree, *currentNode, 2, func(bl *BasicLoad) bool {
		// If we have visited this node, we skip over it
		// can instead remove from tree, therefor improving performance as nodes re removed
		_, ok := visitedLoads[bl.LoadNumber]
		if ok {
			return false
		}
		// So here where we decide multiple costs, like it would be better if we did some advanced smoozy and explored each different option
		possibleDistance := currentMilage
		possibleDistance += DistanceBetweenTwoPoints(currentNode, bl.Pickup) + bl.Distance + bl.DistanceToHome
		if possibleDistance > MAX_DISTANCE {
			return false
		}

		return true
	})

	// Possible loads is where we can go to next
	for _, load := range possibleLoads {
		tempCurrentMillage := currentMilage + DistanceBetweenTwoPoints(currentNode, load.Pickup) + load.Distance
		visitedLoadsCopy := maps.Clone(visitedLoads)
		visitedLoadsCopy[load.LoadNumber] = struct{}{}
		// Add in the array of new results we just acquired
		tempResults := implementation2TruckRecursion(tree, load.DropOff, tempCurrentMillage, visitedLoadsCopy)
		// We are the leaf
		if len(tempResults) == 0 {
			results = possibleTruckPaths{truckPath{load.LoadNumber}}
		} else {
			// We are a node on the trail
			// Adding which node we visited to the result
			thisNodeArray := []int{load.LoadNumber}
			for x := range tempResults {
				tempResults[x] = append(thisNodeArray, tempResults[x]...)
			}
			results = append(results, tempResults...)
		}
	}

	return results
}

func StartingDistance(firstNode *BasicLoad) (distance float64) {
	distance = firstNode.Distance + DistanceBetweenTwoPoints(HOME_VECTOR, firstNode.Pickup)
	return distance
}

// Ways to go about this
// Where does branching occur
// Fleet Level: So each driver is maximized somehow against the fleet
// Car level: Use a measurement to determine if the car has maximized itself well enough

// My evaluation function is going to be percent used driving
// lets try without the queue evaluation. Maybe can still be fast enough?
// but still need to determine which routing was best overall

// For our queue valueation, not ure how you could really get it good
// when exploring for a single trucks route, I can't think of an evaluation
// that would account for total number of trucks needed. One option is
// Just maximizing the number of stops a single truck can do.
// not sure what kind of value that would get out
