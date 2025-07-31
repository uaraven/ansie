package ansie

import (
    "fmt"
    "os"
    "os/signal"
    "sync/atomic"

    "golang.org/x/sys/unix"
)

const stdoutFd = 1

type ScreenError struct {
    Message string
    Cause   error
}

func NewScreenError(message string, cause error) ScreenError {
    return ScreenError{
        Message: message,
        Cause:   cause,
    }
}

func (err ScreenError) Error() string {
    return fmt.Sprintf("screen error: %s, caused by: %v", err.Message, err.Cause)
}

type Screen struct {
    Width          int
    Height         int
    CursorVisible  bool
    signals        chan os.Signal
    initialTermios unix.Termios
    closed         atomic.Bool
}

// NewScreen initializes a new Screen, switches terminal into alternate buffer mode, and retrieves the terminal size.
func NewScreen() (*Screen, error) {
    screen := &Screen{
        CursorVisible: true,
        signals:       make(chan os.Signal, 1),
    }
    screen.closed.Store(false)
    termState, err := unix.IoctlGetTermios(stdoutFd, getTermios)
    if err != nil {
        screen.Close()
        return nil, NewScreenError("Cannot get terminal state", err)
    }
    screen.initialTermios = *termState
    screen.enterAlternateBuffer()
    err = screen.resize()
    if err != nil {
        screen.Close()
        return nil, err
    }
    signal.Notify(screen.signals, unix.SIGWINCH, unix.SIGTERM, unix.SIGINT)
    go func() {
        for sig := range screen.signals {
            if sig == unix.SIGWINCH {
                _ = screen.resize()
            } else if sig == unix.SIGTERM || sig == unix.SIGINT {
                screen.Close()
                return
            }
        }
    }()
    termState.Lflag &^= unix.ECHO | unix.ICANON // Disable echo and canonical mode
    err = unix.IoctlSetTermios(stdoutFd, setTermios, termState)
    if err != nil {
        screen.Close()
        return nil, NewScreenError("Cannot set terminal state", err)
    }

    return screen, nil
}

func (s *Screen) writeEsc(command string) {
    _, _ = os.Stdout.Write([]byte(esc))
    _, _ = os.Stdout.Write([]byte(command))
}

func (s *Screen) resize() error {
    winsize, err := unix.IoctlGetWinsize(1, unix.TIOCGWINSZ)
    if err != nil {
        return NewScreenError("Cannot get window size", err)
    }
    s.Width = int(winsize.Col)
    s.Height = int(winsize.Row)
    return nil
}

func (s *Screen) HideCursor() {
    s.CursorVisible = false
    s.writeEsc("[?25l") // Hide cursor
}

func (s *Screen) ShowCursor() {
    s.CursorVisible = true
    s.writeEsc("[?25h") // Hide cursor
}

func (s *Screen) MoveCursorTo(x, y int) {
    if x < 1 || y < 1 || x > s.Width || y > s.Height {
        return // Invalid coordinates
    }
    s.writeEsc(fmt.Sprintf("[%d;%dH", y, x)) // Move cursor to (x, y)
}

func (s *Screen) Clear() {
    s.writeEsc("[2J") // Clear the screen
    s.writeEsc("[H")  // Move cursor to home position
}

func (s *Screen) Close() {
    if s.closed.CompareAndSwap(false, true) {
        close(s.signals)
        s.enterAlternateBuffer()
        s.Clear()
        s.ShowCursor()
        s.exitAlternateBuffer()
        _ = unix.IoctlSetTermios(stdoutFd, setTermios, &s.initialTermios) // Restore terminal state
    }
}

func (s *Screen) enterAlternateBuffer() {
    s.writeEsc("[?1049h")
}

func (s *Screen) exitAlternateBuffer() {
    s.writeEsc("[?1049l")
}
