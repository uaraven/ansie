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
	if err.Cause == nil {
		return fmt.Sprintf("screen error: %s", err.Message)
	}
	return fmt.Sprintf("screen error: %s, caused by: %v", err.Message, err.Cause)
}

// Terminal interface defines abstract terminal operations.
type Terminal interface {
	// Fd returns the file descriptor of the terminal.
	Fd() int
	// Write writes data to the terminal.
	Write(s string) (n int, err error)
	// IsTerminal checks if the file descriptor is a terminal.
	IsTerminal() bool
	// GetState retrieves the current terminal state.
	GetState() (*unix.Termios, error)
	// SetState sets the terminal state to the provided termios structure.
	SetState(termState *unix.Termios) error
	// GetSize retrieves the size of the terminal.
	GetSize() (*unix.Winsize, error) // Returns width and height of the terminal
}

// FileTerminal implements the Terminal interface backed by an os.File.
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

func (t *FileTerminal) Write(s string) (n int, err error) {
	return t.file.WriteString(s)
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

// Screen represents a terminal screen with methods to manipulate the terminal display.
type Screen struct {
	terminal Terminal
	// Width represents the width of the terminal in characters. It is updated on resize.
	Width int
	// Height represents the height of the terminal in characters. It is updated on resize.
	Height int
	// CursorVisible indicates whether the cursor is currently visible.
	CursorVisible  bool
	signals        chan os.Signal
	initialTermios unix.Termios
	closed         atomic.Bool
}

// NewScreen initializes a new Screen using the standard output file descriptor,
func NewScreen() (*Screen, error) {
	return NewScreenFromFile(os.Stdout)
}

// NewScreenFromTerminal initializes a new Screen using the provided Terminal interface,
func NewScreenFromTerminal(term Terminal) (*Screen, error) {
	if !term.IsTerminal() {
		return nil, NewScreenError(fmt.Sprintf("File descriptor %d is not a valid terminal", term.Fd()), nil)
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
	screen.Clear()
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
	return screen, nil
}

// NewScreenFromFile initializes a new Screen using file descriptor of f parameter,
// switches the terminal into alternate buffer mode, and retrieves the terminal size.
func NewScreenFromFile(f *os.File) (*Screen, error) {
	term, err := NewTerminalFromFile(f)
	if err != nil {
		return nil, NewScreenError("Cannot create terminal for invalid file descriptor", err)
	}
	return NewScreenFromTerminal(term)
}

func (s *Screen) writeEsc(command string) {
	_, _ = s.terminal.Write(esc)
	_, _ = s.terminal.Write(command)
}

func (s *Screen) resize() error {
	winSize, err := s.terminal.GetSize()
	if err != nil {
		return NewScreenError("Cannot get window size", err)
	}
	s.Width = int(winSize.Col)
	s.Height = int(winSize.Row)
	return nil
}

// HideCursor hides the cursor in the terminal.
func (s *Screen) HideCursor() {
	s.CursorVisible = false
	s.writeEsc("?25l") // Hide cursor
}

// ShowCursor shows the cursor in the terminal.
func (s *Screen) ShowCursor() {
	s.CursorVisible = true
	s.writeEsc("?25h") // Show cursor
}

// MoveCursorTo moves the cursor to the specified (x, y) position in the terminal.
// Coordinates are 1-based, where (1, 1) is the top-left corner
func (s *Screen) MoveCursorTo(x, y int) {
	if x < 1 || y < 1 || x > s.Width || y > s.Height {
		return // Invalid coordinates
	}
	s.writeEsc(fmt.Sprintf("%d;%dH", y, x)) // Move cursor to (x, y)
}

// Clear clears the terminal screen and moves the cursor to the home position.
func (s *Screen) Clear() {
	s.writeEsc("2J") // Clear the screen
	s.writeEsc("H")  // Move cursor to home position
}

// Close closes the screen, restores the terminal state, and exits alternate buffer mode.
func (s *Screen) Close() {
	if s.closed.CompareAndSwap(false, true) {
		close(s.signals)
		s.enterAlternateBuffer()
		s.ShowCursor()
		s.exitAlternateBuffer()
		_ = unix.IoctlSetTermios(s.terminal.Fd(), setTermios, &s.initialTermios) // Restore terminal state
	}
}

func (s *Screen) enterAlternateBuffer() {
	s.writeEsc("?1049h")
}

func (s *Screen) exitAlternateBuffer() {
	s.writeEsc("?1049l")
}

// SetRawMode sets the terminal to raw mode or restores it to normal mode.
// In raw mode, input is not processed (no echo, no line buffering).
// This is useful for applications that need to handle input directly, like text editors or games.
func (s *Screen) SetRawMode(rawMode bool) error {
	termState := s.initialTermios
	if rawMode {
		termState.Lflag &^= unix.ECHO | unix.ICANON // Disable echo and canonical mode
	} else {
		termState.Lflag |= unix.ECHO | unix.ICANON // Enable echo and canonical mode
	}
	return s.terminal.SetState(&termState) // Ignore error for simplicity
}
