# ANSI - Go library for ANSI colors in terminal

This library API is inspired by [jansi library](https://github.com/fusesource/jansi).

## Usage

Get the library:

    go get github.com/uaraven/ansi

You will need Go 1.14 or higher to use it.

`ansi` supports basic 7-color, 256 color and true color modes. You can also use various attributes.

Simple colored output:

```go
import . "github.com/uaraven/ansi"

errorMsg := A.A("Error: ").Fg(Red).S("File not found: %s", fileName).Reset().A("Try a different name").String()

```

Underlined text:
```go
import . "github.com/uaraven/ansi"

errorMsg := A.A("This is ").Attr(Underline).A("important").String()

```


Disable color if output is being redirected to file:

```go
import . "github.com/uaraven/ansi"

a := AnsiFor(os.stdout)

errorMsg := a.A("Error: ").Fg(Red).S("File not found: %s", fileName).Reset().A("Try a different name").String()
```