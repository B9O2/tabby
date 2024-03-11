package tabby

import (
	"github.com/B9O2/canvas/containers"
	"github.com/B9O2/canvas/pixel"
)

type TabbyContainer struct {
	width, height uint
	c             containers.Container
}

func (tc *TabbyContainer) Display(space pixel.Pixel) error {
	pm, err := tc.c.Draw(tc.width, tc.height)
	if err != nil {
		return err
	}
	return pm.Display(space)
}

func NewTabbyContainer(width, height uint, c containers.Container) *TabbyContainer {
	return &TabbyContainer{
		width:  width,
		height: height,
		c:      c,
	}
}
