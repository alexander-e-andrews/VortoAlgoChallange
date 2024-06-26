package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"

	"github.com/downflux/go-kd/kd"
)

var loadsOrderFunction = func(loads []*BasicLoad) {
	sort.Slice(loads, func(i, j int) bool { return loads[i].DistanceToHome > loads[j].DistanceToHome })
}

// Implementation 1 will grab the furthest points from home to start from
// Features this is missing is actually iterating through and finding the best solution
func implementation1(tree *kd.KD[*BasicLoad], loads []*BasicLoad) (res Result) {
	res.Results = make(Results, 0)
	// Organize start points from furthest to closest from Home
	// Since further to closest, less function is actually a more function

	loadsOrderFunction(loads)
	// visited loads is the list of all loads that have been visited by everyone
	// For loop is the loop for each driver
	for {
		myVisitedNodes := make([]int, 0)
		// NO more loads to process, no more drivers are needed
		if len(loads) == 0 {
			break
		}
		var currentCost float64
		startingLocation := HOME_VECTOR
		nextLoad := loads[0]
		// Eeewwwww made mistake here, accidently had nextLoad.DistanceToHome which is from end node
		currentCost += DistanceBetweenTwoPoints(startingLocation, nextLoad.Pickup) + nextLoad.Distance
		loads = RemoveFromArray(loads, nextLoad.LoadNumber)
		tree.Remove(*nextLoad.Pickup, nextLoad.Equal)
		myVisitedNodes = append(myVisitedNodes, nextLoad.LoadNumber)
		startingLocation = nextLoad.DropOff
		for {
			possibleLoads := kd.KNN(tree, *startingLocation, 1, func(bl *BasicLoad) bool {
				// If we have visited this node, we skip over it
				// can instead remove from tree, therefor improving performance as nodes re removed
				// So here where we decide multiple costs, like it would be better if we did some advanced smoozy and explored each different option
				possibleCost := currentCost
				possibleCost += DistanceBetweenTwoPoints(startingLocation, bl.Pickup) + bl.Distance + bl.DistanceToHome
				if possibleCost > MAX_DISTANCE {
					return false
				}

				return true
			})
			// We ran out of nodes, we need to say that this is done, and we need a new driver
			if len(possibleLoads) == 0 {
				res.Results = append(res.Results, myVisitedNodes)
				break
			}
			nextLoad = possibleLoads[0]
			tree.Remove(*nextLoad.Pickup, nextLoad.Equal)
			// This can be moved to after this for loop
			loads = RemoveFromArray(loads, nextLoad.LoadNumber)
			// We have a node to add to our system
			myVisitedNodes = append(myVisitedNodes, nextLoad.LoadNumber)
			currentCost += DistanceBetweenTwoPoints(startingLocation, nextLoad.Pickup) + nextLoad.Distance
			startingLocation = nextLoad.DropOff
		}
	}

	return
}

func IsInArray(arr []int, val int) (inside bool) {
	for _, v := range arr {
		if val == v {
			return true
		}
	}
	return false
}

// https://go.dev/wiki/SliceTricks#:~:text=a)%2Bn)%2C%20a...)%0A%7D-,Filter%20(in%20place),-n%20%3A%3D%200%0Afor
func RemoveArrayFromArray(arr []*BasicLoad, remove []int) []*BasicLoad {
	n := 0
	for _, x := range arr {
		// can remove from the remove list as well to possible improve speed
		if IsInArray(remove, x.LoadNumber) {
			arr[n] = x
			n++
		} else {
			panic("what is happening")
		}
	}
	arr = arr[:n]

	return arr
}

func RemoveFromArray(arr []*BasicLoad, remove int) []*BasicLoad {
	//fmt.Printf("Starting Array Length: %d  %+v\n", len(arr), remove)
	n := 0
	for _, x := range arr {
		// can remove from the remove list as well to possible improve speed
		if x.LoadNumber != remove {
			arr[n] = x
			n++
		}
	}
	arr = arr[:n]

	return arr
}

// Going to write as recursive for now, but thats bad and need to change

// Another interesting first sort is running clockwise from center
/*
func sortBoxesClockwise(boxes []MonitorInfoBlock) []MonitorInfoBlock {
	// Find the center of all boxes
	var totalX, totalY float64
	for _, box := range boxes {
		totalX += float64(box.Center.X)
		totalY += float64(box.Center.Y)
	}
	centerX := totalX / float64(len(boxes))
	centerY := totalY / float64(len(boxes))

	// Sort boxes based on angle from center point
	sort.Slice(boxes, func(i, j int) bool {
		p1 := boxes[i].Center
		p2 := boxes[j].Center
		angle1 := math.Atan2((float64(p1.Y) - centerY), (float64(p1.X) - centerX))
		angle2 := math.Atan2((float64(p2.Y) - centerY), (float64(p2.X) - centerX))

		return angle1 < angle2
	})

	return boxes
}
*/

func implementation1ClockWiseFromCenter(tree *kd.KD[*BasicLoad], loads []*BasicLoad) (res Result) {
	res.Results = make(Results, 0)
	// Organize start points from furthest to closest from Home
	// Since further to closest, less function is actually a more function

	sort.Slice(loads, func(i, j int) bool {
		p1 := loads[i].Pickup
		p2 := loads[j].Pickup
		angle1 := math.Atan2((float64(p1.X(1))), (float64(p1.X(0))))
		angle2 := math.Atan2((float64(p2.X(1))), (float64(p2.X(0))))

		return angle1 < angle2
	})

	// visited loads is the list of all loads that have been visited by everyone
	visitedLoads := make([]int, 0, len(loads)) // Wonder how using a map instead would change performance
	// For loop is the loop for each driver
	for {
		myVisitedNodes := make([]int, 0)
		// NO more loads to process, no more drivers are needed
		if len(loads) == 0 {
			break
		}
		var currentCost float64
		startingLocation := HOME_VECTOR
		nextLoad := loads[0]
		// Eeewwwww made mistake here, accidently had nextLoad.DistanceToHome which is from end node
		currentCost += DistanceBetweenTwoPoints(startingLocation, nextLoad.Pickup) + nextLoad.Distance
		loads = RemoveFromArray(loads, nextLoad.LoadNumber)
		visitedLoads = append(visitedLoads, nextLoad.LoadNumber)
		myVisitedNodes = append(myVisitedNodes, nextLoad.LoadNumber)
		startingLocation = nextLoad.DropOff
		for {
			possibleLoads := kd.KNN(tree, *startingLocation, 1, func(bl *BasicLoad) bool {
				// If we have visited this node, we skip over it
				// can instead remove from tree, therefor improving performance as nodes re removed
				if IsInArray(visitedLoads, bl.LoadNumber) {
					return false
				}
				// So here where we decide multiple costs, like it would be better if we did some advanced smoozy and explored each different option
				possibleCost := currentCost
				possibleCost += DistanceBetweenTwoPoints(startingLocation, bl.Pickup) + bl.Distance + bl.DistanceToHome
				if possibleCost > MAX_DISTANCE {
					return false
				}

				return true
			})
			// We ran out of nodes, we need to say that this is done, and we need a new driver
			if len(possibleLoads) == 0 {
				res.Results = append(res.Results, myVisitedNodes)
				break
			}
			nextLoad = possibleLoads[0]
			tree.Remove(*nextLoad.Pickup, nextLoad.Equal)
			loads = RemoveFromArray(loads, nextLoad.LoadNumber)
			// We have a node to add to our system
			visitedLoads = append(visitedLoads, nextLoad.LoadNumber)
			myVisitedNodes = append(myVisitedNodes, nextLoad.LoadNumber)
			currentCost += DistanceBetweenTwoPoints(startingLocation, nextLoad.Pickup) + nextLoad.Distance
			startingLocation = nextLoad.DropOff
		}
	}

	return
}

func implementation1MapHasVisitedNoTreeRemove(tree *kd.KD[*BasicLoad], loads []*BasicLoad) (res Result) {
	res.Results = make(Results, 0)
	// Organize start points from furthest to closest from Home
	// Since further to closest, less function is actually a more function

	sort.Slice(loads, func(i, j int) bool { return loads[i].DistanceToHome > loads[j].DistanceToHome })
	// visited loads is the list of all loads that have been visited by everyone
	visitedLoads := make(map[int]struct{}) // Wonder how using a map instead would change performance
	// For loop is the loop for each driver
	for {
		myVisitedNodes := make([]int, 0)
		// NO more loads to process, no more drivers are needed
		if len(loads) == 0 {
			break
		}
		var currentCost float64
		startingLocation := HOME_VECTOR
		nextLoad := loads[0]
		// Eeewwwww made mistake here, accidently had nextLoad.DistanceToHome which is from end node
		currentCost += DistanceBetweenTwoPoints(startingLocation, nextLoad.Pickup) + nextLoad.Distance
		loads = RemoveFromArray(loads, nextLoad.LoadNumber)
		visitedLoads[nextLoad.LoadNumber] = struct{}{}
		myVisitedNodes = append(myVisitedNodes, nextLoad.LoadNumber)
		startingLocation = nextLoad.DropOff
		for {
			possibleLoads := kd.KNN(tree, *startingLocation, 1, func(bl *BasicLoad) bool {
				// If we have visited this node, we skip over it
				// can instead remove from tree, therefor improving performance as nodes re removed
				_, ok := visitedLoads[bl.LoadNumber]
				if ok {
					return false
				}
				// So here where we decide multiple costs, like it would be better if we did some advanced smoozy and explored each different option
				possibleCost := currentCost
				possibleCost += DistanceBetweenTwoPoints(startingLocation, bl.Pickup) + bl.Distance + bl.DistanceToHome
				if possibleCost > MAX_DISTANCE {
					return false
				}

				return true
			})
			// We ran out of nodes, we need to say that this is done, and we need a new driver
			if len(possibleLoads) == 0 {
				res.Results = append(res.Results, myVisitedNodes)
				break
			}
			nextLoad = possibleLoads[0]
			loads = RemoveFromArray(loads, nextLoad.LoadNumber)
			// We have a node to add to our system
			visitedLoads[nextLoad.LoadNumber] = struct{}{}
			myVisitedNodes = append(myVisitedNodes, nextLoad.LoadNumber)
			currentCost += DistanceBetweenTwoPoints(startingLocation, nextLoad.Pickup) + nextLoad.Distance
			startingLocation = nextLoad.DropOff
		}
	}

	return
}

func implementation1LoadSwapPerformance(loads []*BasicLoad) (res [][]int) {
	minCost := math.MaxFloat64
	theIndex := 0
	loadMap := ConstructLoadMap(loads)
	for x, f := range LoadsOrderFuncs {
		clonedLoads := make([]*BasicLoad, len(loads))
		copy(clonedLoads, loads)
		tree := buildKDTree(clonedLoads)
		loadsOrderFunction = f
		reply := implementation1(tree, clonedLoads)
		cost := calculateCost(reply.Results, loadMap)
		if cost < minCost {
			theIndex = x
			minCost = cost
			res = reply.Results
		}
	}
	f, err := os.OpenFile("./output.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(fmt.Sprint(theIndex)); err != nil {
		panic(err)
	}
	return
}

var LoadsOrderFuncs = []func(loads []*BasicLoad){
	loadsOrderFunction,
	func(loads []*BasicLoad) {
		for i := range loads {
			j := rand.Intn(i + 1)
			loads[i], loads[j] = loads[j], loads[i]
		}
	},
	func(loads []*BasicLoad) {
		sort.Slice(loads, func(i, j int) bool { return loads[i].DistanceToHome < loads[j].DistanceToHome })
	},
	func(loads []*BasicLoad) {
		sort.Slice(loads, func(i, j int) bool { return loads[i].Distance > loads[j].Distance })
	},
	func(loads []*BasicLoad) {
		sort.Slice(loads, func(i, j int) bool { return loads[i].Distance < loads[j].Distance })
	},
}
