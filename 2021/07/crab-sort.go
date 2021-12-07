package main

import (
	"flag"
	"os"
	"sort"

	"github.com/aaronireland/advent-of-code/pkg/input"
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
)

var (
	filePath             string
	fuelUsageFuncEnabled bool
)

func init() {
	log.SetHandler(text.New(os.Stderr))
	flag.StringVar(&filePath, "input", "./input.txt", "The input data text file path")
	flag.BoolVar(&fuelUsageFuncEnabled, "fuel", false, "Use complex fuel usage function")
}

// fuelUsage computes fuel required to travel distance. Part One uses a simple calculation where one unit of fuel is
// required per unit of distance traveled. Part Two uses a recursive formula where each unit of distance requires one
// more unit of fuel than the prior distance.
func fuelUsage(distance int) int {
	if fuelUsageFuncEnabled {
		if distance < 2 {
			return distance
		}

		return distance + fuelUsage(distance-1)
	}
	return distance
}

func mostEfficientFuelUsage(crabs []int) int {
	sort.Ints(crabs)
	var (
		positionRange     = crabs[len(crabs)-1] - crabs[0]
		fuelUsageCh       = make(chan int, positionRange)
		mostEfficient int = -1
	)

	for pos := crabs[0]; pos <= crabs[len(crabs)-1]; pos++ {
		go func(c0 int, c []int, result chan int) {
			var fuelUsed int
			for _, cx := range c {
				if cx > c0 {
					fuelUsed += fuelUsage(cx - c0)
				} else {
					fuelUsed += fuelUsage(c0 - cx)
				}
			}
			result <- fuelUsed
		}(pos, crabs, fuelUsageCh)
	}

	for i := 0; i < positionRange; i++ {
		pos := <-fuelUsageCh
		if mostEfficient < 0 || pos < mostEfficient {
			mostEfficient = pos
		}
	}

	return mostEfficient

}

func main() {
	flag.Parse()

	crabs, err := input.CsvInts(filePath)
	if err != nil {
		log.WithError(err).Fatal("unable to read input file data")
	}
	bestPosition := mostEfficientFuelUsage(crabs)

	log.Infof("Minimum fuel required to sort crabs: %d", bestPosition)
}
