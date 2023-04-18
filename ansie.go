// Package ansi
//
// Adds support for ansi colours in the terminal.
//
// Provides fluent API similar to jansi library for Java (https://github.com/fusesource/jansi)
//
// ansi supports basic 7-colour, 256 colour and true colour modes. You can also use various attributes,
// such as underline, strike-out, etc.
//
// errMsg := Ansi.A("Error: ").Fg(Red).S("File not found: %s", fileName).Reset().A("Try a different name").String()
//
// See https://github.com/uaraven/ansi for more details
//
// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: (c) 2022 Oleksiy Voronin <ovoronin@gmail.com>

package ansie

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const esc = "\033["

type Attribute = int

type Colour = int

//goland:noinspection ALL
const (
	Reset Attribute = iota
	Bold
	Faint
	Italic
	Underline
	SlowBlink
	RapidBlink
	Reverse
	Conceal
	CrossOut

	NoBold Attribute = iota + 11
	Normal
	NoItalic
	NoUnderline
	NoBlink
	_
	NoReverse
	NoConceal
	NoCrossOut
)

type AnsiBuffer struct {
	enabled bool
	content strings.Builder
}

// NewAnsi creates a new AnsiBuffer. It doesn't assume anything about the device that the output will be
// directed to.
func NewAnsi() *AnsiBuffer {
	return &AnsiBuffer{enabled: true}
}

// Ansi is a default instance of AnsiBuffer
var Ansi = NewAnsiFor(os.Stdout)

// NewAnsiFor creates a new AnsiBuffer for a given device. It will not automatically print to this device,
// but it will disable ANSI colours if the device doesn't seem to support them, like when redirecting
// standard output into a file or piping it to another program
func NewAnsiFor(f *os.File) *AnsiBuffer {
	o, err := f.Stat()
	if err != nil {
		panic(err)
	}
	enabled := (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice
	return &AnsiBuffer{enabled: enabled}
}

func (ap *AnsiBuffer) CursorLeft(count int) *AnsiBuffer {
	ap.writeAnsiCommand('D', ';', count)
	return ap
}

func (ap *AnsiBuffer) CursorRight(count int) *AnsiBuffer {
	ap.writeAnsiCommand('C', ';', count)
	return ap
}

func (ap *AnsiBuffer) CursorUp(count int) *AnsiBuffer {
	ap.writeAnsiCommand('A', ';', count)
	return ap
}

func (ap *AnsiBuffer) CursorDown(count int) *AnsiBuffer {
	ap.writeAnsiCommand('B', ';', count)
	return ap
}

// Clear resets the internal buffer, after calling it the AnsiBuffer is in a clean state,
// as if just created
func (ap *AnsiBuffer) Clear() *AnsiBuffer {
	ap.content = strings.Builder{}
	return ap
}

// GetBuffer retrieves internal buffer as a string without clearing it
func (ap *AnsiBuffer) GetBuffer() string {
	return ap.content.String()
}

// IsEnabled returns true if colour output is enabled
func (ap *AnsiBuffer) IsEnabled() bool {
	return ap.enabled
}

// SetEnabled enables or disables the colour output. This does not affect the string already in AnsiBuffer's buffer
func (ap *AnsiBuffer) SetEnabled(value bool) {
	ap.enabled = value
}

// String converts the internal buffer to a string. The buffer is cleared after the call
func (ap *AnsiBuffer) String() string {
	s := ap.content.String()
	ap.Clear()
	return s
}

// Reset resets all the colours and attributes to defaults
func (ap *AnsiBuffer) Reset() *AnsiBuffer {
	ap.writeAnsiSeq(0)
	return ap
}

// Fg sets foreground colour. When using SysColour constants, like SysRed or SysYellow, the most basic and most compatible
// ANSI sequence will be used. Using any of other colour constants or integer values will use 256-colour ANSI sequence
func (ap *AnsiBuffer) Fg(colour Colour) *AnsiBuffer {
	if colour < 8 {
		ap.writeAnsiSeq(30 + colour)
	} else {
		ap.writeAnsiSeq(38, 5, colour)
	}
	return ap
}

// Bg sets background colour to one of standard 8 colours
func (ap *AnsiBuffer) Bg(colour Colour) *AnsiBuffer {
	if colour < 8 {
		ap.writeAnsiSeq(40 + colour)
	} else {
		ap.writeAnsiSeq(48, 5, colour)
	}
	return ap
}

// FgHi sets foreground colour to the high intensity version of one of standard 8 colours
//
// If used with one of 256 colour codes, it will just set the colour, without modifying the intensity
func (ap *AnsiBuffer) FgHi(colour Colour) *AnsiBuffer {
	if colour < 7 {
		ap.writeAnsiSeq(90 + colour)
		return ap
	} else {
		return ap.Fg(colour)
	}
}

// Attr sets font attribute
func (ap *AnsiBuffer) Attr(attr Attribute) *AnsiBuffer {
	ap.writeAnsiSeq(attr)
	return ap
}

// BgHi sets background colour to the high intensity version of one of standard 8 colours
//
// If used with one of 256 colour codes, it will just set the colour, without modifying the intensity
func (ap *AnsiBuffer) BgHi(colour Colour) *AnsiBuffer {
	if colour < 7 {
		ap.writeAnsiSeq(100 + colour)
		return ap
	} else {
		return ap.Bg(colour)
	}
}

// FgRgb sets foreground colour using "true colour" RGB colour
func (ap *AnsiBuffer) FgRgb(r, g, b uint) *AnsiBuffer {
	ap.writeAnsiSeq(38, 2, int(clip(r, 255)), int(clip(g, 255)), int(clip(b, 255)))
	return ap
}

// FgRgbI sets foreground colour using "true colour" RGB colour represented as a single integer
func (ap *AnsiBuffer) FgRgbI(i uint) *AnsiBuffer {
	r := (i >> 16) & 0xFF
	g := (i >> 8) & 0xFF
	b := i & 0xFF
	ap.writeAnsiSeq(38, 2, int(clip(r, 255)), int(clip(g, 255)), int(clip(b, 255)))
	return ap
}

// BgRgb sets foreground colour using "true colour" RGB colour
func (ap *AnsiBuffer) BgRgb(r, g, b uint) *AnsiBuffer {
	ap.writeAnsiSeq(48, 2, int(clip(r, 255)), int(clip(g, 255)), int(clip(b, 255)))
	return ap
}

// BgRgbI sets background colour using "true colour" RGB colour represented as a single integer
func (ap *AnsiBuffer) BgRgbI(i uint) *AnsiBuffer {
	r := (i >> 16) & 0xFF
	g := (i >> 8) & 0xFF
	b := i & 0xFF
	ap.writeAnsiSeq(48, 2, int(clip(r, 255)), int(clip(g, 255)), int(clip(b, 255)))
	return ap
}

// FgRgb6 sets foreground RGB colour as supported by 256-colour ANSI sequence
// In this mode each colour component is represented by a value from 0 to 5.
// Values beyond range of [0..5] are clipped
// R, G an B values are combined to get one of 216 colours supported by terminal
func (ap *AnsiBuffer) FgRgb6(r, g, b uint) *AnsiBuffer {
	colour := RgbTo216Colours(r, g, b)
	return ap.Fg(colour)
}

// BgRgb6 sets background RGB colour converted to 9-bit colour (3;3;3) as supported by 256-colour ANSI sequence
func (ap *AnsiBuffer) BgRgb6(r, g, b uint) *AnsiBuffer {
	colour := RgbTo216Colours(r, g, b)
	return ap.Bg(colour)
}

// FgGray sets foreground colour that is the shade of gray. intensity is a value in a range [0..23]. It is converted to
// one of standard 24 gray shades in the 256-colour palette
func (ap *AnsiBuffer) FgGray(intensity uint) *AnsiBuffer {
	return ap.Fg(Colour(232 + clip(intensity, 23)))
}

// BgGray sets background colour that is the shade of gray. intensity is a value in a range [0..24]. It is converted to
// one of standard 24 gray shades in the 256-colour palette
func (ap *AnsiBuffer) BgGray(intensity uint) *AnsiBuffer {
	return ap.Bg(Colour(232 + clip(intensity, 23)))
}

// FgGrayF sets foreground colour that is the shade of gray. intensity is a floating point value in a range [0..1].
// It is converted to one of standard 24 gray shades in the 256-colour palette
func (ap *AnsiBuffer) FgGrayF(intensity float64) *AnsiBuffer {
	gray := ap.shadeOfGrayColour(intensity)
	return ap.Fg(gray)
}

// BgGrayF sets background colour that is the shade of gray. intensity is a floating point value in a range [0..1].
// It is converted to one of standard 24 gray shades in the 256-colour palette
func (ap *AnsiBuffer) BgGrayF(intensity float64) *AnsiBuffer {
	gray := ap.shadeOfGrayColour(intensity)
	return ap.Bg(gray)
}

// A adds text to the AnsiBuffer's buffer. The text will be output with the current colours and attributes
func (ap *AnsiBuffer) A(text string) *AnsiBuffer {
	ap.content.WriteString(text)
	return ap
}

// S adds formatted (similar to fmt.Sprintf) text to the AnsiBuffer's buffer. The text will be output with the current colours and attributes
func (ap *AnsiBuffer) S(format string, params ...interface{}) *AnsiBuffer {
	ap.content.WriteString(fmt.Sprintf(format, params...))
	return ap
}

// CR adds carriage return character (ASCII 13) to the AnsiBuffer's buffer
func (ap *AnsiBuffer) CR() *AnsiBuffer {
	ap.content.WriteRune('\n')
	return ap
}

// Esc allows to add custom Esc sequence to the buffer
// The sequence that will be added is:
//
// ESC[codes sep codes sep codes sep command
//
// i.e. Esc('m', ':', 4, 3) will create sequence 'ESC[4:3m' which will create curly underline in iTerm2
func (ap *AnsiBuffer) Esc(command rune, sep rune, codes ...int) *AnsiBuffer {
	ap.writeAnsiCommand(command, sep, codes...)
	return ap
}

// EscM allows to add custom SGR sequences to the output
// EscM(38, 2, 0, 0, 255) will add sequence 'ESC[38;2;0;0;255m' to enable bright blue RGB colour
func (ap *AnsiBuffer) EscM(codes ...int) *AnsiBuffer {
	ap.writeAnsiSeq(codes...)
	return ap
}

// RgbTo216Colours converts a colour represented as R,G,B values of 0 to 5 to one of 216 colours
// in 256-colour palette
func RgbTo216Colours(r uint, g uint, b uint) Colour {
	colour := 16 + 36*clip(r, 5) + 6*clip(g, 5) + clip(b, 5)
	return Colour(colour)
}

func (ap *AnsiBuffer) writeAnsiCommand(command rune, sep rune, codes ...int) {
	if ap.enabled {
		ap.content.WriteString(esc)
		for _, code := range codes {
			ap.content.WriteString(strconv.Itoa(code))
			ap.content.WriteRune(sep)
		}
		ap.content.WriteRune(command)
	}
}

func (ap *AnsiBuffer) writeAnsiSeq(codes ...int) {
	ap.writeAnsiCommand('m', ';', codes...)
}

func clip(c uint, high uint) uint {
	if c > high {
		return high
	} else {
		return c
	}
}

func (ap *AnsiBuffer) shadeOfGrayColour(intensity float64) int {
	if intensity < 0 {
		intensity = 0
	}
	if intensity > 1 {
		intensity = 1
	}
	gray := 232 + int(23*intensity)
	return gray
}
