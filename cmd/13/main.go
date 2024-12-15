package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
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

	reader := fileReader("inputs/13.txt")
	defer reader.Close()

	buttonPrize := ioConvert(reader)
	fmt.Println(buttonPrize)
	flag := strings.Compare(args[1], "v2") == 0
	fmt.Println(flag)
}

func fileReader(path string) *bufioCloser {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(fmt.Errorf("open failed: %w", err).Error())
		panic("file panic")
	}

	return &bufioCloser{bufio.NewReader(file), file}
}

type buttonPrize struct {
	a, b, prize pair
}

type pair struct {
	x, y int
}

func ioConvert(reader *bufioCloser) []buttonPrize {
	arr := []buttonPrize{}
	i := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)

		var x, y int
		switch {
		case i%4 == 0:
			fmt.Sscanf(line, "Button A: X+%d, Y+%d", &x, &y)
			arr = append(arr, buttonPrize{a: pair{x, y}})
		case i%4 == 1:
			fmt.Sscanf(line, "Button B: X+%d, Y+%d", &x, &y)
			arr[len(arr)-1].b = pair{x, y}
		case i%4 == 2:
			fmt.Sscanf(line, "Prize: X=%d, Y=%d", &x, &y)
			arr[len(arr)-1].prize = pair{x, y}
		}
		i++
	}
	return arr
}
