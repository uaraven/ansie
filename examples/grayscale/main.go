package main

import (
	"fmt"

	. "github.com/uaraven/ansie"
)

func main() {
	a := NewAnsi()
	fmt.Println(a.Attr(Underline).A("Grayscale example").Reset().String())
	for g := range 24 {
		fmt.Print(a.FgGray(uint(g)).A("█").String())
	}
	fmt.Println()
	fmt.Println(a.Attr(Italic).Attr(Underline).A("0..1 grayscale example").Reset().String())
	for g := 0.0; g <= 1.0; g += 0.042 {
		fmt.Print(a.FgGrayF(g).A("█").String())
	}
	fmt.Println()
}
