package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		args = append(args, "v1")
	}

	content, err := os.ReadFile("inputs/9.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	fullFileFlag := strings.Compare(args[1], "v2") == 0
	sum := sum(content, fullFileFlag)
	fmt.Println(sum)
}

type pair struct {
	idx, size int
}

func sum(content []byte, fullFlag bool) int {
	strContent := string(content)
	var filesArray, emptyArray []int

	for i, c := range strContent {
		if i%2 == 0 {
			filesArray = append(filesArray, atoiWrap(string(c)))
		} else {
			emptyArray = append(emptyArray, atoiWrap(string(c)))
		}
	}

	result := []rune{}
	startEmptyIdx, endEmptyIdx := 0, len(emptyArray)-1
	var elementsIdxSize, emptyIdxSize []pair
	for i := 0; i < len(filesArray); i++ {
		elementIdxSize := int(filesArray[i])
		elementsIdxSize = append(elementsIdxSize, pair{i, elementIdxSize})
		fill(&result, rune(i), elementIdxSize)
		if startEmptyIdx <= endEmptyIdx {
			emptySpaceIdx := len(result)
			emptySpaceSize := int(emptyArray[startEmptyIdx])
			if emptySpaceSize > 0 {
				emptyIdxSize = append(emptyIdxSize, pair{emptySpaceIdx, emptySpaceSize})
			}
			fill(&result, rune('𠔐'), emptySpaceSize)
			startEmptyIdx++
		}
	}

	// x < y is -ve, x > y is +ve and x == y is 0
	slices.SortFunc(emptyIdxSize, func(x, y pair) int {
		if x.idx < y.idx {
			return -1
		} else if x.idx > y.idx {
			return 1
		} else if x.size > y.size {
			return -1
		} else {
			return 1
		}
	})

	printRune(result)
	fmt.Println(emptyIdxSize)

	var elements []rune
	var elementsIndexArr []int
	for i, ch := range result {
		if ch != rune('𠔐') {
			elements = append(elements, ch)
			elementsIndexArr = append(elementsIndexArr, i)
		}
	}

	if !fullFlag {
		movedCount := 0
		elementIndex := len(elements) - 1
		for i := 0; i < len(result); i++ {
			if result[i] == rune('𠔐') && i <= elementsIndexArr[elementIndex] {
				result[i] = elements[elementIndex]
				result[elementsIndexArr[elementIndex]] = rune('𠔀')
				elementIndex--
				movedCount++
			}
		}
	} else {
		for {
			elementIndex := len(elements) - 1
			origElem := elementsIdxSize[elementIndex]

		}
	}

	printRune(result)

	ans := 0
	for i, char := range result {
		if char == rune('𠔀') || char == rune('𠔐') {
			continue
		}
		ans += i * int(char)
	}

	return ans
}

func fill(arr *[]rune, val rune, times int) {
	for times > 0 {
		*arr = append(*arr, val)
		times--
	}
}

func atoiWrap(number string) int {
	num, err := strconv.Atoi(number)
	if err != nil {
		fmt.Println(fmt.Errorf("atoi failed: %w", err).Error())
		panic("atoi panic")
	}
	return num
}

func printRune(result []rune) {
	for _, c := range result {
		if c != rune('𠔐') && c != rune('𠔀') {
			fmt.Print(c)
		} else if c == rune('𠔐') {
			fmt.Print(string('𠔐'))
		} else {
			fmt.Print(string('𠔀'))
		}
	}
	fmt.Println()
}

func Insert(ts []pair, t pair) []pair {
	i, _ := customSearch(ts, t)
	ts = append(ts, *new(pair))
	copy(ts[i+1:], ts[i:])
	ts[i] = t
	return ts
}

func customSearch(ts []pair, t pair) (int, bool) {
	for i, j := range ts {
		if j.size >= t.size {
			return i, j.idx > t.idx
		}
	}
	return -1, false
}
