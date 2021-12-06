package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/apex/log"
)

var days int

func read_input(file_path string) (fish []int64) {
	file, err := os.Open(file_path)
	if err != nil {
		log.WithError(err).Fatal("cannot read input.txt")
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ",")
		for _, f := range line {
			fInt, err := strconv.ParseInt(f, 10, 64)
			if err == nil {
				fish = append(fish, fInt)
			}
		}
	}
	return fish
}

func grow_fish(fish int, days int, populationCh chan int64) {
	var totalCount int64
	var population = make(map[int]int64)
	population[fish] += 1
	for day := 0; day < days; day++ {
		var newFish int64 = population[0]
		population[0] = population[1]
		population[1] = population[2]
		population[2] = population[3]
		population[3] = population[4]
		population[4] = population[5]
		population[5] = population[6]
		population[6] = population[7] + newFish
		population[7] = population[8]
		population[8] = newFish
	}
	for _, populationCount := range population {
		totalCount += populationCount
	}
	populationCh <- totalCount
}

func main() {
	flag.IntVar(&days, "days", 0, "number of days to grow population")
	flag.Parse()
	population := read_input("./input.txt")
	var totalLanternFish int64
	popGrowthCh := make(chan int64)
	for _, f := range population {
		go grow_fish(int(f), days, popGrowthCh)
	}

	for i := 0; i < len(population); i++ {
		totalLanternFish += <-popGrowthCh
	}

	fmt.Printf("Total Lantern Fish Population is: %d\n", totalLanternFish)
}
