package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	searchWord  = "XMAS"
	searchWord2 = "MAS"
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

	reader := fileReader("inputs/4.txt")
	defer reader.Close()

	lines := ioConvert(reader)
	count := XmasCounter(lines, args[1])
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
	var lines [][]rune
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

func XmasCounter(lines [][]rune, stage string) int {
	var word string
	// var stageFlag bool
	if strings.Compare(stage, "v2") != 0 {
		// stageFlag = false
		word = searchWord
	} else {
		// stageFlag = true
		word = searchWord2
	}
	positions := findAllOccurences(lines, word)
	// if stageFlag {
	// 	positions = findIntersections(positions)
	// }
	return len(positions)
}

func findAllOccurences(grid [][]rune, word string) [][]int {
	var directions [][]int
	var stageFlag bool
	if strings.Compare(word, searchWord) == 0 {
		stageFlag = false
		directions = [][]int{
			{0, 1}, {1, 0}, {1, 1}, {1, -1},
			{0, -1}, {-1, 0}, {-1, -1}, {-1, 1},
		}
	} else {
		stageFlag = true
		directions = [][]int{
			{1, 1}, {1, -1}, {-1, -1}, {-1, 1},
		}
	}

	var positions [][]int
	seen := make(map[string]int)

	for i := range grid {
		for j := range grid[i] {
			for _, dir := range directions {
				if search(grid, word, i, j, dir) {
					if stageFlag {
						// position of char 'A'
						// storing in hashmap
						// multiple same A is a match
						position := []int{i + dir[0], j + dir[1]}
						key := arrayToString(position)
						if _, found := seen[key]; found {
							positions = append(positions, position)
							seen[key]++
						} else {
							seen[key] = 1
						}
					} else {
						positions = append(positions, []int{i, j})
					}
				}
			}
		}
	}

	return positions
}

func search(grid [][]rune, word string, row, col int, dir []int) bool {
	for k := 0; k < len(word); k++ {
		r := row + k*dir[0]
		c := col + k*dir[1]
		if r < 0 || r >= len(grid) || c < 0 || c >= len(grid[0]) || grid[r][c] != rune(word[k]) {
			return false
		}
	}
	return true
}

// func findIntersections(positions [][]int) [][]int {
// 	seen := make(map[string]int)
// 	intersections := [][]int{}
// 	for _, array := range positions {
// 		key := arrayToString(array)
// 		if _, found := seen[key]; found {
// 			intersections = append(intersections, array)
// 			seen[key]++
// 		} else {
// 			seen[key] = 1
// 		}
// 	}
// 	return intersections
// }

func arrayToString(position []int) string {
	return fmt.Sprint(position)
}
