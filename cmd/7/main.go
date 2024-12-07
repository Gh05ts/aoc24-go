package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/exp/maps"
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

	reader := fileReader("inputs/7.txt")
	defer reader.Close()

	targetInputs := ioConvert(reader)
	count := validCount(targetInputs, args[1])
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

type pair struct {
	x, y int
}

func ioConvert(reader *bufioCloser) map[pair][]int {
	targetInputs := make(map[pair][]int)
	index := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		target, inputs := stringToIntSlice(line)
		targetInputs[pair{index, target}] = inputs
		index++
	}
	return targetInputs
}

func stringToIntSlice(str string) (int, []int) {
	stringSlice := strings.Split(str, ":")
	target := atoiWrap(stringSlice[0])
	intSlice := []int{}
	inputsSlice := strings.Fields(stringSlice[1])
	for _, str := range inputsSlice {
		num := atoiWrap(str)
		intSlice = append(intSlice, num)
	}
	return target, intSlice
}

func atoiWrap(number string) int {
	num, err := strconv.Atoi(number)
	if err != nil {
		fmt.Println(fmt.Errorf("atoi failed: %w", err).Error())
		panic("atoi panic")
	}
	return num
}

func validCount(targetInputs map[pair][]int, stage string) int {
	count := 0
	allowed := strings.Compare(stage, "v2") == 0
	for target, inputs := range targetInputs {
		memo := make(map[string]bool)
		if canFormTarget(inputs, target.y, inputs[0], 1, memo, allowed) {
			count += target.y
		}
		maps.Clear(memo)
	}
	return count
}

func canFormTarget(nums []int, target int, current int, index int, memo map[string]bool, concatAllowed bool) bool {
	key := fmt.Sprintf("%d-%d", index, current)

	if result, exists := memo[key]; exists {
		return result
	}

	if index == len(nums)-1 {
		if current+nums[index] == target {
			memo[key] = true
			return true
		}
		if current*nums[index] == target {
			memo[key] = true
			return true
		}
		if concatAllowed {
			concatenated := concatenateNumbers(current, nums[index])
			if concatenated == target {
				memo[key] = true
				return true
			}
		}

		memo[key] = false
		return false
	}

	if canFormTarget(nums, target, current+nums[index], index+1, memo, concatAllowed) {
		memo[key] = true
		return true
	}
	if canFormTarget(nums, target, current*nums[index], index+1, memo, concatAllowed) {
		memo[key] = true
		return true
	}
	if concatAllowed {
		concatenated := concatenateNumbers(current, nums[index])
		if canFormTarget(nums, target, concatenated, index+1, memo, concatAllowed) {
			memo[key] = true
			return true
		}
	}

	memo[key] = false
	return false
}

func concatenateNumbers(a, b int) int {
	strA := strconv.Itoa(a)
	strB := strconv.Itoa(b)
	concatenatedStr := strA + strB
	concatenated, _ := strconv.Atoi(concatenatedStr)
	return concatenated
}
