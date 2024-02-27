package tabby

import (
	"fmt"
	"testing"
)

type TestApplication struct {
	*BaseApplication
}

func (t TestApplication) Detail() (string, string) {
	return "test", "Test App"
}

func (t TestApplication) Main(args Arguments) error {
	fmt.Println("ok")
	return nil
}

func NewTestApplication() *TestApplication {
	return &TestApplication{
		NewBaseApplication(nil),
	}
}
func TestTabby(t *testing.T) {
	tb := NewTabby("test", NewTestApplication())
	tb.Run(nil)
}
