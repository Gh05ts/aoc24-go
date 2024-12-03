package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	lenDont = 7
)

// type Set map[int]bool
// type MapOfSets map[string]Set

// func addToSet(mos MapOfSets, key string, value int) {
// 	if _, ok := mos[key]; !ok {
// 		mos[key] = make(Set)
// 	}
// 	mos[key][value] = true
// }

func main() {
	args := os.Args

	if len(args) < 2 {
		return
	}

	content, err := os.ReadFile("inputs/3.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	strContent := string(content)

	// mos := make(MapOfSets)
	var dontIndices, doIndices []int
	re := regexp.MustCompile(`(do\(\)|don't\(\))`)
	for _, match := range re.FindAllStringIndex(strContent, -1) {
		start := match[0]
		end := match[1]

		if end-start == lenDont {
			dontIndices = append(dontIndices, start)
		} else {
			doIndices = append(doIndices, start)
		}
	}

	result := 0
	dontIdx := 0
	doIdx := 0
	allowFlag := true
	re = regexp.MustCompile(`mul\([0-9]+,[0-9]+\)`)
	for _, match := range re.FindAllStringIndex(strContent, -1) {
		reg := regexp.MustCompile(`\d+`)
		mul := 0
		start := match[0]
		end := match[1]
		allowFlag = check(&doIdx, &dontIdx, start, doIndices, dontIndices) || strings.Compare(args[1], "stage") != 0
		for _, num := range reg.FindAllString(strContent[start:end], -1) {
			if allowFlag {
				if mul == 0 {
					mul = atoiWrap(num)
				} else {
					mul *= atoiWrap(num)
				}
			} else {
				break
			}
		}
		result += mul
	}

	fmt.Println(result)
}

func check(doIdx, dontIdx *int, curIdx int, doIndices, dontIndices []int) bool {
	for *doIdx < len(doIndices) && doIndices[*doIdx] < curIdx {
		*doIdx += 1
	}
	for *dontIdx < len(dontIndices) && dontIndices[*dontIdx] < curIdx {
		*dontIdx += 1
	}

	if *dontIdx <= 0 {
		return true
	}
	if *dontIdx >= 0 && *doIdx <= 0 {
		return false
	}
	return doIndices[*doIdx-1] > dontIndices[*dontIdx-1]
}

func atoiWrap(number string) int {
	num, err := strconv.Atoi(number)
	if err != nil {
		fmt.Println(fmt.Errorf("atoi failed: %w", err).Error())
		panic("atoi panic")
	}
	return num
}
