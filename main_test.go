package main

import (
	"sort"
	"testing"
)

// Test my different implementation counts

func TestImplementation1HyperSpace(t *testing.T){
	filePath := "./Training Problems/problem20.txt"
	// filePath := "myTestProblems/CloseTrip.txt"
	loads, err := loadLoadFile(filePath)
	if err != nil {
		t.Error(err)
		return
	}
	sort.Slice(loads, func(i, j int) bool { return loads[i].DistanceToHome > loads[j].DistanceToHome })
	loadMap := ConstructLoadMap(loads)
	// With just {47,185,151} got [0 0 0] as one of the options ,{165, 146}
	results := implementation1(loads)
	
	t.Log(results)
	cost := calculateCost(results, loadMap)
	t.Log(results)
	t.Log(cost)
}

func TestImplementation1(t *testing.T) {
	results, err := BasicRun("./Training Problems/problem10.txt", implementation1)
	if err != nil {
		t.Error(err)
	}

	ensureNoRepeatNodes(t, results)
}

func ensureNoRepeatNodes(t *testing.T, results Results) {
	// can also do this by keeping track of min, max and so far sum of all the nodes we counted
	// then comparing to the sum of all numbers between 1 and max, <- this doesn't actually work, in like a most specific possibility, 
	mappable := make(map[int]struct{})
	min := 100
	max := 0
	total := 0
	for _, nodeArray := range results {
		for _, v := range nodeArray {
			total += v
			if v > max {
				max = v
			}
			if min > v {
				min = v
			}
			_, ok := mappable[v]
			if ok {
				t.Fail()
				return
			}
			mappable[v] = struct{}{}
		}
	}
	if min != 1 {
		t.Fail()
	}

	if len(mappable) != max {
		t.Log(len(mappable), " ", max)
	}

	maxFloat := float64(max)
	rb := int((maxFloat * (maxFloat + 1)) / 2)
	if total != rb {
		t.Fail()
	}

}


func TestImplementation4Basic(t *testing.T) {
	results, err := BasicRun("./Training Problems/problem10.txt", implementation4)
	if err != nil {
		t.Error(err)
	}

	ensureNoRepeatNodes(t, results)
}