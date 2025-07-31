package ansie

import (
    . "github.com/onsi/gomega"
    "testing"
)

func TestNewScreen(t *testing.T) {
    g := NewGomegaWithT(t)
    s, err := NewScreen()
    if err == nil {
        defer s.Close()
    }
    g.Expect(err).To(BeNil(), "Expected no error when creating a new screen")
    g.Expect(s.Width).ToNot(Equal(0), "Expected terminal width to be non-zero")
    g.Expect(s.Height).ToNot(Equal(0), "Expected terminal width to be non-zero")
}
