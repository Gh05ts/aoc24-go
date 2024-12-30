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
	startLocation, _, _, gridWide, boxIDs := prepareInput(grid)

	stage2 := strings.Compare(args[1], "v2") == 0
	if stage2 {
		sum := getSumWide(gridWide, ops, startLocation, boxIDs)
		fmt.Println(sum)
	} else {
		sum := getSum(grid, ops, startLocation)
		fmt.Println(sum)
	}
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

func prepareInput(grid [][]rune) (startLoc pair, boxes, walls []pair, gridWide [][]rune, boxID [][]int) {
	var startLocation pair
	boxSlice, wallsSlice := []pair{}, []pair{}
	gridUpdated := [][]rune{}
	boxId := [][]int{}
	nextID := 1
	for i := range grid {
		row := []rune{}
		rowID := []int{}
		for j := range grid[i] {
			if grid[i][j] == rune('O') {
				row = append(row, 'O')
				row = append(row, 'O')
				id := nextID
				rowID = append(rowID, id)
				rowID = append(rowID, id)
				nextID++
				boxSlice = append(boxSlice, pair{i, j})
			} else if grid[i][j] == rune('#') {
				row = append(row, '#')
				row = append(row, '#')
				rowID = append(rowID, 0)
				rowID = append(rowID, 0)
				wallsSlice = append(wallsSlice, pair{i, j})
			} else if grid[i][j] == rune('@') {
				row = append(row, '@')
				row = append(row, '.')
				rowID = append(rowID, 0)
				rowID = append(rowID, 0)
				startLocation = pair{i, j}
			} else {
				rowID = append(rowID, 0)
				rowID = append(rowID, 0)
				row = append(row, '.')
				row = append(row, '.')
			}
		}
		gridUpdated = append(gridUpdated, row)
		boxId = append(boxId, rowID)
	}
	// gridPrinter(gridUpdated)
	return startLocation, boxSlice, wallsSlice, gridUpdated, boxId
}

func getSum(grid [][]rune, ops []rune, startLocation pair) int {
	sum := 0
	dirs := map[rune]pair{
		rune('<'): {0, -1},
		rune('>'): {0, 1},
		rune('^'): {-1, 0},
		rune('v'): {1, 0},
	}

	for _, char := range ops {
		dir := dirs[char]
		boxes := []pair{}
		move := false
		crawler := startLocation
		for {
			crawler.x += dir.x
			crawler.y += dir.y
			if grid[crawler.x][crawler.y] == rune('#') {
				break
			}
			if grid[crawler.x][crawler.y] == rune('.') {
				move = true
				break
			}
			if grid[crawler.x][crawler.y] == rune('O') {
				boxes = append(boxes, crawler)
			}
		}

		if move {
			grid[startLocation.x][startLocation.y] = rune('.')
			startLocation.x += dir.x
			startLocation.y += dir.y
			for _, location := range boxes {
				grid[location.x][location.y] = rune('.')
			}
			for _, location := range boxes {
				grid[location.x+dir.x][location.y+dir.y] = rune('O')
			}
			grid[startLocation.x][startLocation.y] = rune('@')
		}
	}

	for i := range grid {
		for j := range grid {
			if grid[i][j] == rune('O') {
				sum += 100*i + j
			}
		}
	}
	return sum
}

func getSumWide(grid [][]rune, ops []rune, startLocation pair, boxIDs [][]int) int {
	sum := 0
	dirs := map[rune]pair{
		rune('<'): {0, -1},
		rune('>'): {0, 1},
		rune('^'): {-1, 0},
		rune('v'): {1, 0},
	}
	startLocation.y *= 2

	for _, char := range ops {
		dir := dirs[char]
		boxes := []pair{}
		visited := map[int]bool{}
		move := true
		boxes = append(boxes, startLocation)
		for i := 0; i < len(boxes); i++ {
			row := boxes[i].x + dir.x
			col := boxes[i].y + dir.y
			if grid[row][col] == rune('#') {
				move = false
				break
			}
			if grid[row][col] == rune('O') {
				if _, found := visited[boxIDs[row][col]]; !found {
					visited[boxIDs[row][col]] = true
					boxes = append(boxes, pair{row, col})
					for _, dir := range []int{col + 1, col - 1} {
						if boxIDs[row][col] == boxIDs[row][dir] {
							boxes = append(boxes, pair{row, dir})
						}
					}
				}
			}
		}
		tmp := map[pair]int{}
		if move {
			grid[startLocation.x][startLocation.y] = rune('.')
			startLocation.x += dir.x
			startLocation.y += dir.y
			boxes = boxes[1:]
			for _, location := range boxes {
				grid[location.x][location.y] = rune('.')
				tmp[location] = boxIDs[location.x][location.y]
				boxIDs[location.x][location.y] = 0
			}
			for _, location := range boxes {
				boxIDs[location.x+dir.x][location.y+dir.y] = tmp[location]
				grid[location.x+dir.x][location.y+dir.y] = rune('O')
			}
			grid[startLocation.x][startLocation.y] = rune('@')
		}
	}

	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == rune('O') && boxIDs[i][j] != boxIDs[i][j-1] {
				grid[i][j] = rune('[')
				grid[i][j+1] = rune(']')
				sum += 100*i + j
			}
		}
	}
	gridPrinter(grid)
	return sum
}

// func opPrinter(ops []rune) {
// 	for _, char := range ops {
// 		fmt.Print(string(char))
// 	}
// 	fmt.Println()
// }

func gridPrinter(grid [][]rune) {
	for i := range grid {
		for j := range grid[i] {
			fmt.Print(string(grid[i][j]))
		}
		fmt.Println()
	}
}
