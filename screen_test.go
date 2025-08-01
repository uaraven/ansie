package ansie

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"golang.org/x/sys/unix"
)

func TestNewScreen(t *testing.T) {
	g := NewGomegaWithT(t)
	m := NewMockTerminal(80, 24)
	s, err := NewScreenFromTerminal(m)
	if err == nil {
		defer s.Close()
	}
	g.Expect(err).To(BeNil(), "Expected no error when creating a new screen")
	g.Expect(s.Width).To(Equal(80), "Expected terminal width to be 80")
	g.Expect(s.Height).To(Equal(24), "Expected terminal height to be 24")
	state, err := m.GetState()
	g.Expect(err).To(BeNil(), "Expected no error when reading terminal state")
	g.Expect(state.Lflag&uint64(unix.ECHO)).ToNot(Equal(uint64(0)), "Expected terminal to have ECHO flag set")
	g.Expect(state.Lflag&uint64(unix.ICANON)).ToNot(Equal(uint64(0)), "Expected terminal to have ICANON flag set")

	g.Expect(m.Buffer.String()).To(ContainSubstring("\033[?1049h"), "Expected no output in mock terminal buffer")
}

func TestCloseScreen(t *testing.T) {
	g := NewGomegaWithT(t)
	m := NewMockTerminal(80, 24)
	s, err := NewScreenFromTerminal(m)
	g.Expect(err).To(BeNil(), "Expected no error when creating a new screen")
	g.Expect(m.Buffer.String()).To(ContainSubstring("\033[?1049h"), "Expected switch to alternative buffer")
	g.Expect(m.Buffer.String()).ToNot(ContainSubstring("\033[?25l"), "Expected cursor not to be hidden")
	s.SetCursorVisible(false)
	m.ResetBuffer() // Reset buffer to check after closing
	s.Close()
	g.Expect(m.Buffer.String()).To(ContainSubstring("\u001B[?1049l"), "Expected switch from alternative buffer")
	g.Expect(m.Buffer.String()).ToNot(ContainSubstring("\u001B[2J"), "Expected screen not cleared")
	g.Expect(m.Buffer.String()).To(ContainSubstring("[?25h"), "Expected cursor to be shown")
}

func TestScreenResize(t *testing.T) {
	g := NewGomegaWithT(t)
	m := NewMockTerminal(80, 24)
	s, err := NewScreenFromTerminal(m)
	if err == nil {
		defer s.Close()
	}
	m.SetSize(100, 30, s.signals)
	time.Sleep(100 * time.Millisecond) // Allow time for resize signal to be processed
	g.Expect(s.Width).To(Equal(100), "Expected terminal width to be updated")
	g.Expect(s.Height).To(Equal(30), "Expected terminal height to be updated")
}

func TestScreen_SetRawMode(t *testing.T) {
	g := NewGomegaWithT(t)
	m := NewMockTerminal(80, 24)
	s, err := NewScreenFromTerminal(m)
	if err == nil {
		defer s.Close()
	}
	_ = s.SetRawMode(true)
	state, err := m.GetState()
	g.Expect(err).To(BeNil(), "Expected no error when reading terminal state")
	g.Expect(state.Lflag&uint64(unix.ECHO)).To(Equal(uint64(0)), "Expected terminal to have ECHO flag cleared")
	g.Expect(state.Lflag&uint64(unix.ICANON)).To(Equal(uint64(0)), "Expected terminal to have ICANON flag cleared")

}

func TestScreen_SetCursorVisible(t *testing.T) {
	g := NewGomegaWithT(t)
	m := NewMockTerminal(80, 24)
	s, err := NewScreenFromTerminal(m)
	if err == nil {
		defer s.Close()
	}
	g.Expect(err).To(BeNil(), "Expected no error when creating a new screen")
	s.SetCursorVisible(false)
	g.Expect(m.Buffer.String()).To(ContainSubstring("\u001B[?25l"), "Expected cursor to be hidden")
	s.SetCursorVisible(true)
	g.Expect(m.Buffer.String()).To(ContainSubstring("\u001B[?25h"), "Expected cursor to be shown")
}

func TestScreen_MoveCursorTo(t *testing.T) {
	g := NewGomegaWithT(t)
	m := NewMockTerminal(80, 24)
	s, err := NewScreenFromTerminal(m)
	if err == nil {
		defer s.Close()
	}
	g.Expect(err).To(BeNil(), "Expected no error when creating a new screen")
	s.MoveCursorTo(10, 5)
	g.Expect(m.Buffer.String()).To(ContainSubstring("\u001B[5;10H"), "Expected cursor to move to (10, 5)")
	m.ResetBuffer()
	s.MoveCursorTo(0, 5)
	g.Expect(m.Buffer.String()).ToNot(ContainSubstring("\u001B["), "Expected cursor to not move")
	m.ResetBuffer()
	s.MoveCursorTo(20, 25)
	g.Expect(m.Buffer.String()).ToNot(ContainSubstring("\u001B["), "Expected cursor to not move")
	s.MoveCursorTo(80, 24)
	g.Expect(m.Buffer.String()).To(ContainSubstring("\u001B[24;80H"), "Expected cursor to move to (80, 24)")
}
