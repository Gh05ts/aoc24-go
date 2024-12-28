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

	flag := strings.Compare(args[1], "v2") == 0
	buttonPrize := ioConvert(reader, flag)

	tokens := fewestTokens(buttonPrize)
	fmt.Println(tokens)
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
	x, y int64
}

func ioConvert(reader *bufioCloser, flag bool) []buttonPrize {
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
			arr = append(arr, buttonPrize{a: pair{int64(x), int64(y)}})
		case i%4 == 1:
			fmt.Sscanf(line, "Button B: X+%d, Y+%d", &x, &y)
			arr[len(arr)-1].b = pair{int64(x), int64(y)}
		case i%4 == 2:
			fmt.Sscanf(line, "Prize: X=%d, Y=%d", &x, &y)
			// 10000000000000
			if flag {
				arr[len(arr)-1].prize = pair{int64(x) + 10000000000000, int64(y) + 10000000000000}
			} else {
				arr[len(arr)-1].prize = pair{int64(x), int64(y)}
			}
		}
		i++
	}
	return arr
}

// i*ax+j*bx=px
// i = px - j*bx / ax
// i*ay+j*by=py
// (px-j*bx)/ax *ay+j*by=py
// (px-j*bx)*ay+j*by*ax=py*ax
// px.ay - j*ay*bx + j*by*ax = py*ax
// j*(by*ax - ay*bx) = py*ax-px.ay
// j = (py*ax - px*ay) / (by*ax - bx*ay)
// i = (py*bx - px*by) / (bx*ay - by*ax)

func fewestTokens(groups []buttonPrize) int64 {
	tokens := int64(0)
	for _, group := range groups {
		// for i := int64(0); i*group.a.x <= group.prize.x; i++ {
		// 	j := (group.prize.x - i*group.a.x) / group.b.x
		// 	if i*group.a.x+j*group.b.x == group.prize.x && i*group.a.y+j*group.b.y == group.prize.y {
		// 		tokens += i*3 + j
		// 	}
		// }
		i := (group.prize.y*group.b.x - group.prize.x*group.b.y) / (group.b.x*group.a.y - group.b.y*group.a.x)
		j := (group.prize.y*group.a.x - group.prize.x*group.a.y) / (group.b.y*group.a.x - group.b.x*group.a.y)
		if i*group.a.x+j*group.b.x == group.prize.x && i*group.a.y+j*group.b.y == group.prize.y {
			tokens += i*3 + j
		}
	}
	return tokens
}
