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
		args = append(args, "v1")
	}

	reader := fileReader("inputs/5.txt")
	defer reader.Close()

	rules, ops := ioConvert(reader)
	count := MiddleSum(rules, ops, args[1])
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

func ioConvert(reader *bufioCloser) (rulesArray, opsArray [][]int) {
	var rules [][]int
	var ops [][]int
	opsFlag := false
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if opsFlag || line == "" {
			if !opsFlag {
				opsFlag = true
				continue
			}
			ops = append(ops, stringToIntArray(line, ","))
		} else {
			rules = append(rules, stringToIntArray(line, "|"))
		}
	}
	return rules, ops
}

func stringToIntArray(str string, sep string) []int {
	stringArray := strings.Split(str, sep)
	intArray := []int{}
	for _, str := range stringArray {
		num := atoiWrap(str)
		intArray = append(intArray, num)
	}
	return intArray
}

func atoiWrap(number string) int {
	num, err := strconv.Atoi(number)
	if err != nil {
		fmt.Println(fmt.Errorf("atoi failed: %w", err).Error())
		panic("atoi panic")
	}
	return num
}

func MiddleSum(ops, rules [][]int, stage string) int {
	count := 0
	_, backwardMap := opsToOpsMap(ops)
	failed := []int{}
	for i, val := range rules {
		seen := []int{}
		violated := false
		for _, cur := range val {
			// loop seen and check if any matches?
			// can this be improved?
			for _, val := range seen {
				if !backwardMap[cur][val] {
					violated = true
					break
				}
			}
			seen = append(seen, cur)
		}
		if !violated {
			middle := seen[len(seen)/2]
			count += middle
		} else {
			failed = append(failed, i)
		}
	}
	if strings.Compare(stage, "v2") == 0 {
		count = corrections(rules, failed, backwardMap)
	}
	return count
}

func opsToOpsMap(ops [][]int) (forwardMap, backwardMap map[int]map[int]bool) {
	fwMap := make(map[int]map[int]bool)
	bwMap := make(map[int]map[int]bool)
	for _, op := range ops {
		key := op[0]
		val := op[1]
		if _, found := fwMap[key]; !found {
			fwMap[key] = make(map[int]bool)
		}
		fwMap[key][val] = true
		if _, found := bwMap[val]; !found {
			bwMap[val] = make(map[int]bool)
		}
		bwMap[val][key] = true
	}
	return fwMap, bwMap
}

func corrections(rules [][]int, failedIndices []int, backwardMap map[int]map[int]bool) int {
	count := 0
	for _, index := range failedIndices {
		rule := rules[index]
		correctMap := make(map[int]int)
		for _, val := range rule {
			lenBack := 0
			for _, check := range rule {
				if _, found := backwardMap[val][check]; found {
					lenBack++
				}
			}
			correctMap[lenBack] = val
		}
		count += correctMap[len(rule)/2]
	}
	return count
}
