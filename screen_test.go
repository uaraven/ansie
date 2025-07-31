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
	g.Expect(s.Width).To(Equal(80), "Expected terminal width to be non-zero")
	g.Expect(s.Height).To(Equal(24), "Expected terminal width to be non-zero")
	state, err := m.GetState()
	g.Expect(err).To(BeNil(), "Expected no error when reading terminal state")
	g.Expect(state.Lflag&unix.ECHO).ToNot(Equal(unix.ECHO), "Expected terminal to not have ECHO flag set")
	g.Expect(state.Lflag&unix.ICANON).ToNot(Equal(unix.ICANON), "Expected terminal to not have ICANON flag set")

	g.Expect(m.Buffer.String()).To(ContainSubstring("[?1049h"), "Expected no output in mock terminal buffer")
}

func TestCloseScreen(t *testing.T) {
	g := NewGomegaWithT(t)
	m := NewMockTerminal(80, 24)
	s, err := NewScreenFromTerminal(m)
	g.Expect(err).To(BeNil(), "Expected no error when creating a new screen")
	g.Expect(m.Buffer.String()).To(ContainSubstring("[?1049h"), "Expected switch to alternative buffer")
	g.Expect(m.Buffer.String()).ToNot(ContainSubstring("[?25l"), "Expected cursor not to be hidden")
	s.Close()
	g.Expect(m.Buffer.String()).To(ContainSubstring("[?1049l"), "Expected switch from alternative buffer")
	g.Expect(m.Buffer.String()).ToNot(ContainSubstring("[2J"), "Expected screen not cleared")
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

func TestCursorManipulation(t *testing.T) {
	g := NewGomegaWithT(t)
	m := NewMockTerminal(80, 24)
	s, err := NewScreenFromTerminal(m)
	if err == nil {
		defer s.Close()
	}
	s.HideCursor()
	g.Expect(m.Buffer.String()).To(ContainSubstring("[?25l"), "Expected cursor to be hidden")
	s.ShowCursor()
	g.Expect(m.Buffer.String()).To(ContainSubstring("[?25h"), "Expected cursor to be shown")
	s.MoveCursorTo(10, 5)
	g.Expect(m.Buffer.String()).To(ContainSubstring("[5;10H"), "Expected cursor to move to (10, 5)")
}
