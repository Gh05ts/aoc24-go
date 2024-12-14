package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
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

	reader := fileReader("inputs/10.txt")
	defer reader.Close()

	grid, startLocations := ioConvert(reader)
	count := findScore(grid, startLocations)
	fmt.Println(count)
	// printGrid(grid)
	// printStartLocations(startLocations)
	// fmt.Println(grid)
}

func fileReader(path string) *bufioCloser {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(fmt.Errorf("open failed: %w", err).Error())
		panic("file panic")
	}

	return &bufioCloser{bufio.NewReader(file), file}
}

func ioConvert(reader *bufioCloser) ([][]int, []pair) {
	lines := [][]int{}
	startLocations := []pair{}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		lines, startLocations = lineConv(lines, line, startLocations)
	}
	return lines, startLocations
}

type pair struct {
	x, y int
}

func lineConv(lines [][]int, line string, startLocations []pair) ([][]int, []pair) {
	len := len(lines)
	lines = append(lines, []int{})
	for i, char := range line {
		value := atoiWrap(string(char))
		if value == 0 {
			startLocations = append(startLocations, pair{len, i})
		}
		lines[len] = append(lines[len], value)
	}
	return lines, startLocations
}

func atoiWrap(number string) int {
	num, err := strconv.Atoi(number)
	if err != nil {
		fmt.Println(fmt.Errorf("atoi failed: %w", err).Error())
		panic("atoi panic")
	}
	return num
}

func findScore(grid [][]int, startLocations []pair) int {
	count := 0
	hashMap := make(map[pair]int)
	startEnd := make(map[pair]map[pair]bool)

	LenLimit := len(grid)
	BreadthLimit := len(grid[0])
	for _, startLocation := range startLocations {
		dfs(grid, startLocation.x, startLocation.y, LenLimit, BreadthLimit, startLocation, hashMap, startEnd)
	}
	for _, v := range hashMap {
		count += v
	}
	return count
}

func dfs(grid [][]int, i, j, length, breadth int, startLocation pair, finalLocations map[pair]int, startEnd map[pair]map[pair]bool) {
	if grid[i][j] == 9 {
		finalLocations[pair{i, j}]++
	}

	directions := [4][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	for _, dir := range directions {
		ni, nj := i+dir[0], j+dir[1]
		if ni >= 0 && ni < length && nj >= 0 && nj < breadth && grid[ni][nj] == grid[i][j]+1 {
			dfs(grid, ni, nj, length, breadth, startLocation, finalLocations, startEnd)
		}
	}
}
