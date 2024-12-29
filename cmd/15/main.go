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

	reader := fileReader("inputs/15.txt")
	defer reader.Close()

	grid, ops := ioConvert(reader)
	startLocation, _, _ := prepareInput(grid)

	flag := strings.Compare(args[1], "v2") == 0
	sum := getSum(grid, ops, startLocation, flag)
	fmt.Println(sum)
}

func fileReader(path string) *bufioCloser {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(fmt.Errorf("open failed: %w", err).Error())
		panic("file panic")
	}

	return &bufioCloser{bufio.NewReader(file), file}
}

type pair struct {
	x, y int
}

func ioConvert(reader *bufioCloser) (gridSlice [][]rune, opsArray []rune) {
	var grid [][]rune
	var ops []rune
	opsFlag := false
	i := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if opsFlag || line == "" {
			if !opsFlag {
				opsFlag = true
				continue
			}
			ops = append(ops, []rune(line)...)
		} else {
			grid = append(grid, []rune(line))
		}
		i++
	}
	return grid, ops
}

func prepareInput(grid [][]rune) (startLoc pair, boxes, walls []pair) {
	var startLocation pair
	boxSlice, wallsSlice := []pair{}, []pair{}
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == rune('O') {
				boxSlice = append(boxSlice, pair{i, j})
			} else if grid[i][j] == rune('#') {
				wallsSlice = append(wallsSlice, pair{i, j})
			} else if grid[i][j] == rune('@') {
				startLocation = pair{i, j}
			}
		}
	}
	return startLocation, boxSlice, wallsSlice
}

func getSum(grid [][]rune, ops []rune, startLocation pair, stage2 bool) int {
	count := 0
	dirs := map[rune]pair{rune('<'): {0, -1}, rune('>'): {0, 1}, rune('^'): {-1, 0}, rune('v'): {1, 0}}
	for _, char := range ops {
		dir := dirs[char]
		boxs := []pair{}
		move := false
		tracker := startLocation
		for {
			tracker.x += dir.x
			tracker.y += dir.y
			if grid[tracker.x][tracker.y] == rune('#') {
				break
			}
			if grid[tracker.x][tracker.y] == rune('.') {
				move = true
				break
			}
			if grid[tracker.x][tracker.y] == rune('O') {
				boxs = append(boxs, tracker)
			}
		}
		if move {
			grid[startLocation.x][startLocation.y] = '.'
			startLocation.x += dir.x
			startLocation.y += dir.y
			for _, location := range boxs {
				grid[location.x][location.y] = rune('.')
			}
			for _, location := range boxs {
				grid[location.x+dir.x][location.y+dir.y] = rune('O')
			}
			grid[startLocation.x][startLocation.y] = '@'
		}
	}
	for i := range grid {
		for j := range grid {
			if grid[i][j] == rune('O') {
				count += 100*i + j
			}
		}
	}
	return count
}

// func opPrinter(ops []rune) {
// 	for _, char := range ops {
// 		fmt.Print(string(char))
// 	}
// 	fmt.Println()
// }

// func gridPrinter(grid [][]rune) {
// 	for i := range grid {
// 		for j := range grid[i] {
// 			fmt.Print(string(grid[i][j]))
// 		}
// 		fmt.Println()
// 	}
// }
