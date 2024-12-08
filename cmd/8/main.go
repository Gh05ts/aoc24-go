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
	for _, locations := range runeLocations {
		for i := 0; i < len(locations); i++ {
			for j := i + 1; j < len(locations); j++ {
				addValidLocations(locations, i, j, height, width, uniqueLocations)
			}
		}
	}
	return len(uniqueLocations)
}

func addValidLocations(locations []pair, i, j, height, width int, uniqueLocations map[pair]bool) {
	pointA := locations[i]
	pointB := locations[j]

	var point1, point2 pair
	diffHeight := 0
	diffLen := 0
	if pointA.y > pointB.y { // diagonal
		diffHeight = pointB.x - pointA.x
		diffLen = pointA.y - pointB.y
		point1 = pair{pointA.x - diffHeight, pointA.y + diffLen}
		point2 = pair{pointB.x + diffHeight, pointB.y - diffLen}
	} else {
		diffHeight = pointB.x - pointA.x
		diffLen = pointB.y - pointA.y
		point1 = pair{pointA.x - diffHeight, pointA.y - diffLen}
		point2 = pair{pointB.x + diffHeight, pointB.y + diffLen}
	}
	if isValid(point1, height, width) {
		uniqueLocations[point1] = true
	}
	if isValid(point2, height, width) {
		uniqueLocations[point2] = true
	}
}

func isValid(point pair, len, wid int) bool {
	return point.x >= 0 && point.x < len && point.y >= 0 && point.y < wid
}

func diagUpRight() {

}

func diagDownLeft() {

}

func diagUpLeft() {

}

func diagDownRight() {

}
