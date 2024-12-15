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

	fmt.Println(string('^'))

	grid, ops := ioConvert(reader)
	startLocation, boxes, walls := prepareInput(grid)

	flag := strings.Compare(args[1], "v2") == 0
	sum := getSum(grid, ops, startLocation, boxes, walls, flag)
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

func getSum(grid [][]rune, ops []rune, startLocation pair, boxes, walls []pair, stage2 bool) int {
	count := 0
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
