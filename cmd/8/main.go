package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type bufioCloser struct {
	*bufio.Reader
	io.Closer
}

func main() {
	args := os.Args

	if len(args) < 2 {
		args = append(args, "v1")
	}

	reader := fileReader("inputs/8.txt")
	defer reader.Close()

	grid := ioConvert(reader)
	count := findAllLocations(grid, args[1])
	fmt.Println(count)
}

func fileReader(path string) *bufioCloser {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(fmt.Errorf("open failed: %w", err).Error())
		panic("file panic")
	}

	return &bufioCloser{bufio.NewReader(file), file}
}

func ioConvert(reader *bufioCloser) [][]rune {
	lines := [][]rune{}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		lines = append(lines, []rune(line))
	}
	return lines
}

type pair struct {
	x, y int
}

func findAllLocations(grid [][]rune, stage string) int {
	runeLocations := make(map[rune][]pair)
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] != rune('.') {
				if _, found := runeLocations[grid[i][j]]; !found {
					runeLocations[grid[i][j]] = []pair{}
				}
				runeLocations[grid[i][j]] = append(runeLocations[grid[i][j]], pair{i, j})
			}
		}
	}

	height := len(grid)
	width := len(grid[0])
	uniqueLocations := make(map[pair]bool)
	recursiveFlag := strings.Compare(stage, "v2") == 0
	for _, locations := range runeLocations {
		for i := 0; i < len(locations); i++ {
			for j := i + 1; j < len(locations); j++ {
				addAllValidLocationsForPoints(locations, i, j, height, width, uniqueLocations, recursiveFlag)
			}
		}
	}
	return len(uniqueLocations)
}

// struct to encapsulate point with all the info for operations on it
// opX, opY are the signs for add/ sub with 0 as sub and 1 as add
// pX and pY are the values to be added to point.x and point.y
// bX and bY are the bounds
type pairOp struct {
	point    pair
	opX, opY int
	pX, pY   int
	bX, bY   int
}

func addAllValidLocationsForPoints(locations []pair, i, j, height, width int, uniqueLocations map[pair]bool, loop bool) {
	pointA := locations[i]
	pointB := locations[j]

	// in stage 2, even the original location is considered an antitode
	if loop {
		uniqueLocations[pointA] = true
		uniqueLocations[pointB] = true
	}

	diffHeight := 0
	diffLen := 0
	if pointA.y > pointB.y { // +ve slope diagonal
		diffHeight = pointB.x - pointA.x
		diffLen = pointA.y - pointB.y
		populateDiagonalWithOp(pairOp{pointA, 0, 1, diffHeight, diffLen, height, width}, loop, uniqueLocations)
		populateDiagonalWithOp(pairOp{pointB, 1, 0, diffHeight, diffLen, height, width}, loop, uniqueLocations)
	} else {
		diffHeight = pointB.x - pointA.x
		diffLen = pointB.y - pointA.y
		populateDiagonalWithOp(pairOp{pointA, 0, 0, diffHeight, diffLen, height, width}, loop, uniqueLocations)
		populateDiagonalWithOp(pairOp{pointB, 1, 1, diffHeight, diffLen, height, width}, loop, uniqueLocations)
	}
}

// if loop that means stage 2 search exhaustively otherwise break preemptively
// the point inside the datastructure pairOp is updated to mark moving to next point
func populateDiagonalWithOp(po pairOp, loop bool, uniqueLocations map[pair]bool) {
	for {
		nextPoint := ApplyOp(po)
		if !isValid(nextPoint, po.bX, po.bY) {
			break
		}
		po.point = nextPoint
		uniqueLocations[nextPoint] = true
		if loop {
			continue
		} else {
			break
		}
	}
}

// function for abstracting adding to diagonals
// 0 means - and 1 means +
// opX goes to inp.x and opY goes to inp.y
// pX and pY are the parameters to be applied to inp.x and inp.y respectively
func ApplyOp(po pairOp) pair {
	if po.opX == 0 && po.opY == 0 {
		return pair{po.point.x - po.pX, po.point.y - po.pY}
	} else if po.opX == 0 {
		return pair{po.point.x - po.pX, po.point.y + po.pY}
	} else if po.opX == 1 && po.opY == 1 {
		return pair{po.point.x + po.pX, po.point.y + po.pY}
	} else {
		return pair{po.point.x + po.pX, po.point.y - po.pY}
	}
}

func isValid(point pair, len, wid int) bool {
	return point.x >= 0 && point.x < len && point.y >= 0 && point.y < wid
}
