package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"strconv"
	"sort"
	"math"
)

func main() {
	file, err := os.Open("inputs/1.txt")
	if err != nil {
		fmt.Println("Bad file")
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	var list1, list2 []int

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		
		numbers := strings.Fields(line)
		num1, err := strconv.Atoi(numbers[0])
		if err != nil {
			fmt.Println("can't convert")
			continue
		}
		list1 = append(list1, num1)

		num2, err := strconv.Atoi(numbers[1])
		if err != nil {
			fmt.Println("can't convert")
			continue
		}
		list2 = append(list2, num2)
	}

	sort.Ints(list1)
	sort.Ints(list2)

	result := 0

	for i := range list1 {
		result += int(math.Abs(float64(list1[i]) - float64(list2[i])))
	}

	fmt.Println(result)
}
