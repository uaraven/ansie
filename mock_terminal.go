package ansie

import (
	"os"
	"strings"

	"golang.org/x/sys/unix"
)

type MockTerminal struct {
	FileDesc      int
	Width         int
	Height        int
	CursorX       int
	CursorY       int
	CursorVisible bool
	State         unix.Termios
	Buffer        strings.Builder
}

func NewMockTerminal(width, height int) *MockTerminal {
	return &MockTerminal{
		FileDesc:      1,
		Width:         width,
		Height:        height,
		CursorX:       1,
		CursorY:       1,
		CursorVisible: true,
		State: unix.Termios{
			Lflag: unix.ECHO | unix.ICANON,
		},
	}
}

// Fd implements Terminal.
func (m *MockTerminal) Fd() int {
	return m.FileDesc
}

// GetState implements Terminal.
func (m *MockTerminal) GetState() (*unix.Termios, error) {
	return &m.State, nil
}

// IsTerminal implements Terminal.
func (m *MockTerminal) IsTerminal() bool {
	return true
}

// SetState implements Terminal.
func (m *MockTerminal) SetState(termState *unix.Termios) error {
	m.State = *termState
	return nil
}

// Write implements Terminal.
func (m *MockTerminal) Write(p []byte) (n int, err error) {
	return m.Buffer.Write(p)
}

// GetSize implements Terminal.
func (m *MockTerminal) GetSize() (*unix.Winsize, error) {
	return &unix.Winsize{
		Row:    uint16(m.Height),
		Col:    uint16(m.Width),
		Xpixel: 0,
		Ypixel: 0,
	}, nil
}

func (m *MockTerminal) SetSize(w, h int, c chan os.Signal) {
	m.Width = w
	m.Height = h
	if c != nil {
		c <- os.Signal(unix.SIGWINCH) // Simulate a signal for testing purposes
	}
}

var _ Terminal = (*MockTerminal)(nil)
