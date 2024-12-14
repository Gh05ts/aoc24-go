package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		args = append(args, "v1")
	}

	content, err := os.ReadFile("inputs/11.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	input := ioConvert(content)
	biggerInput := strings.Compare(args[1], "v2") == 0
	count := countStones(input, biggerInput)
	fmt.Println(count)
}

func ioConvert(content []byte) []int64 {
	input := []int64{}
	strContent := string(content)
	strInput := strings.Fields(strContent)
	for _, str := range strInput {
		input = append(input, atoiWrap(str))
	}
	return input
}

func atoiWrap(number string) int64 {
	num, err := strconv.Atoi(number)
	if err != nil {
		fmt.Println(fmt.Errorf("atoi failed: %w", err).Error())
		panic("atoi panic")
	}
	return int64(num)
}

func countStones(stones []int64, bigger bool) int64 {
	freqMap, copyFreqMap := make(map[int64]int64), make(map[int64]int64)
	for _, v := range stones {
		freqMap[v]++
		copyFreqMap[v]++
	}
	var limit int
	if bigger {
		limit = 75
	} else {
		limit = 25
	}
	start := time.Now()
	for i := 0; i < limit; i++ {
		for k, v := range freqMap {
			if v == 0 {
				continue
			}
			strVal := strconv.FormatInt(k, 10)
			lenStrVal := len(strVal)
			copyFreqMap[k] -= v
			switch {
			case k == 0:
				copyFreqMap[1] += v
			case lenStrVal%2 == 0:
				val1, val2 := atoiWrap(strVal[:lenStrVal/2]), atoiWrap(strVal[lenStrVal/2:])
				copyFreqMap[val1] += v
				copyFreqMap[val2] += v
			default:
				copyFreqMap[k*2024] += v
			}
		}
		for key, value := range copyFreqMap {
			freqMap[key] = value
		}
	}
	count := int64(0)
	for _, value := range copyFreqMap {
		count += value
	}
	end := time.Now()
	fmt.Println("time taken:", end.Sub(start))
	return count
}
