package main

import (
	"encoding/csv"
	"errors"	
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/downflux/go-geometry/nd/vector"
)

const MAX_DISTANCE = float64(720)
var HOME_VECTOR = vector.New(0,0)

var ErrCSVWrongLength = errors.New("CSV had wrong number of entries for line")
var ErrCSVWrongPointsLength = errors.New("CSV had wrong number of entries for a point value")
var ErrInvalidLoadDistance = errors.New("load by default is too long")

// Bad Naming
type Result struct {
	CalculatedCost float64
	NumberOfDrivers int
	Results Results
}

type BasicLoad struct {
	LoadNumber     int // could be char
	Pickup         *vector.V
	DropOff        *vector.V
	Distance       float64
	DistanceToHome float64 // Maybe not needed, but the distance from the end of this load, to get back home. If its too big, we know we can't add it
	// don't have distance to start because ideally we are coming from a different load
	sync.Mutex // Can lock a basicLoad if we are planning on using it, for multi-threaded solution. 
}

// Assuming how locations are sorted, so returning the starting locale
func (bl *BasicLoad) P() vector.V {return *bl.Pickup}
func (bl *BasicLoad) Equal(q *BasicLoad) bool { return bl.LoadNumber == q.LoadNumber}

type Point struct {
	X float64
	Y float64
}

func loadLoadFile(filePath string) (loads []*BasicLoad, err error) {
	loads = make([]*BasicLoad, 0, 10) // All the test loads have at least 10, so going allocate 10 to start
	f, err := os.Open(filePath)
	if err != nil {
		return
	}

	// https://pkg.go.dev/encoding/csv#example-Reader
	csvReader := csv.NewReader(f)
	csvReader.Comma = ' '
	// Removing title field
	_, err = csvReader.Read()
	if err != nil {
		return
	}

	for {
		record := BasicLoad{}
		var line []string
		line, err = csvReader.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return
		}
		if len(line) != 3 {
			err = ErrCSVWrongLength
			return
		}
		record.LoadNumber, err = strconv.Atoi(line[0])
		if err != nil {
			return
		}
		record.Pickup, err = parsePoint(line[1])
		if err != nil {
			return
		}
		record.DropOff, err = parsePoint(line[2])
		if err != nil {
			return
		}

		record.Distance = DistanceBetweenTwoPoints(record.Pickup, record.DropOff)
		record.DistanceToHome = DistanceBetweenTwoPoints(record.DropOff, HOME_VECTOR)
		loads = append(loads, &record)
	}

	return
}

func parsePoint(pointString string) (point *vector.V, err error) {
	points := strings.Split(strings.TrimSuffix(strings.TrimPrefix(pointString, "("), ")"), ",")
	if len(points) != 2 {
		err = ErrCSVWrongPointsLength
		return
	}
	
	// point.X, err = strconv.ParseFloat(points[0], 64)
	X, err := strconv.ParseFloat(points[0], 64)
	if err != nil {
		return
	}

	// point.Y, err = strconv.ParseFloat(points[1], 64)
	Y, err := strconv.ParseFloat(points[1], 64)
	if err != nil {
		return
	}

	point = vector.New(X, Y)
	return
}

func DistanceBetweenTwoPoints(a, b *vector.V) (distance float64) {
	xDif := a.X(0) - b.X(0)
	yDif := a.X(1) - b.X(1)
	distance = math.Sqrt((xDif * xDif) + (yDif * yDif))
	return
}

// Prompt says the actions are always possible
func isImpossible(load *BasicLoad) (impossible bool) {
	return (load.Distance + load.DistanceToHome + DistanceBetweenTwoPoints(HOME_VECTOR, load.Pickup)) > MAX_DISTANCE
}

func CalculateSolutionScore(numDrivers int, totalDistance float64) (score float64) {
	score = (float64(numDrivers) * 500) + totalDistance
	return
}
