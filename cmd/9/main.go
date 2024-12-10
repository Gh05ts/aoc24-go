package main

import (
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
	printRune(result)
	// fmt.Println(spaces)
	// fmt.Println(elements)
	// fmt.Println(sizeSpaceIdx)

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
			curElement := elements[i]
			bestFit := search(sizeSpaceIdx, curElement)
			if bestFit != -1 {
				emptySpace := spaces[bestFit]
				for i := 0; i < curElement.size; i++ {
					result[curElement.idx+i] = rune('𠔀')
					result[emptySpace.idx+i] = curElement.value
				}
				sizeSpaceIdx[emptySpace.size] = sizeSpaceIdx[emptySpace.size][1:]
				newLength := emptySpace.size - curElement.size
				spaces[bestFit].idx += curElement.size
				spaces[bestFit].size -= curElement.size
				if newLength != 0 {
					arr := sizeSpaceIdx[newLength]
					idx := -1
					for i := 0; i < len(arr); i++ {
						if arr[i].resultIdx >= emptySpace.idx+curElement.size {
							idx = i
							break
						}
					}
					if idx == -1 {
						idx = len(arr)
					}
					sizeSpaceIdx[newLength] = append(sizeSpaceIdx[newLength][:idx], append([]pair{{emptySpace.idx + curElement.size, bestFit}}, sizeSpaceIdx[newLength][idx:]...)...)
				}
			}
		}
	}

	fmt.Println(len(result))
	printRune(result)
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

func prepareInput(filesArray, emptyArray []int) (final []rune, elements, spaces []tuple, emptySizeSort map[int][]pair) {
	result := []rune{}
	elementsIdxSize, emptyIdxSize := []tuple{}, []tuple{}
	emptySizeSortIdx := make(map[int][]pair)

	// h := &PairHeap{}
	// heap.Init(h)
	// heap.Push(h, pair{0, 0})
	// // heap.Fix(h, 0)
	// // heap.Remove(h, 0)
	// for h.Len() > 0 {
	// 	fmt.Println(heap.Pop(h))
	// 	fmt.Println((*h)[0])
	// }

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
					emptySizeSortIdx[emptySpaceSize] = []pair{}
				}
				emptySizeSortIdx[emptySpaceSize] = append(emptySizeSortIdx[emptySpaceSize], pair{emptySpaceIdx, len(emptyIdxSize)})
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
			// fmt.Print(c)
			fmt.Print(string('#'))
		} else if c == rune('𠔐') {
			// fmt.Print(string('𠔐'))
			fmt.Print(string(' '))
		} else {
			// fmt.Print(string('𠔀'))
			fmt.Print(string(' '))
		}
	}
	fmt.Println()
}

// func searchAndUpdate(ts map[int][]int, size int) int {
// 	minIdx, sizeMatched := search(ts, size)
// 	update(ts, sizeMatched, size-sizeMatched)
// 	return minIdx
// }

func search(ts map[int][]pair, element tuple) int {
	minIdx := math.MaxInt
	for i := element.size; i < 10; i++ {
		value := ts[i]
		if len(value) > 0 && element.idx > value[0].resultIdx && value[0].resultIdx < minIdx {
			minIdx = value[0].emptyIdx
		}
	}

	if minIdx == math.MaxInt {
		fmt.Println("failed", element)
		return -1
	}
	return minIdx
}

func update(ts map[int][]int, size, updatedSize int) {
	// remove
	if updatedSize == 0 {

	}
	// remove and add

}

func insert(ts map[int][]int, size, updatedSize int) {

}

// for removing an index
// arr = append(arr[:indexToRemove], arr[indexToRemove+1])

// for finding the index
// index := sort.Search(len(arr), func(i int) bool { return arr[i] >= target })

// 7861696610269
// 7704727317120 too high
// 7658461057225
// 6327037990710
// 6327037990710
// 6347435485773
// 6320029754031 too low

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
	return x
}
