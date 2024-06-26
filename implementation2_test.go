package main

import (
	"fmt"
	"sort"
	"testing"
)

func TestImplementation2TruckRecursionX(t *testing.T) {
	loads, err := loadLoadFile("myTestProblems/FarTrip.txt")
	if err != nil {
		t.Error(err)
		return
	}
	sort.Slice(loads, func(i, j int) bool { return loads[i].DistanceToHome > loads[j].DistanceToHome })
	tree := buildKDTree(loads)
	firstNode := loads[0]
	distanceSoFar := firstNode.Distance + DistanceBetweenTwoPoints(HOME_VECTOR, firstNode.Pickup)
	visitedNodes := make(map[int]struct{})
	visitedNodes[firstNode.LoadNumber] = struct{}{}
	results := implementation2TruckRecursion(tree, firstNode.DropOff, distanceSoFar, visitedNodes)
	for x := range results {
		fmt.Println(results[x])
	}
}

func TestImplementation2TruckRecursionFleetLevel(t *testing.T) {
	loads, err := loadLoadFile("myTestProblems/CloseTrip.txt")
	if err != nil {
		t.Error(err)
		return
	}
	sort.Slice(loads, func(i, j int) bool { return loads[i].DistanceToHome > loads[j].DistanceToHome })
	tree := buildKDTree(loads)

	results := Implementation2PossibleNewTruckRoutes(tree, loads, [][]int{})
	for x := range results {
		fmt.Println(results[x])
	}
}

func TestImplementation2RecursionStartingPoint2(t *testing.T) {
	filePath := "./Training Problems/problem20.txt"
	// filePath := "myTestProblems/CloseTrip.txt"
	loads, err := loadLoadFile(filePath)
	if err != nil {
		t.Error(err)
		return
	}
	sort.Slice(loads, func(i, j int) bool { return loads[i].DistanceToHome > loads[j].DistanceToHome })
	tree := buildKDTree(loads)

	results := Implementation2RecursionStartingPoint2(tree, loads, PathsOfSomeTrucks{})
	for x := range results {
		fmt.Println(results[x])
	}
}

// Okay somehow even just setting it to 1 as the search extra, we use 1000 more milage points
func TestImplementation2RecursionStartingPoint2FromExistingPath(t *testing.T) {
	filePath := "./Training Problems/problem20.txt"
	// filePath := "myTestProblems/CloseTrip.txt"
	loads, err := loadLoadFile(filePath)
	if err != nil {
		t.Error(err)
		return
	}
	sort.Slice(loads, func(i, j int) bool { return loads[i].DistanceToHome > loads[j].DistanceToHome })
	tree := buildKDTree(loads)

	// With just {47,185,151} got [0 0 0] as one of the options ,{165, 146}
	results := Implementation2RecursionStartingPoint2(tree, loads, PathsOfSomeTrucks{{47, 185, 151}, {165, 146}, {23, 106}})
	loadMap := ConstructLoadMap(loads)
	sol, cost := calculateMinimumCost(results, loadMap)
	t.Log(sol)
	t.Log(cost)
}


func TestImp1LoadsOrderPerformance(t *testing.T){
	filePath := "./Training Problems/problem13.txt"
	loads, err := loadLoadFile(filePath)
	if err != nil {
		t.Error(err)
		return
	}
	sort.Slice(loads, func(i, j int) bool { return loads[i].DistanceToHome > loads[j].DistanceToHome })
	best := implementation1LoadSwapPerformance(loads)
	t.Log(best)
}