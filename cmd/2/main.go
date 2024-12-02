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
		return
	}

	reader := fileReader("inputs/2.txt")
	defer reader.Close()

	fmt.Println(ioCount(reader, args[1]))
}

func fileReader(path string) *bufioCloser {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(fmt.Errorf("open failed: %w", err).Error())
		panic("file panic")
	}

	return &bufioCloser{bufio.NewReader(file), file}
}

func ioCount(reader *bufioCloser, strict string) int {
	count := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		numbers := strings.Fields(line)

		// v1 deprecated
		// if (atoiWrap(numbers[0]) < atoiWrap(numbers[1]) && isSortedAsc(numbers)) || (atoiWrap(numbers[0]) > atoiWrap(numbers[1]) && isSortedDesc(numbers)) {
		// 	count += 1
		// }

		if canBeSortedByRemovingOneAscending(numbers, strict) || canBeSortedByRemovingOneDescending(numbers, strict) {
			count += 1
		}

	}
	return count
}

func canBeSortedByRemovingOneAscending(arr []string, strict string) bool {
	violations := 0
	index := -1
	for i := 0; i < len(arr)-1; i++ {
		one := atoiWrap(arr[i])
		two := atoiWrap(arr[i+1])
		if one > two || absIntDiff(one, two) < 1 || absIntDiff(one, two) > 3 {
			violations++
			index = i
			if violations > 1 {
				return false
			}
		}
	}
	// If no violations, already sorted
	if violations == 0 {
		return true
	}
	// Check by removing arr[index] or arr[index+1]
	return ((index == 0 || arr[index-1] <= arr[index+1]) || (index+1 == len(arr)-1 || arr[index] <= arr[index+2])) && strings.Compare(strict, "strict") != 0
}

func canBeSortedByRemovingOneDescending(arr []string, strict string) bool {
	violations := 0
	index := -1
	for i := 0; i < len(arr)-1; i++ {
		one := atoiWrap(arr[i])
		two := atoiWrap(arr[i+1])
		if one < two || absIntDiff(one, two) < 1 || absIntDiff(one, two) > 3 {
			violations++
			index = i
			if violations > 1 {
				return false
			}
		}
	}
	// If no violations, already sorted
	if violations == 0 {
		return true
	}
	// Check by removing arr[index] or arr[index+1]
	return ((index == 0 || arr[index-1] >= arr[index+1]) || (index+1 == len(arr)-1 || arr[index] >= arr[index+2])) && strings.Compare(strict, "strict") != 0
}

func atoiWrap(number string) int {
	num, err := strconv.Atoi(number)
	if err != nil {
		fmt.Println(fmt.Errorf("atoi failed: %w", err).Error())
		panic("atoi panic")
	}
	return num
}

// func isSortedDesc(arr []string) bool {
// 	fmt.Println(arr)
// 	for i := 0; i < len(arr)-1; i++ {
// 		one := atoiWrap(arr[i])
// 		two := atoiWrap(arr[i+1])
// 		if one < two || (absIntDiff(one, two) < 1 || absIntDiff(one, two) > 3) {
// 			return false
// 		}
// 	}
// 	return true
// }

// func isSortedAsc(arr []string) bool {
// 	fmt.Println(arr)
// 	for i := 0; i < len(arr)-1; i++ {
// 		one := atoiWrap(arr[i])
// 		two := atoiWrap(arr[i+1])
// 		if one > two || (absIntDiff(one, two) < 1 || absIntDiff(one, two) > 3) {
// 			return false
// 		}
// 	}
// 	return true
// }

func absIntDiff(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}
