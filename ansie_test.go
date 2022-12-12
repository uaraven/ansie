/*
 * SPDX-License-Identifier: MIT
 * SPDX-FileCopyrightText: (c) 2022 Oleksiy Voronin <ovoronin@gmail.com>
 */
package ansie

import (
	. "github.com/onsi/gomega"
	"log"
	"os"
	"testing"
)

func TestAnsiBuffer_Reset(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.Reset().String()

	g.Expect(s).To(Equal("\033[0;m"))
}

func TestAnsiBuffer_Clear(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.FgHi(Blue).A("some text")
	g.Expect(s.GetBuffer()).To(Equal("\033[94;msome text"))
	s.Clear()
	g.Expect(s.GetBuffer()).To(Equal(""))
}

func TestAnsiBuffer_FgHi(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.FgHi(Blue).A("some text").String()
	g.Expect(s).To(Equal("\033[94;msome text"))
	s = a.FgHi(BrightYellow).A("text").String()
	g.Expect(s).To(Equal("\033[38;5;11;mtext"))
}

func TestAnsiBuffer_BgHi(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.BgHi(Blue).A("some text").String()
	g.Expect(s).To(Equal("\033[104;msome text"))
	s = a.BgHi(BrightYellow).A("text").String()
	g.Expect(s).To(Equal("\033[48;5;11;mtext"))
}

func TestAnsiBuffer_ImplicitClear(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.FgHi(Blue).A("some text")
	g.Expect(s.GetBuffer()).To(Equal("\033[94;msome text"))
	g.Expect(s.String()).To(Equal("\033[94;msome text"))
	g.Expect(s.GetBuffer()).To(Equal(""))
}

func TestAnsiBuffer_Fg(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.Fg(Red).A("text").Reset().String()

	g.Expect(s).To(Equal("\033[31;mtext\033[0;m"))
}

func TestAnsiBuffer_Fg256(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.Fg(DarkGoldenrod).A("text").Reset().String()

	g.Expect(s).To(Equal("\033[38;5;136;mtext\033[0;m"))
}

func TestAnsiBuffer_Fg256rgb(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.FgRgb6(5, 5, 0).A("text").Reset().String()

	g.Expect(s).To(Equal("\033[38;5;226;mtext\033[0;m"))
}

func TestAnsiBuffer_Bg(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.Bg(Red).A("text").Reset().String()

	g.Expect(s).To(Equal("\033[41;mtext\033[0;m"))
}

func TestAnsiBuffer_Bg256(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.Bg(DarkGoldenrod).A("text").Reset().String()

	g.Expect(s).To(Equal("\033[48;5;136;mtext\033[0;m"))
}

func TestAnsiBuffer_Bg256rgb(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.BgRgb6(5, 5, 0).A("text").Reset().String()

	g.Expect(s).To(Equal("\033[48;5;226;mtext\033[0;m"))
}

func TestAnsiBuffer_FgRgb(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.FgRgb(255, 0, 10).A("text").Reset().String()

	g.Expect(s).To(Equal("\033[38;2;255;0;10;mtext\033[0;m"))
}

func TestAnsiBuffer_BgRgb(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.BgRgb(255, 1, 10).A("text").Reset().String()

	g.Expect(s).To(Equal("\033[48;2;255;1;10;mtext\033[0;m"))
}

func TestAnsiBuffer_FgRgbI(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.FgRgbI(0xFF000A).A("text").Reset().String()

	g.Expect(s).To(Equal("\033[38;2;255;0;10;mtext\033[0;m"))
}

func TestAnsiBuffer_BgRgbI(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.BgRgbI(0xFF000A).A("text").Reset().String()

	g.Expect(s).To(Equal("\033[48;2;255;0;10;mtext\033[0;m"))
}

func TestAnsiBuffer_Attr(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.Attr(Underline).A("text").String()

	g.Expect(s).To(Equal("\033[4;mtext"))
}

func TestAnsiBuffer_Custom(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.EscM(11, 2, 4).A("text").String()

	g.Expect(s).To(Equal("\033[11;2;4;mtext"))
}

func TestAnsiBuffer_Custom2(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.Esc('n', ':', 4, 3).A("text").String()

	g.Expect(s).To(Equal("\033[4:3:ntext"))
}

func TestAnsiBuffer_Disable(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	a.Fg(Red).A("red")
	a.SetEnabled(false)
	a.Fg(White).A("white")
	a.SetEnabled(true)
	s := a.Fg(Blue).A("blue").String()

	g.Expect(s).To(Equal("\033[31;mredwhite\033[34;mblue"))
}

func TestAnsiBuffer_IsEnabled(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	g.Expect(a.IsEnabled()).To(BeTrue())
	a.SetEnabled(false)
	g.Expect(a.IsEnabled()).To(BeFalse())
	a.SetEnabled(true)
	g.Expect(a.IsEnabled()).To(BeTrue())
}

func TestAnsiBuffer_FgGray(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()

	s := a.FgGrayF(0.5).A("text").String()
	g.Expect(s).To(Equal("\033[38;5;243;mtext"))

	s = a.FgGray(12).A("text").String()
	g.Expect(s).To(Equal("\033[38;5;244;mtext"))
}

func TestAnsiBuffer_BgGray(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()

	s := a.BgGrayF(0.5).A("text").String()
	g.Expect(s).To(Equal("\033[48;5;243;mtext"))

	s = a.BgGray(12).A("text").String()
	g.Expect(s).To(Equal("\033[48;5;244;mtext"))
}

func TestAnsiBuffer_Gray(t *testing.T) {
	g := NewGomegaWithT(t)

	g.Expect(Ansi.shadeOfGrayColour(-1)).To(Equal(232))
	g.Expect(Ansi.shadeOfGrayColour(22)).To(Equal(255))
}

func TestAnsiBuffer_CR(t *testing.T) {
	g := NewGomegaWithT(t)

	a := NewAnsi()
	s := a.A("line1").CR().A("line2").String()
	g.Expect(s).To(Equal("line1\nline2"))
}

func TestAnsiBuffer_S(t *testing.T) {
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

func TestClip(t *testing.T) {
	g := NewGomegaWithT(t)

	g.Expect(clip(150, 4)).To(Equal(uint(4)))
	g.Expect(clip(3, 4)).To(Equal(uint(3)))
}
