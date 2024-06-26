package main

import (
	"testing"

	"github.com/downflux/go-geometry/nd/vector"
)

func TestLoadTestLoad(t *testing.T) {
	filePath := "./Training Problems/problem1.txt"
	loads, err := loadLoadFile(filePath)
	if err != nil {
		t.Error(err)
	}

	if len(loads) != 10 {
		t.Error("recieved wrong number of loads", loads)
	}
}

func TestSubsetLength(t *testing.T){
	// Give the points and lets calculate the total distance
	totalDistance := float64(0)
	vectors := []vector.V{}
	vectors = append(vectors, *HOME_VECTOR)
	// 87
	vectors = append(vectors, *vector.New(221.1650211140732, 258.03443664251506), *vector.New(236.4917656338128,209.45076030564394))
	// 23
	vectors = append(vectors, *vector.New(237.57608273227433,205.13238611597814), *vector.New(180.30514919654038,151.8564521983921))
	// 153
	vectors = append(vectors, *vector.New(116.52141569710801,156.84704344824164), *vector.New(86.4454664122545, 94.99201923066752))

	vectors = append(vectors, *HOME_VECTOR)

	for x := 0; x < len(vectors) - 1; x ++ {
		totalDistance += DistanceBetweenTwoPoints(&vectors[x], &vectors[x+1])
	}
	t.Log(totalDistance)
}