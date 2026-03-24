package main

import (
	"fmt"

	. "github.com/uaraven/ansie"
)

func main() {
	a := NewAnsi()
	fmt.Println(a.Attr(Underline).A("Standard colours example").Reset().String())
	for fg := range 8 {
		for bg := range 8 {
			fmt.Print(a.Fg(fg).Bg(bg).A("◘").Reset().String())
		}
		fmt.Println()
	}
}
