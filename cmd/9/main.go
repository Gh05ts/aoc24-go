package main

import (
	"container/heap"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
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

	filesArray, emptyArray := ioConvert(content)
	fullFileFlag := strings.Compare(args[1], "v2") == 0

	sum := sum(filesArray, emptyArray, fullFileFlag)
	fmt.Println(sum)
}

type tuple struct {
	idx, size int
	value     rune
}

type pair struct {
	resultIdx, emptyIdx int
}

func ioConvert(content []byte) ([]int, []int) {
	strContent := string(content)
	var filesArray, emptyArray []int

	for i, c := range strContent {
		if i%2 == 0 {
			filesArray = append(filesArray, atoiWrap(string(c)))
		} else {
			emptyArray = append(emptyArray, atoiWrap(string(c)))
		}
	}
	return filesArray, emptyArray
}

func atoiWrap(number string) int {
	num, err := strconv.Atoi(number)
	if err != nil {
		fmt.Println(fmt.Errorf("atoi failed: %w", err).Error())
		panic("atoi panic")
	}
	return num
}

func sum(filesArray, emptyArray []int, fullFileSwap bool) int64 {
	result, elements, spaces, sizeSpaceIdx := prepareInput(filesArray, emptyArray)

	if !fullFileSwap {
		lastElementIdx := len(elements) - 1
		lastElement := elements[lastElementIdx]
		breakAll := false
		for _, spaceTuple := range spaces {
			for i := 0; i < spaceTuple.size; i++ {
				if lastElement.size == 0 {
					lastElementIdx--
					lastElement = elements[lastElementIdx]
				}
				if spaceTuple.idx+i > lastElement.idx+lastElement.size-1 {
					breakAll = true
					break
				}
				result[spaceTuple.idx+i] = lastElement.value
				result[lastElement.idx+lastElement.size-1] = rune('𠔀')
				lastElement.size -= 1
			}
			if breakAll {
				break
			}
		}
	} else {
		i := len(result) - 1
		for i >= 0 {
			if result[i] == rune('𠔐') || result[i] == rune('𠔀') {
				i -= 1
				continue
			}

			fileWidth := 0
			fileID := result[i]
			for i >= 0 && fileID == result[i] {
				fileWidth++
				i -= 1
			}

			smallestIdx := math.MaxInt
			resultIdx := math.MaxInt
			bestWidth := -1
			for size := fileWidth; size < 10; size++ {
				if sizeSpaceIdx[size] != nil && sizeSpaceIdx[size].Len() > 0 && resultIdx > (*sizeSpaceIdx[size])[0].resultIdx {
					smallestIdx = (*sizeSpaceIdx[size])[0].emptyIdx
					resultIdx = (*sizeSpaceIdx[size])[0].resultIdx
					bestWidth = size
				}
			}
			if smallestIdx == math.MaxInt || resultIdx > i {
				continue
			}

			emptySpace := heap.Pop(sizeSpaceIdx[bestWidth])
			if typed, ok := emptySpace.(pair); ok {
				newLength := spaces[typed.emptyIdx].size - fileWidth
				for j := 0; j < fileWidth; j++ {
					result[i+j+1] = rune('𠔀')
					result[spaces[typed.emptyIdx].idx+j] = fileID
				}
				spaces[typed.emptyIdx].idx += fileWidth
				spaces[typed.emptyIdx].size -= fileWidth
				if newLength > 0 {
					heap.Push(sizeSpaceIdx[newLength], pair{typed.resultIdx + fileWidth, typed.emptyIdx})
				}
			}
		}
	}

	checksum := getCheckSum(result)
	return checksum
}

func getCheckSum(result []rune) int64 {
	var sum int64 = 0
	for i, char := range result {
		if char == rune('𠔀') || char == rune('𠔐') {
			continue
		}
		sum += int64(i) * int64(char)
	}
	return sum
}

func prepareInput(filesArray, emptyArray []int) (final []rune, elements, spaces []tuple, emptySizeSort map[int]*PairHeap) {
	result := []rune{}
	elementsIdxSize, emptyIdxSize := []tuple{}, []tuple{}
	emptySizeSortIdx := make(map[int]*PairHeap)

	startEmptyIdx, endEmptyIdx := 0, len(emptyArray)-1
	for i := 0; i < len(filesArray); i++ {
		elementIdxSize := int(filesArray[i])
		elementsIdxSize = append(elementsIdxSize, tuple{len(result), elementIdxSize, rune(i)})
		fill(&result, rune(i), elementIdxSize)

		if startEmptyIdx <= endEmptyIdx {
			emptySpaceIdx := len(result)
			emptySpaceSize := int(emptyArray[startEmptyIdx])
			if emptySpaceSize > 0 {
				if _, found := emptySizeSortIdx[emptySpaceSize]; !found {
					emptySizeSortIdx[emptySpaceSize] = &PairHeap{}
					heap.Init(emptySizeSortIdx[emptySpaceSize])
				}
				heap.Push(emptySizeSortIdx[emptySpaceSize], pair{emptySpaceIdx, len(emptyIdxSize)})
				emptyIdxSize = append(emptyIdxSize, tuple{emptySpaceIdx, emptySpaceSize, rune('𠔐')})
			}
			fill(&result, rune('𠔐'), emptySpaceSize)
			startEmptyIdx++
		}
	}
	return result, elementsIdxSize, emptyIdxSize, emptySizeSortIdx
}

func fill(arr *[]rune, val rune, times int) {
	for times > 0 {
		*arr = append(*arr, val)
		times--
	}
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

type PairHeap []pair

func (h PairHeap) Len() int           { return len(h) }
func (h PairHeap) Less(i, j int) bool { return h[i].emptyIdx < h[j].emptyIdx }
func (h PairHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *PairHeap) Push(x any) {
	*h = append(*h, x.(pair))
}

func (h *PairHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return pair(x)
}
