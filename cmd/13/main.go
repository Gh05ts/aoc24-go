package main

import (
	"fmt"
)

func main() {
	fmt.Println("hello")
}

// file, err := os.Open("input.txt")
// if err != nil {
// fmt.Println("Error opening file:", err)
// return
// }
// defer file.Close()
// var name string
//  var age int
// for {
// _, err := fmt.Fscanf(file, "%s %d\n", &name, &age)
// if err != nil {
// break
// }
// fmt.Printf("Name: %s, Age: %d\n", name, age)
// }

// FMT.SSCANF()
// "A:X+13,Y+14",
// "A:X+%D, Y+%D",
// &AX, &AY,
// )
