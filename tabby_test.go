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

func (t TestApplication) Main(args Arguments) (*TabbyContainer, error) {
	fmt.Println(args.Get("aa").(int))
	return nil, nil
}

func NewTestApplication() *TestApplication {
	return &TestApplication{
		NewBaseApplication(true, nil),
	}
}
func TestTabby(t *testing.T) {
	ta := NewTestApplication()
	ta.SetParam("aa", "", Int(10))
	tb := NewTabby("test", ta)
	_, err := tb.Run([]string{"-aa", "30"})
	if err != nil {
		fmt.Println(err)
	}
}
