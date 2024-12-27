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

	reader := fileReader("inputs/12.txt")
	defer reader.Close()

	grid := ioConvert(reader)
	flag := strings.Compare(args[1], "v2") == 0
	// printGrid(grid)
	price := totalPrice(grid, flag)
	fmt.Println(price)
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

func totalPrice(grid [][]rune, discount bool) int {
	price := 0
	visited := make(map[pair]bool)
	groups := [][]pair{}

	lenLimit := len(grid)
	breadthLimit := len(grid[0])
	for i := range grid {
		for j := range grid[i] {
			if !visited[pair{i, j}] {
				currentGroup := new([]pair)
				dfs(grid, i, j, lenLimit, breadthLimit, currentGroup, visited)
				groups = append(groups, *currentGroup)
			}
		}
	}

	for i := range groups {
		perimeter := 0
		corners := 0
		for _, val := range groups[i] {
			if discount {
				corners += countCorners(groups[i], val)
			} else {
				perimeter += getPerimeter(grid, val, lenLimit, breadthLimit)
			}
		}
		// fmt.Println("i: ", i, "area:", len(groups[i]), "perimeter:", perimeter, "corners:", corners)
		if discount {
			price += len(groups[i]) * corners
		} else {
			price += len(groups[i]) * perimeter
		}
	}

	return price
}

func dfs(grid [][]rune, i, j, length, breadth int, groups *[]pair, visited map[pair]bool) {
	visited[pair{i, j}] = true
	*groups = append(*groups, pair{i, j})

	directions := [4][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	for _, dir := range directions {
		ni, nj := i+dir[0], j+dir[1]
		if ni >= 0 && ni < length && nj >= 0 && nj < breadth && grid[ni][nj] == grid[i][j] && !visited[pair{ni, nj}] {
			dfs(grid, ni, nj, length, breadth, groups, visited)
		}
	}
}

func getPerimeter(grid [][]rune, val pair, lenLimit, breadthLimit int) int {
	perimeter := 0
	fourDirections := [4][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	curChar := grid[val.x][val.y]
	for _, direction := range fourDirections {
		ni, nj := val.x+direction[0], val.y+direction[1]
		if ni >= 0 && ni < lenLimit && nj >= 0 && nj < breadthLimit && grid[ni][nj] != curChar {
			perimeter += 1
		} else if ni < 0 || nj < 0 || ni == lenLimit || nj == breadthLimit {
			perimeter += 1
		}
	}
	return perimeter
}

func countCorners(group []pair, cur pair) int {
	corners := 0
	sevenDirection := [7][2]int{{-1, 0}, {0, -1}, {-1, -1}, {1, 0}, {1, -1}, {0, 1}, {-1, 1}}

	x, y := cur.x, cur.y
	north := pair{x + sevenDirection[0][0], y + sevenDirection[0][1]}
	west := pair{x + sevenDirection[1][0], y + sevenDirection[1][1]}
	northWest := pair{x + sevenDirection[2][0], y + sevenDirection[2][1]}
	south := pair{x + sevenDirection[3][0], y + sevenDirection[3][1]}
	southWest := pair{x + sevenDirection[4][0], y + sevenDirection[4][1]}
	east := pair{x + sevenDirection[5][0], y + sevenDirection[5][1]}
	northEast := pair{x + sevenDirection[6][0], y + sevenDirection[6][1]}

	if !search(group, north) {
		sameEdge := search(group, west) && !search(group, northWest)
		if !sameEdge {
			corners += 1
		}
	}

	if !search(group, south) {
		sameEdge := search(group, west) && !search(group, southWest)
		if !sameEdge {
			corners += 1
		}
	}

	if !search(group, west) {
		sameEdge := search(group, north) && !search(group, northWest)
		if !sameEdge {
			corners += 1
		}
	}

	if !search(group, east) {
		sameEdge := search(group, north) && !search(group, northEast)
		if !sameEdge {
			corners += 1
		}
	}
	return corners
}

func search(haystack []pair, needle pair) bool {
	for _, hay := range haystack {
		if needle == hay {
			return true
		}
	}
	return false
}

// func printGrid(grid [][]rune) {
// 	for i, _ := range grid {
// 		for j, _ := range grid[i] {
// 			fmt.Print(string(grid[i][j]))
// 		}
// 		fmt.Println()
// 	}
// }

// func printGroups(grid [][]rune, groups [][]pair) {
// 	for i := range groups {
// 		for _, val := range groups[i] {
// 			fmt.Print(string(grid[val.x][val.y]))
// 		}
// 		fmt.Println()
// 	}
// }
