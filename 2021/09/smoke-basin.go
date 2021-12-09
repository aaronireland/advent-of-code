package main

import (
	"flag"
	"os"
	"sort"
	"strconv"

	"github.com/aaronireland/advent-of-code/pkg/input"
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
)

var filePath string

func init() {
	flag.StringVar(&filePath, "input", "./input.txt", "The path to the input data file")
	log.SetHandler(text.New(os.Stderr))
}

type gridMap [][]uint8

func newGrid(data []string) (grid gridMap, err error) {
	for _, line := range data {
		var row []uint8
		for _, char := range line {
			height, err := strconv.ParseInt(string(char), 10, 8)
			if err != nil {
				return grid, err
			}
			row = append(row, uint8(height))
		}
		grid = append(grid, row)
	}
	return grid, nil
}

func (g gridMap) IsLowestPoint(i, j int) bool {
	var up, down, left, right bool
	up = (i == 0) || g[i-1][j] > g[i][j]
	down = (i == len(g)-1) || g[i+1][j] > g[i][j]
	left = (j == 0) || g[i][j-1] > g[i][j]
	right = (j == len(g[i])-1) || g[i][j+1] > g[i][j]

	return up && down && left && right
}

func (g gridMap) RiskLevel(i, j int) int {
	if g.IsLowestPoint(i, j) {
		return int(g[i][j]) + 1
	}
	return 0
}

func (g gridMap) BasinSize(i, j int) int {
	var (
		calculateBasinArea func(*[][]int, []int) int
		nodesVisited       [][]int
	)

	// Recursively cycle through the coordinates and count unseen nodes that have a height value less than 9. Ensure that
	// each copy of the function is updating the same visited nodes array but counting nodes inside the closure
	calculateBasinArea = func(visited *[][]int, node []int) int {
		i, j := node[0], node[1]
		size := 0
		if g.exists(i, j) && g[i][j] < 9 && !seen(*visited, node) {
			size++
		}
		*visited = append(*visited, []int{i, j})

		up, down, left, right := []int{-1, 0}, []int{1, 0}, []int{0, -1}, []int{0, 1}
		directions := [][]int{left, right, up, down}

		for _, direction := range directions {
			nextRow, nextColumn := i+direction[0], j+direction[1]
			if g.exists(nextRow, nextColumn) && !seen(*visited, []int{nextRow, nextColumn}) && g[nextRow][nextColumn] < 9 {
				size += calculateBasinArea(visited, []int{nextRow, nextColumn})
			}
		}
		// There's nowhere left to go for this vector!
		return size
	}
	return calculateBasinArea(&nodesVisited, []int{i, j})
}

func seen(visited [][]int, node []int) bool {
	for _, vNode := range visited {
		if vNode[0] == node[0] && vNode[1] == node[1] {
			return true
		}
	}
	return false
}

func (g gridMap) exists(i, j int) bool {
	rowOnMap := i >= 0 && i < len(g)
	if rowOnMap {
		return j >= 0 && j < len(g[i])
	}
	return false
}

func main() {
	flag.Parse()
	data, err := input.ReadLines(filePath)
	if err != nil {
		log.WithError(err).Fatal("unable to read input data file")
	}

	topology, err := newGrid(data)
	if err != nil {
		log.WithError(err).Fatal("unable to generate topological map")
	}

	var (
		totalRisk int
		lowPoints [][]int
	)
	for i, row := range topology {
		for j := 0; j < len(row); j++ {
			risk := topology.RiskLevel(i, j)
			totalRisk += risk

			if risk > 0 {
				lowPoints = append(lowPoints, []int{i, j})
			}
		}
	}

	log.Infof("Total Risk is %d", totalRisk)

	var basinSizes []int
	for _, pt := range lowPoints {
		basin := topology.BasinSize(pt[0], pt[1])
		basinSizes = append(basinSizes, basin)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(basinSizes)))

	log.Infof("Area of three largest basins is: %d", basinSizes[0]*basinSizes[1]*basinSizes[2])
}
