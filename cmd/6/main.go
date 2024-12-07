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

	reader := fileReader("inputs/6.txt")
	defer reader.Close()

	grid := ioConvert(reader)
	count := count(grid, args[1])
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

func count(grid [][]rune, stage string) int {
	startingLocation := [2]int{}
	breakFlag := false
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == rune('^') {
				breakFlag = true
				startingLocation[0] = i
				startingLocation[1] = j
				grid[i][j] = '.'
				break
			}
		}
		if breakFlag {
			break
		}
	}
	visited, _ := traverse(grid, startingLocation)
	if strings.Compare(stage, "v2") == 0 {
		// start := time.Now()
		count := loopCheck(grid, startingLocation, visited)
		// end := time.Now()
		// fmt.Println("time taken:", end.Sub(start))
		return count
	} else {
		return len(visited)
	}
}

type pair struct {
	x, y int
}

type pairDir struct {
	x, y, dir int
}

func loopCheck(grid [][]rune, startingLoc [2]int, visited map[pair]bool) int {
	count := 0
	for location := range visited {
		grid[location.x][location.y] = '#'
		_, isLoop := traverse(grid, startingLoc)
		if isLoop {
			count++
		}
		grid[location.x][location.y] = '.'
	}
	return count
}

func traverse(grid [][]rune, startingLoc [2]int) (map[pair]bool, bool) {
	directions := [4][2]int{
		{-1, 0}, {0, 1}, {1, 0}, {0, -1},
	}
	directionIdx := 0
	x, y := startingLoc[0], startingLoc[1]
	visited := make(map[pair]bool)
	visited[pair{x, y}] = true

	// used for checking cycles
	visitedDir := make(map[pairDir]bool)
	visitedDir[pairDir{x, y, directionIdx}] = true
	for {
		nextX, nextY := x+directions[directionIdx][0], y+directions[directionIdx][1]
		if _, found := visitedDir[pairDir{nextX, nextY, directionIdx}]; found {
			return visited, true
		}
		if !(nextX >= 0 && nextX < len(grid[0]) && nextY >= 0 && nextY < len(grid)) {
			break
		} else if grid[nextX][nextY] == rune('.') {
			x, y = nextX, nextY
			visited[pair{x, y}] = true
			visitedDir[pairDir{x, y, directionIdx}] = true
		} else {
			directionIdx = (directionIdx + 1) % len(directions)
		}
	}

	return visited, false
}
