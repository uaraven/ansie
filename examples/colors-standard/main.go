package main

import (
	"fmt"
	"strconv"

	. "github.com/uaraven/ansie"
)

func main() {
	a := NewAnsi()
	fmt.Println(a.Attr(Underline).A("Standard colours example").Reset().String())
	fmt.Print("    ")
	for fg := range 8 {
		fmt.Print(a.Fg(fg).A(strconv.FormatInt(int64(fg), 10)).A(" ").Reset().
			FgHi(fg).A(strconv.FormatInt(int64(fg), 10)).A(" ").Reset().
			String())
	}
	fmt.Println()
	for bg := range 8 {
		fmt.Print(a.Bg(bg).A(strconv.FormatInt(int64(bg), 10)).Reset().A("   "))
		for fg := range 8 {
			fmt.Print(a.Bg(bg).Fg(fg).A("•").Reset().A(" ").FgHi(fg).Bg(bg).A("•").Reset().A(" ").String())
		}
		fmt.Println()
		fmt.Print(a.BgHi(bg).A(strconv.FormatInt(int64(bg), 10)).Reset().A("   "))
		for fg := range 8 {
			fmt.Print(a.BgHi(bg).Fg(fg).A("•").Reset().A(" ").FgHi(fg).BgHi(bg).A("•").Reset().A(" ").String())
		}
		fmt.Println()
	}
}
