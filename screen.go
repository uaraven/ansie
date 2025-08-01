package ansie

import (
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"

	"golang.org/x/sys/unix"
)

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

// Terminal interface defines abstract terminal operations.
type Terminal interface {
	// Fd returns the file descriptor of the terminal.
	Fd() int
	// Write writes data to the terminal.
	Write(p []byte) (n int, err error)
	// IsTerminal checks if the file descriptor is a terminal.
	IsTerminal() bool
	// GetState retrieves the current terminal state.
	GetState() (*unix.Termios, error)
	// SetState sets the terminal state to the provided termios structure.
	SetState(termState *unix.Termios) error
	// GetSize retrieves the size of the terminal.
	GetSize() (*unix.Winsize, error) // Returns width and height of the terminal
}

type FileTerminal struct {
	file *os.File
}

// GetSize implements Terminal.
func (t *FileTerminal) GetSize() (*unix.Winsize, error) {
	return unix.IoctlGetWinsize(t.Fd(), unix.TIOCGWINSZ)
}

var _ Terminal = (*FileTerminal)(nil)

func NewTerminalFromFile(file *os.File) (*FileTerminal, error) {
	if file == nil {
		return nil, fmt.Errorf("file cannot be nil")
	}
	return &FileTerminal{file: file}, nil
}

func (t *FileTerminal) Fd() int {
	return int(t.file.Fd())
}

func (t *FileTerminal) Write(p []byte) (n int, err error) {
	return t.file.Write(p)
}

func (t *FileTerminal) IsTerminal() bool {
	_, err := unix.IoctlGetTermios(t.Fd(), getTermios)
	return err == nil
}

func (t *FileTerminal) GetState() (*unix.Termios, error) {
	termState, err := unix.IoctlGetTermios(t.Fd(), getTermios)
	if err != nil {
		return nil, err
	}
	return termState, nil
}

func (t *FileTerminal) SetState(termState *unix.Termios) error {
	return unix.IoctlSetTermios(t.Fd(), setTermios, termState)
}

type Screen struct {
	terminal       Terminal
	Width          int
	Height         int
	CursorVisible  bool
	signals        chan os.Signal
	initialTermios unix.Termios
	closed         atomic.Bool
}

func NewScreen() (*Screen, error) {
	return NewScreenFromFile(os.Stdout)
}

func NewScreenFromTerminal(term Terminal) (*Screen, error) {
	if !term.IsTerminal() {
		return nil, NewScreenError("File descriptor is not a terminal", nil)
	}
	screen := &Screen{
		terminal:      term,
		CursorVisible: true,
		signals:       make(chan os.Signal, 1),
	}
	screen.closed.Store(false)
	termState, err := term.GetState()
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
			switch sig {
			case unix.SIGWINCH:
				_ = screen.resize()
			case unix.SIGTERM, unix.SIGINT:
				screen.Close()
				return
			}
		}
	}()
	termState.Lflag &^= unix.ECHO | unix.ICANON // Disable echo and canonical mode
	err = term.SetState(termState)
	if err != nil {
		screen.Close()
		return nil, NewScreenError("Cannot set terminal state", err)
	}

	return screen, nil
}

// NewScreen initializes a new Screen, switches terminal into alternate buffer mode, and retrieves the terminal size.
func NewScreenFromFile(f *os.File) (*Screen, error) {
	term, err := NewTerminalFromFile(f)
	if err != nil {
		return nil, NewScreenError("Cannot create terminal for invalid file descriptor", err)
	}
	return NewScreenFromTerminal(term)
}

func (s *Screen) writeEsc(command string) {
	_, _ = s.terminal.Write([]byte(esc))
	_, _ = s.terminal.Write([]byte(command))
}

func (s *Screen) resize() error {
	winsize, err := s.terminal.GetSize()
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
		s.ShowCursor()
		s.exitAlternateBuffer()
		_ = unix.IoctlSetTermios(int(s.terminal.Fd()), setTermios, &s.initialTermios) // Restore terminal state
	}
}

func (s *Screen) enterAlternateBuffer() {
	s.writeEsc("[?1049h")
}

func (s *Screen) exitAlternateBuffer() {
	s.writeEsc("[?1049l")
}
