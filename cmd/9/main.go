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

func sum(filesArray, emptyArray []int, fullFileSwap bool) int {
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
		for i := len(elements) - 1; i >= 0; i-- {
			fileSize := elements[i].size

			resultIdx, matchedSize := search(sizeSpaceIdx, fileSize)
			if resultIdx == math.MaxInt || resultIdx > elements[i].idx {
				continue
			}

			emptySpace := heap.Pop(sizeSpaceIdx[matchedSize])
			if space, ok := emptySpace.(pair); ok {
				for j := 0; j < fileSize; j++ {
					result[elements[i].idx+j] = rune('𠔀')
					result[spaces[space.emptyIdx].idx+j] = elements[i].value
				}

				newEmptySize := spaces[space.emptyIdx].size - fileSize
				spaces[space.emptyIdx].idx += fileSize
				spaces[space.emptyIdx].size -= fileSize
				if newEmptySize > 0 {
					heap.Push(sizeSpaceIdx[newEmptySize], pair{space.resultIdx + fileSize, space.emptyIdx})
				}
			}
		}
	}

	checksum := getCheckSum(result)
	return checksum
}

func search(sizeSpaceIdx map[int]*PairHeap, fileSize int) (int, int) {
	resultIdx := math.MaxInt
	matchedSize := math.MaxInt
	for size := fileSize; size < 10; size++ {
		if sizeSpaceIdx[size] != nil && sizeSpaceIdx[size].Len() > 0 && resultIdx > (*sizeSpaceIdx[size])[0].resultIdx {
			resultIdx = (*sizeSpaceIdx[size])[0].resultIdx
			matchedSize = size
		}
	}
	return resultIdx, matchedSize
}

func getCheckSum(result []rune) int {
	sum := 0
	for i, char := range result {
		if char == rune('𠔀') || char == rune('𠔐') {
			continue
		}
		sum += i * int(char)
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

// for debugging
// func printRune(result []rune) {
// 	for _, c := range result {
// 		if c != rune('𠔐') && c != rune('𠔀') {
// 			fmt.Print(c)
// 		} else if c == rune('𠔐') {
// 			fmt.Print(string('𠔐'))
// 		} else {
// 			fmt.Print(string('𠔀'))
// 		}
// 	}
// 	fmt.Println()
// }

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
