# ANSI - Go library for ANSI colours in terminal

This library API is inspired by [jansi library](https://github.com/fusesource/jansi).

## Usage

Get the library:

    go get github.com/uaraven/ansi

You will need Go 1.14 or higher to use it.

`ansi` supports basic 7-colour, 256 colour and true colour modes. You can also use various attributes.

Simple coloured output:

```go
import . "github.com/uaraven/ansi"

errorMsg := Ansi.A("Error: ").Fg(Red).S("File not found: %s", fileName).Reset().A("Try a different name").String()

```

Underlined text:
```go
import . "github.com/uaraven/ansi"

errorMsg := Ansi.A("This is ").Attr(Underline).A("important").String()

```

Disable colour if output is being redirected to file:

```go
import . "github.com/uaraven/ansi"

a := NewAnsiFor(os.stdout)

errorMsg := a.A("Error: ").Fg(Red).S("File not found: %s", fileName).Reset().A("Try a different name").String()
```

## Compatibility

### Basic colours

`ansi` supports basic terminal colours: Black, SysRed, SysGreen, SysYellow, SysBlue, SysMagenta, SysCyan and SysWhite.
These colours should be compatible with all terminals of the last 30 years.

You can use `AnsiPrinter.Fg()` and `AnsiPrinter.Bg()` to set these colours.

### 256 colours

Extended colours, like `Black`, `Red`, `Yellow` or `LightGoldenrod3` will use 256-colour ANSI sequences and should also 
work with any terminal of the last 30 years.

You also use `AnsiPrinter.Fg()` and `AnsiPrinter.Bg()` to set these colours.

You can set 256 colours using RGB values. 
The 24-bit RGB colours will be converted to 9-bit representation (3 bit per colour) suitable for 256-colour mode. 
Use following functions to set 9-bit colours: `AnsiPrinter.FgRgb3(r,g,b int)`, `AnsiPrinter.FgRgb3I(rgb int)`, `AnsiPrinter.BgRgb3(r,g,b int)` and `AnsiPrinter.BgRgb3I(rgb int)`

You can select shades of gray using intensity values for foreground and backround using 
`AnsiPrinter.FgGray(intensity)` and `AnsiPrinter.BgGray(intensity)`. `intensity` is a floating point value in the range
of 0 to 1. Any values beyond this range will be clipped.

### "True colour"

`ansi` supports full 24-bit colours using RGB values. The 24-bit RGB colours will be converted to 9-bit representation
(3 bit per colour) suitable for 256-colour mode.
Use following functions to set 9-bit colours: `AnsiPrinter.FgRgb(r,g,b int)`, `AnsiPrinter.FgRgbI(rgb int)`, `AnsiPrinter.BgRgb(r,g,b int)` and `AnsiPrinter.BgRgbI(rgb int)` 

### Colour names

`ansi` defines constants for the 256-colour palette with the names taken from [here](https://www.ditig.com/256-colors-cheat-sheet) and
matching Xterm names.

Note that some of the Xterm colour names are duplicated with the different palette index value. In the case of duplicates
one random name was selected and added to `ansi` constants.


## License

`ansi` is distributed under the terms of MIT license.

Copy of the license text is available in the [license.txt](license.txt) file.