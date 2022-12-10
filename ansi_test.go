/*
 * SPDX-License-Identifier: MIT
 * SPDX-FileCopyrightText: (c) 2022 Oleksiy Voronin <ovoronin@gmail.com>
 */
package ansi

import (
	. "github.com/onsi/gomega"
	"log"
	"os"
	"testing"
)

func TestAnsiPrinter_Reset(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.Reset().String()

	g.Expect(s).To(Equal("\033[0;m"))
}

func TestAnsiPrinter_Clear(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.FgHi(SysBlue).A("some text")
	g.Expect(s.GetBuffer()).To(Equal("\033[94;msome text"))
	s.Clear()
	g.Expect(s.GetBuffer()).To(Equal(""))
}

func TestAnsiPrinter_FgHi(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.FgHi(SysBlue).A("some text").String()
	g.Expect(s).To(Equal("\033[94;msome text"))
	s = a.FgHi(Yellow).A("text").String()
	g.Expect(s).To(Equal("\033[38;5;11;mtext"))
}

func TestAnsiPrinter_BgHi(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.BgHi(SysBlue).A("some text").String()
	g.Expect(s).To(Equal("\033[104;msome text"))
	s = a.BgHi(Yellow).A("text").String()
	g.Expect(s).To(Equal("\033[48;5;11;mtext"))
}

func TestAnsiPrinter_ImplicitClear(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.FgHi(SysBlue).A("some text")
	g.Expect(s.GetBuffer()).To(Equal("\033[94;msome text"))
	g.Expect(s.String()).To(Equal("\033[94;msome text"))
	g.Expect(s.GetBuffer()).To(Equal(""))
}

func TestAnsiPrinter_Fg(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.Fg(SysRed).A("text").Reset().String()

	g.Expect(s).To(Equal("\033[31;mtext\033[0;m"))
}

func TestAnsiPrinter_Fg256(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.Fg(DarkGoldenrod).A("text").Reset().String()

	g.Expect(s).To(Equal("\033[38;5;136;mtext\033[0;m"))
}

func TestAnsiPrinter_Fg256rgb(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.FgRgb3(255, 255, 0).A("text").Reset().String()

	g.Expect(s).To(Equal("\033[38;5;226;mtext\033[0;m"))
}

func TestAnsiPrinter_Fg256rgbI(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.FgRgb3I(0xFFFF00).A("text").Reset().String()

	g.Expect(s).To(Equal("\033[38;5;226;mtext\033[0;m"))
}

func TestAnsiPrinter_Bg(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.Bg(SysRed).A("text").Reset().String()

	g.Expect(s).To(Equal("\033[41;mtext\033[0;m"))
}

func TestAnsiPrinter_Bg256(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.Bg(DarkGoldenrod).A("text").Reset().String()

	g.Expect(s).To(Equal("\033[48;5;136;mtext\033[0;m"))
}

func TestAnsiPrinter_Bg256rgb(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.BgRgb3(255, 255, 0).A("text").Reset().String()

	g.Expect(s).To(Equal("\033[48;5;226;mtext\033[0;m"))
}

func TestAnsiPrinter_Bg256rgbI(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.BgRgb3I(0xFFFF00).A("text").Reset().String()

	g.Expect(s).To(Equal("\033[48;5;226;mtext\033[0;m"))
}

func TestAnsiPrinter_FgRgb(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.FgRgb(255, 0, 10).A("text").Reset().String()

	g.Expect(s).To(Equal("\033[38;2;255;0;10;mtext\033[0;m"))
}

func TestAnsiPrinter_BgRgb(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.BgRgb(255, 1, 10).A("text").Reset().String()

	g.Expect(s).To(Equal("\033[48;2;255;1;10;mtext\033[0;m"))
}

func TestAnsiPrinter_FgRgbI(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.FgRgbI(0xFF000A).A("text").Reset().String()

	g.Expect(s).To(Equal("\033[38;2;255;0;10;mtext\033[0;m"))
}

func TestAnsiPrinter_BgRgbI(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.BgRgbI(0xFF000A).A("text").Reset().String()

	g.Expect(s).To(Equal("\033[48;2;255;0;10;mtext\033[0;m"))
}

func TestAnsiPrinter_Attr(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.Attr(Underline).A("text").String()

	g.Expect(s).To(Equal("\033[4;mtext"))
}

func TestAnsiPrinter_Custom(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.EscM(11, 2, 4).A("text").String()

	g.Expect(s).To(Equal("\033[11;2;4;mtext"))
}

func TestAnsiPrinter_Custom2(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.Esc('n', ':', 4, 3).A("text").String()

	g.Expect(s).To(Equal("\033[4:3:ntext"))
}

func TestAnsiPrinter_Disable(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	a.Fg(SysRed).A("red")
	a.SetEnabled(false)
	a.Fg(SysWhite).A("white")
	a.SetEnabled(true)
	s := a.Fg(SysBlue).A("blue").String()

	g.Expect(s).To(Equal("\033[31;mredwhite\033[34;mblue"))
}

func TestAnsiPrinter_IsEnabled(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	g.Expect(a.IsEnabled()).To(BeTrue())
	a.SetEnabled(false)
	g.Expect(a.IsEnabled()).To(BeFalse())
	a.SetEnabled(true)
	g.Expect(a.IsEnabled()).To(BeTrue())
}

func TestAnsiPrinter_FgGray(t *testing.T) {
	g := NewGomegaWithT(t)

	s := Ansi.FgGray(0.5).A("text").String()

	g.Expect(s).To(Equal("\033[38;5;243;mtext"))
}

func TestAnsiPrinter_BgGray(t *testing.T) {
	g := NewGomegaWithT(t)

	s := Ansi.BgGray(0.5).A("text").String()

	g.Expect(s).To(Equal("\033[48;5;243;mtext"))
}

func TestAnsiPrinter_Gray(t *testing.T) {
	g := NewGomegaWithT(t)

	g.Expect(Ansi.shadeOfGrayColor(-1)).To(Equal(232))
	g.Expect(Ansi.shadeOfGrayColor(22)).To(Equal(255))
}

func TestAnsiPrinter_CR(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.A("line1").CR().A("line2").String()
	g.Expect(s).To(Equal("line1\nline2"))
}

func TestAnsiPrinter_S(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.Fg(Fuchsia).S("this:%s", "is fuchsia").String()
	g.Expect(s).To(Equal("\033[38;5;13;mthis:is fuchsia"))
}

func TestAnsiFor(t *testing.T) {
	g := NewGomegaWithT(t)

	file, err := os.CreateTemp("", "ansi-test")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = os.Remove(file.Name()) }()

	a := NewAnsiFor(file)
	g.Expect(a.IsEnabled()).To(BeFalse())

	g.Expect(func() { NewAnsiFor(nil) }).To(Panic())
}
