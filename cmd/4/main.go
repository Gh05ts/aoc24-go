package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	searchWord = "XMAS"
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
	count := XmasCounter(lines)
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

func XmasCounter(lines [][]rune) int {
	positions := findAllOccurences(lines, searchWord)
	if len(positions) > 0 {
		fmt.Println("Word found!")
	} else {
		fmt.Println("Word not found.")
	}
	return len(positions)
}

func findAllOccurences(grid [][]rune, word string) [][]int {
	directions := [][]int{
		{0, 1}, {1, 0}, {1, 1}, {1, -1},
		{0, -1}, {-1, 0}, {-1, -1}, {-1, 1},
	}

	var positions [][]int

	for i := range grid {
		for j := range grid[i] {
			for _, dir := range directions {
				if search(grid, word, i, j, dir) {
					positions = append(positions, []int{i, j})
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
