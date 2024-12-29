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

const (
	width  = 101
	height = 103
)

func main() {
	args := os.Args

	if len(args) < 2 {
		args = append(args, "v1")
	}

	reader := fileReader("inputs/14.txt")
	defer reader.Close()

	robots := ioConvert(reader)
	countTurn := strings.Compare(args[1], "v2") == 0

	if !countTurn {
		safetyFactor := safetyFactor(robots)
		fmt.Println(safetyFactor)
	} else {
		turn := findTree(robots)
		fmt.Println(turn)
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
	x, y   int
	vx, vy int
}

func ioConvert(reader *bufioCloser) []pair {
	robots := []pair{}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)

		var x, y, vx, vy int
		fmt.Sscanf(line, "p=%d,%d v=%d,%d", &x, &y, &vx, &vy)
		robots = append(robots, pair{x, y, vx, vy})
	}
	return robots
}

func safetyFactor(io []pair) int {
	collector := [2][2]int{}
	for _, val := range io {
		px := ((val.x+100*val.vx)%width + width) % width
		py := ((val.y+100*val.vy)%height + height) % height
		if px != width/2 && py != height/2 {
			if px > width/2 && py > height/2 {
				collector[1][1]++
			} else if px > width/2 && py < height/2 {
				collector[1][0]++
			} else if px < width/2 && py > height/2 {
				collector[0][1]++
			} else {
				collector[0][0]++
			}
		}
	}
	return collector[0][0] * collector[0][1] * collector[1][0] * collector[1][1]
}

func findTree(io []pair) int {
	turn := 0
	for i := 0; i < 10000; i++ {
		positions := [width][height]rune{}
		for _, val := range io {
			px := ((val.x+i*val.vx)%width + width) % width
			py := ((val.y+i*val.vy)%height + height) % height
			positions[px][py] = 'X'
		}

		count := 0
		for i := range positions {
			for j := range positions[i] {
				if j < width-1-j && positions[i][j] == 'X' && positions[i][width-1-j] == 'X' {
					count++
				}
			}
		}

		if count > 40 {
			printGrid(positions)
			return turn
		}
		turn++
	}
	return turn
}

func printGrid(grid [width][height]rune) {
	lg := len(grid)
	lh := len(grid[0])

	for j := 0; j < lh; j++ {
		for i := 0; i < lg; i++ {
			if grid[i][j] == rune('X') {
				fmt.Print(string('*')) // grid[i][j]
			} else {
				fmt.Print(string(' '))
			}
		}
		fmt.Println()
	}
	fmt.Println()

	// for i := range grid {
	// 	for j := range grid[i] {
	// 		if grid[i][j] == rune('X') {
	// 			fmt.Print(string(grid[i][j]))
	// 		} else {
	// 			fmt.Print(string(' '))
	// 		}
	// 	}
	// 	fmt.Println()
	// }
	// fmt.Println()
}
