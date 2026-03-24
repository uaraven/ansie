package main

import (
	"fmt"

	. "github.com/uaraven/ansie"
)

func main() {
	a := NewAnsi()
	fmt.Println(a.Attr(Underline).A("256 colour example").Reset().String())
	for r := range 6 {
		for g := range 6 {
			for b := range 6 {
				fmt.Print(a.FgRgb6(uint(r), uint(g), uint(b)).A("█").String())
			}
		}
		fmt.Println()
	}
}
