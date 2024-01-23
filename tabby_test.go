package tabby

import "testing"

type TestApplication struct {
	*BaseApplication
}

func (t TestApplication) Name() string {
	return "Test Application"
}

func (t TestApplication) Main() error {

}

func NewTestApplication() *TestApplication {
	return &TestApplication{
		NewBaseApplication(nil),
	}
}
func TestTabby(t *testing.T) {

	tb := NewTabby("test", NewTestApplication())
	tb
}
