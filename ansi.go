// Package ansi
//
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

type Color = int

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
	Black Color = 0
	// SysRed is a color from standard 7-color terminal palette
	SysRed Color = -1
	// SysGreen is a color from standard 7-color terminal palette
	SysGreen Color = -2
	// SysYellow is a color from standard 7-color terminal palette
	SysYellow Color = -3
	// SysBlue is a color from standard 7-color terminal palette
	SysBlue Color = -4
	// SysMagenta is a color from standard 7-color terminal palette
	SysMagenta Color = -5
	// SysCyan is a color from standard 7-color terminal palette
	SysCyan Color = -6
	// SysWhite is a Red color from standard 7-color terminal palette
	SysWhite Color = -7

	Maroon            = 1
	Green             = 2
	Olive             = 3
	Navy              = 4
	Purple            = 5
	Teal              = 6
	Silver            = 7
	Grey              = 8
	Red               = 9
	Lime              = 10
	Yellow            = 11
	Blue              = 12
	Fuchsia           = 13
	Aqua              = 14
	White             = 15
	Grey0             = 16
	NavyBlue          = 17
	DarkBlue          = 18
	Blue3             = 20
	Blue1             = 21
	DarkGreen         = 22
	DeepSkyBlue4      = 25
	DodgerBlue3       = 26
	DodgerBlue2       = 27
	Green4            = 28
	SpringGreen4      = 29
	Turquoise4        = 30
	DeepSkyBlue3      = 32
	DodgerBlue1       = 33
	Green3            = 34
	SpringGreen3      = 35
	DarkCyan          = 36
	LightSeaGreen     = 37
	DeepSkyBlue2      = 38
	DeepSkyBlue1      = 39
	SpringGreen2      = 42
	Cyan3             = 43
	DarkTurquoise     = 44
	Turquoise2        = 45
	Green1            = 46
	SpringGreen1      = 48
	MediumSpringGreen = 49
	Cyan2             = 50
	Cyan1             = 51
	Purple4           = 55
	Purple3           = 56
	BlueViolet        = 57
	Grey37            = 59
	MediumPurple4     = 60
	SlateBlue3        = 62
	RoyalBlue1        = 63
	Chartreuse4       = 64
	PaleTurquoise4    = 66
	SteelBlue         = 67
	SteelBlue3        = 68
	CornflowerBlue    = 69
	DarkSeaGreen4     = 71
	CadetBlue         = 72
	SkyBlue3          = 74
	Chartreuse3       = 76
	PaleGreen3        = 77
	SeaGreen3         = 78
	Aquamarine3       = 79
	MediumTurquoise   = 80
	SteelBlue1        = 81
	SeaGreen2         = 83
	SeaGreen1         = 85
	DarkSlateGray2    = 87
	DarkRed           = 88
	DarkMagenta       = 91
	Orange4           = 94
	LightPink4        = 95
	Plum4             = 96
	MediumPurple3     = 98
	SlateBlue1        = 99
	Wheat4            = 101
	Grey53            = 102
	LightSlateGrey    = 103
	MediumPurple      = 104
	LightSlateBlue    = 105
	Yellow4           = 106
	DarkSeaGreen      = 108
	LightSkyBlue3     = 110
	SkyBlue2          = 111
	Chartreuse2       = 112
	DarkSlateGray3    = 116
	SkyBlue1          = 117
	Chartreuse1       = 118
	LightGreen        = 120
	Aquamarine1       = 122
	DarkSlateGray1    = 123
	DeepPink4         = 125
	MediumVioletRed   = 126
	DarkViolet        = 128
	MediumOrchid3     = 133
	MediumOrchid      = 134
	DarkGoldenrod     = 136
	RosyBrown         = 138
	Grey63            = 139
	MediumPurple2     = 140
	MediumPurple1     = 141
	DarkKhaki         = 143
	NavajoWhite3      = 144
	Grey69            = 145
	LightSteelBlue3   = 146
	LightSteelBlue    = 147
	DarkOliveGreen3   = 149
	DarkSeaGreen3     = 150
	LightCyan3        = 152
	LightSkyBlue1     = 153
	GreenYellow       = 154
	DarkOliveGreen2   = 155
	PaleGreen1        = 156
	DarkSeaGreen1     = 158
	PaleTurquoise1    = 159
	Red3              = 160
	DeepPink3         = 162
	Magenta3          = 164
	DarkOrange3       = 166
	IndianRed         = 167
	HotPink3          = 168
	HotPink2          = 169
	Orchid            = 170
	Orange3           = 172
	LightSalmon3      = 173
	LightPink3        = 174
	Pink3             = 175
	Plum3             = 176
	Violet            = 177
	Gold3             = 178
	LightGoldenrod3   = 179
	Tan               = 180
	MistyRose3        = 181
	Thistle3          = 182
	Plum2             = 183
	Yellow3           = 184
	Khaki3            = 185
	LightYellow3      = 187
	Grey84            = 188
	LightSteelBlue1   = 189
	Yellow2           = 190
	DarkOliveGreen1   = 192
	Honeydew2         = 194
	LightCyan1        = 195
	Red1              = 196
	DeepPink2         = 197
	DeepPink1         = 199
	Magenta2          = 200
	Magenta1          = 201
	OrangeRed1        = 202
	IndianRed1        = 204
	HotPink           = 206
	MediumOrchid1     = 207
	DarkOrange        = 208
	Salmon1           = 209
	LightCoral        = 210
	PaleVioletRed1    = 211
	Orchid2           = 212
	Orchid1           = 213
	Orange1           = 214
	SandyBrown        = 215
	LightSalmon1      = 216
	LightPink1        = 217
	Pink1             = 218
	Plum1             = 219
	Gold1             = 220
	LightGoldenrod2   = 222
	NavajoWhite1      = 223
	MistyRose1        = 224
	Thistle1          = 225
	Yellow1           = 226
	LightGoldenrod1   = 227
	Khaki1            = 228
	Wheat1            = 229
	Cornsilk1         = 230
	Grey100           = 231
	Grey3             = 232
	Grey7             = 233
	Grey11            = 234
	Grey15            = 235
	Grey19            = 236
	Grey23            = 237
	Grey27            = 238
	Grey30            = 239
	Grey35            = 240
	Grey39            = 241
	Grey42            = 242
	Grey46            = 243
	Grey50            = 244
	Grey54            = 245
	Grey58            = 246
	Grey62            = 247
	Grey66            = 248
	Grey70            = 249
	Grey74            = 250
	Grey78            = 251
	Grey82            = 252
	Grey85            = 253
	Grey89            = 254
	Grey93            = 255
)

type AnsiPrinter struct {
	enabled bool
	content strings.Builder
}

// Ansi creates a new AnsiPrinter. It doesn't assume anything about the device that the output will be
// directed to.
func Ansi() *AnsiPrinter {
	return &AnsiPrinter{enabled: true}
}

// A is a default instance of AnsiPrinter
var A = Ansi()

// AnsiFor creates a new AnsiPrinter for a given device. It will not automatically print to this device,
// but it will disable ANSI colors if the device doesn't seem to support them, like when redirecting
// standard output into a file or piping it to another program
func AnsiFor(f *os.File) *AnsiPrinter {
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
func (ap *AnsiPrinter) Fg(color Color) *AnsiPrinter {
	if color < 0 {
		ap.writeAnsiSeq(30 + (-color))
	} else {
		ap.writeAnsiSeq(38, 5, color)
	}
	return ap
}

// Bg sets background color to one of standard 8 colors
func (ap *AnsiPrinter) Bg(color Color) *AnsiPrinter {
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
func (ap *AnsiPrinter) FgHi(color Color) *AnsiPrinter {
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
func (ap *AnsiPrinter) BgHi(color Color) *AnsiPrinter {
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
