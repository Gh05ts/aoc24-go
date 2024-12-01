package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Since bufioReader has no close on it
// creating this type to help close the file in the calling function
type bufioCloser struct {
	*bufio.Reader
	io.Closer
}

/*
input is in the format:
2 8
6 3
3 1

we store the columns in input as lists and sort them
2 1
3 3
6 8

then we traverse and get absolute diff and add them up
*/
func main() {
	reader := fileReader("inputs/1.txt")
	defer reader.Close()

	column1, column2 := ioColumns(reader)
	sort.Ints(column1)
	sort.Ints(column2)

	result := 0
	for i := range column1 {
		result += absIntDiff(column1[i], column2[i])
	}

	fmt.Println(result)
}

// allow main function to be just logical ops and move other things to the side
func fileReader(path string) *bufioCloser {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(fmt.Errorf("open failed: %w", err).Error())
		panic("file panic")
	}

	return &bufioCloser{bufio.NewReader(file), file}
}

// read lines and get fields and append to 2 lists
// first list corresponds to first column's values
// second list corresponds to second column's values
// convert to numbers from strings for math operations
func ioColumns(reader *bufioCloser) (column1, column2 []int) {
	var col1, col2 []int
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		numbers := strings.Fields(line)

		num1 := atoiWrap(numbers[0])
		col1 = append(col1, num1)

		num2 := atoiWrap(numbers[1])
		col2 = append(col2, num2)
	}
	return col1, col2
}

func atoiWrap(number string) int {
	num, err := strconv.Atoi(number)
	if err != nil {
		fmt.Println(fmt.Errorf("atoi failed: %w", err).Error())
		panic("atoi panic")
	}
	return num
}

func absIntDiff(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}
