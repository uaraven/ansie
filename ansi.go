// Package ansi
//
// Adds support for ansi colors in the terminal.
//
// Provides fluent API similar to jansi library for Java (https://github.com/fusesource/jansi)
//
// errMsg := Ansi.A("Error: ").Fg(Red).S("File not found: %s", fileName).Reset().A("Try a different name").String()
//
// See https://github.com/uaraven/ansi for more details
//
// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: (c) 2022 Oleksiy Voronin <ovoronin@gmail.com>

package ansi

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

//goland:noinspection ALL
const (
	Black Colour = 0
	// SysBlack is a color from standard 7-color terminal palette
	SysBlack Colour = 0
	// SysRed is a color from standard 7-color terminal palette
	SysRed Colour = -1
	// SysGreen is a color from standard 7-color terminal palette
	SysGreen Colour = -2
	// SysYellow is a color from standard 7-color terminal palette
	SysYellow Colour = -3
	// SysBlue is a color from standard 7-color terminal palette
	SysBlue Colour = -4
	// SysMagenta is a color from standard 7-color terminal palette
	SysMagenta Colour = -5
	// SysCyan is a color from standard 7-color terminal palette
	SysCyan Colour = -6
	// SysWhite is a Red color from standard 7-color terminal palette
	SysWhite Colour = -7

	Maroon            Colour = 1
	Green             Colour = 2
	Olive             Colour = 3
	Navy              Colour = 4
	Purple            Colour = 5
	Teal              Colour = 6
	Silver            Colour = 7
	Grey              Colour = 8
	Red               Colour = 9
	Lime              Colour = 10
	Yellow            Colour = 11
	Blue              Colour = 12
	Fuchsia           Colour = 13
	Aqua              Colour = 14
	White             Colour = 15
	Grey0             Colour = 16
	NavyBlue          Colour = 17
	DarkBlue          Colour = 18
	Blue3             Colour = 20
	Blue1             Colour = 21
	DarkGreen         Colour = 22
	DeepSkyBlue4      Colour = 25
	DodgerBlue3       Colour = 26
	DodgerBlue2       Colour = 27
	Green4            Colour = 28
	SpringGreen4      Colour = 29
	Turquoise4        Colour = 30
	DeepSkyBlue3      Colour = 32
	DodgerBlue1       Colour = 33
	Green3            Colour = 34
	SpringGreen3      Colour = 35
	DarkCyan          Colour = 36
	LightSeaGreen     Colour = 37
	DeepSkyBlue2      Colour = 38
	DeepSkyBlue1      Colour = 39
	SpringGreen2      Colour = 42
	Cyan3             Colour = 43
	DarkTurquoise     Colour = 44
	Turquoise2        Colour = 45
	Green1            Colour = 46
	SpringGreen1      Colour = 48
	MediumSpringGreen Colour = 49
	Cyan2             Colour = 50
	Cyan1             Colour = 51
	Purple4           Colour = 55
	Purple3           Colour = 56
	BlueViolet        Colour = 57
	Grey37            Colour = 59
	MediumPurple4     Colour = 60
	SlateBlue3        Colour = 62
	RoyalBlue1        Colour = 63
	Chartreuse4       Colour = 64
	PaleTurquoise4    Colour = 66
	SteelBlue         Colour = 67
	SteelBlue3        Colour = 68
	CornflowerBlue    Colour = 69
	DarkSeaGreen4     Colour = 71
	CadetBlue         Colour = 72
	SkyBlue3          Colour = 74
	Chartreuse3       Colour = 76
	PaleGreen3        Colour = 77
	SeaGreen3         Colour = 78
	Aquamarine3       Colour = 79
	MediumTurquoise   Colour = 80
	SteelBlue1        Colour = 81
	SeaGreen2         Colour = 83
	SeaGreen1         Colour = 85
	DarkSlateGray2    Colour = 87
	DarkRed           Colour = 88
	DarkMagenta       Colour = 91
	Orange4           Colour = 94
	LightPink4        Colour = 95
	Plum4             Colour = 96
	MediumPurple3     Colour = 98
	SlateBlue1        Colour = 99
	Wheat4            Colour = 101
	Grey53            Colour = 102
	LightSlateGrey    Colour = 103
	MediumPurple      Colour = 104
	LightSlateBlue    Colour = 105
	Yellow4           Colour = 106
	DarkSeaGreen      Colour = 108
	LightSkyBlue3     Colour = 110
	SkyBlue2          Colour = 111
	Chartreuse2       Colour = 112
	DarkSlateGray3    Colour = 116
	SkyBlue1          Colour = 117
	Chartreuse1       Colour = 118
	LightGreen        Colour = 120
	Aquamarine1       Colour = 122
	DarkSlateGray1    Colour = 123
	DeepPink4         Colour = 125
	MediumVioletRed   Colour = 126
	DarkViolet        Colour = 128
	MediumOrchid3     Colour = 133
	MediumOrchid      Colour = 134
	DarkGoldenrod     Colour = 136
	RosyBrown         Colour = 138
	Grey63            Colour = 139
	MediumPurple2     Colour = 140
	MediumPurple1     Colour = 141
	DarkKhaki         Colour = 143
	NavajoWhite3      Colour = 144
	Grey69            Colour = 145
	LightSteelBlue3   Colour = 146
	LightSteelBlue    Colour = 147
	DarkOliveGreen3   Colour = 149
	DarkSeaGreen3     Colour = 150
	LightCyan3        Colour = 152
	LightSkyBlue1     Colour = 153
	GreenYellow       Colour = 154
	DarkOliveGreen2   Colour = 155
	PaleGreen1        Colour = 156
	DarkSeaGreen1     Colour = 158
	PaleTurquoise1    Colour = 159
	Red3              Colour = 160
	DeepPink3         Colour = 162
	Magenta3          Colour = 164
	DarkOrange3       Colour = 166
	IndianRed         Colour = 167
	HotPink3          Colour = 168
	HotPink2          Colour = 169
	Orchid            Colour = 170
	Orange3           Colour = 172
	LightSalmon3      Colour = 173
	LightPink3        Colour = 174
	Pink3             Colour = 175
	Plum3             Colour = 176
	Violet            Colour = 177
	Gold3             Colour = 178
	LightGoldenrod3   Colour = 179
	Tan               Colour = 180
	MistyRose3        Colour = 181
	Thistle3          Colour = 182
	Plum2             Colour = 183
	Yellow3           Colour = 184
	Khaki3            Colour = 185
	LightYellow3      Colour = 187
	Grey84            Colour = 188
	LightSteelBlue1   Colour = 189
	Yellow2           Colour = 190
	DarkOliveGreen1   Colour = 192
	Honeydew2         Colour = 194
	LightCyan1        Colour = 195
	Red1              Colour = 196
	DeepPink2         Colour = 197
	DeepPink1         Colour = 199
	Magenta2          Colour = 200
	Magenta1          Colour = 201
	OrangeRed1        Colour = 202
	IndianRed1        Colour = 204
	HotPink           Colour = 206
	MediumOrchid1     Colour = 207
	DarkOrange        Colour = 208
	Salmon1           Colour = 209
	LightCoral        Colour = 210
	PaleVioletRed1    Colour = 211
	Orchid2           Colour = 212
	Orchid1           Colour = 213
	Orange1           Colour = 214
	SandyBrown        Colour = 215
	LightSalmon1      Colour = 216
	LightPink1        Colour = 217
	Pink1             Colour = 218
	Plum1             Colour = 219
	Gold1             Colour = 220
	LightGoldenrod2   Colour = 222
	NavajoWhite1      Colour = 223
	MistyRose1        Colour = 224
	Thistle1          Colour = 225
	Yellow1           Colour = 226
	LightGoldenrod1   Colour = 227
	Khaki1            Colour = 228
	Wheat1            Colour = 229
	Cornsilk1         Colour = 230
	Grey100           Colour = 231
	Grey3             Colour = 232
	Grey7             Colour = 233
	Grey11            Colour = 234
	Grey15            Colour = 235
	Grey19            Colour = 236
	Grey23            Colour = 237
	Grey27            Colour = 238
	Grey30            Colour = 239
	Grey35            Colour = 240
	Grey39            Colour = 241
	Grey42            Colour = 242
	Grey46            Colour = 243
	Grey50            Colour = 244
	Grey54            Colour = 245
	Grey58            Colour = 246
	Grey62            Colour = 247
	Grey66            Colour = 248
	Grey70            Colour = 249
	Grey74            Colour = 250
	Grey78            Colour = 251
	Grey82            Colour = 252
	Grey85            Colour = 253
	Grey89            Colour = 254
	Grey93            Colour = 255
)

type AnsiPrinter struct {
	enabled bool
	content strings.Builder
}

// NewAnsi creates a new AnsiPrinter. It doesn't assume anything about the device that the output will be
// directed to.
func NewAnsi() *AnsiPrinter {
	return &AnsiPrinter{enabled: true}
}

// Ansi is a default instance of AnsiPrinter
var Ansi = NewAnsi()

// NewAnsiFor creates a new AnsiPrinter for a given device. It will not automatically print to this device,
// but it will disable ANSI colors if the device doesn't seem to support them, like when redirecting
// standard output into a file or piping it to another program
func NewAnsiFor(f *os.File) *AnsiPrinter {
	o, err := f.Stat()
	if err != nil {
		panic(err)
	}
	enabled := (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice
	return &AnsiPrinter{enabled: enabled}
}

// Clear resets the internal buffer, after calling it the AnsiPrinter is in a clean state,
// as if just created
func (ap *AnsiPrinter) Clear() *AnsiPrinter {
	ap.content = strings.Builder{}
	return ap
}

// GetBuffer retrieves internal buffer as a string without clearing it
func (ap *AnsiPrinter) GetBuffer() string {
	return ap.content.String()
}

// IsEnabled returns true if color output is enabled
func (ap *AnsiPrinter) IsEnabled() bool {
	return ap.enabled
}

// SetEnabled enables or disables the color output. This does not affect the string already in AnsiPrinter's buffer
func (ap *AnsiPrinter) SetEnabled(value bool) {
	ap.enabled = value
}

// String converts the internal buffer to a string. The buffer is cleared after the call
func (ap *AnsiPrinter) String() string {
	s := ap.content.String()
	ap.Clear()
	return s
}

// Reset resets all the colors and attributes to defaults
func (ap *AnsiPrinter) Reset() *AnsiPrinter {
	ap.writeAnsiSeq(0)
	return ap
}

// Fg sets foreground color. When using SysColor constants, like SysRed or SysYellow, the most basic and most compatible
// ANSI sequence will be used. Using any of other color constants or integer values will use 256-color ANSI sequence
func (ap *AnsiPrinter) Fg(color Colour) *AnsiPrinter {
	if color < 0 {
		ap.writeAnsiSeq(30 + (-color))
	} else {
		ap.writeAnsiSeq(38, 5, color)
	}
	return ap
}

// Bg sets background color to one of standard 8 colors
func (ap *AnsiPrinter) Bg(color Colour) *AnsiPrinter {
	if color < 0 {
		ap.writeAnsiSeq(40 + (-color))
	} else {
		ap.writeAnsiSeq(48, 5, color)
	}
	return ap
}

// FgHi sets foreground color to the high intensity version of one of standard 8 colors
//
// If used with one of 256 color codes, it will just set the color, without modifying the intensity
func (ap *AnsiPrinter) FgHi(color Colour) *AnsiPrinter {
	if color < 0 {
		color = -color
		ap.writeAnsiSeq(90 + color)
		return ap
	} else {
		return ap.Fg(color)
	}
}

func (ap *AnsiPrinter) Attr(attr Attribute) *AnsiPrinter {
	ap.writeAnsiSeq(attr)
	return ap
}

// BgHi sets background color to the high intensity version of one of standard 8 colors
//
// If used with one of 256 color codes, it will just set the color, without modifying the intensity
func (ap *AnsiPrinter) BgHi(color Colour) *AnsiPrinter {
	if color < 0 {
		color = -color
		ap.writeAnsiSeq(100 + color)
		return ap
	} else {
		return ap.Bg(color)
	}
}

// FgRgb sets foreground color using "true color" RGB color
func (ap *AnsiPrinter) FgRgb(r, g, b int) *AnsiPrinter {
	ap.writeAnsiSeq(38, 2, r, g, b)
	return ap
}

// FgRgbI sets foreground color using "true color" RGB color represented as a single integer
func (ap *AnsiPrinter) FgRgbI(i int) *AnsiPrinter {
	r := (i >> 16) & 0xFF
	g := (i >> 8) & 0xFF
	b := i & 0xFF
	ap.writeAnsiSeq(38, 2, r, g, b)
	return ap
}

// BgRgb sets foreground color using "true color" RGB color
func (ap *AnsiPrinter) BgRgb(r, g, b int) *AnsiPrinter {
	ap.writeAnsiSeq(48, 2, r, g, b)
	return ap
}

// BgRgbI sets background color using "true color" RGB color represented as a single integer
func (ap *AnsiPrinter) BgRgbI(i int) *AnsiPrinter {
	r := (i >> 16) & 0xFF
	g := (i >> 8) & 0xFF
	b := i & 0xFF
	ap.writeAnsiSeq(48, 2, r, g, b)
	return ap
}

// FgRgb3 sets foreground RGB color converted to 9-bit color (3;3;3) as supported by 256-color ANSI sequence
func (ap *AnsiPrinter) FgRgb3(r, g, b int) *AnsiPrinter {
	color := ap.threeBitColorCube(r, g, b)
	return ap.Fg(color)
}

// BgRgb3 sets background RGB color converted to 9-bit color (3;3;3) as supported by 256-color ANSI sequence
func (ap *AnsiPrinter) BgRgb3(r, g, b int) *AnsiPrinter {
	color := ap.threeBitColorCube(r, g, b)
	return ap.Bg(color)
}

// FgRgb3I sets foreground RGB color, represented as integer, converting it to 9-bit color (3;3;3) as supported by 256-color ANSI sequence
func (ap *AnsiPrinter) FgRgb3I(rgb int) *AnsiPrinter {
	r := (rgb >> 16) & 0xFF
	g := (rgb >> 8) & 0xFF
	b := rgb & 0xFF
	color := ap.threeBitColorCube(r, g, b)
	return ap.Fg(color)
}

// BgRgb3I sets background RGB color, represented as integer, converting it to 9-bit color (3;3;3) as supported by 256-color ANSI sequence
func (ap *AnsiPrinter) BgRgb3I(rgb int) *AnsiPrinter {
	r := (rgb >> 16) & 0xFF
	g := (rgb >> 8) & 0xFF
	b := rgb & 0xFF
	color := ap.threeBitColorCube(r, g, b)
	return ap.Bg(color)
}

// FgGray sets foreground color that is the shade of gray. intensity is a value in a range [0..1]. It is converted to
// one of standard 24 gray shades in the 256-color palette
func (ap *AnsiPrinter) FgGray(intensity float64) *AnsiPrinter {
	gray := ap.shadeOfGrayColor(intensity)
	return ap.Fg(gray)
}

// BgGray sets background color that is the shade of gray. intensity is a value in a range [0..1]. It is converted to
// one of standard 24 gray shades in the 256-color palette
func (ap *AnsiPrinter) BgGray(intensity float64) *AnsiPrinter {
	gray := ap.shadeOfGrayColor(intensity)
	return ap.Bg(gray)
}

// A adds text to the AnsiPrinter's buffer. The text will be output with the current colors and attributes
func (ap *AnsiPrinter) A(text string) *AnsiPrinter {
	ap.content.WriteString(text)
	return ap
}

// S adds formatted (similar to fmt.Sprintf) text to the AnsiPrinter's buffer. The text will be output with the current colors and attributes
func (ap *AnsiPrinter) S(format string, params ...interface{}) *AnsiPrinter {
	ap.content.WriteString(fmt.Sprintf(format, params...))
	return ap
}

// CR adds carriage return character (ASCII 13) to the AnsiPrinter's buffer
func (ap *AnsiPrinter) CR() *AnsiPrinter {
	ap.content.WriteRune('\n')
	return ap
}

// Esc allows to add custom Esc sequence to the buffer
// The sequence that will be added is:
//
// ESC[codes sep codes sep codes sep command
//
// i.e. Esc('m', ':', 4, 3) will create sequence 'ESC[4:3m' which will create curly underline in iTerm2
func (ap *AnsiPrinter) Esc(command rune, sep rune, codes ...int) *AnsiPrinter {
	ap.writeAnsiCommand(command, sep, codes...)
	return ap
}

// EscM allows to add custom SGR sequences to the output
// EscM(38, 2, 0, 0, 255) will add sequence 'ESC[38;2;0;0;255m' to enable bright blue RGB color
func (ap *AnsiPrinter) EscM(codes ...int) *AnsiPrinter {
	ap.writeAnsiSeq(codes...)
	return ap
}

func (ap *AnsiPrinter) writeAnsiCommand(command rune, sep rune, codes ...int) {
	if ap.enabled {
		ap.content.WriteString(esc)
		for _, code := range codes {
			ap.content.WriteString(strconv.Itoa(code))
			ap.content.WriteRune(sep)
		}
		ap.content.WriteRune(command)
	}
}

func (ap *AnsiPrinter) writeAnsiSeq(codes ...int) {
	ap.writeAnsiCommand('m', ';', codes...)
}

func (ap *AnsiPrinter) threeBitColorCube(r int, g int, b int) int {
	r = (r & 0xFF) * 5 / 255
	g = (g & 0xFF) * 5 / 255
	b = (b & 0xFF) * 5 / 255
	color := 16 + 36*r + 6*g + b
	return color
}

func (ap *AnsiPrinter) shadeOfGrayColor(intensity float64) int {
	if intensity < 0 {
		intensity = 0
	}
	if intensity > 1 {
		intensity = 1
	}
	gray := 232 + int(23*intensity)
	return gray
}
