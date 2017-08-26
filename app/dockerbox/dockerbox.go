package dockerbox

import (
	"fmt"

	"github.com/gliderlabs/com"
)

func init() {
	com.DefaultRegistry.Register(&com.Object{Value: &Component{}})
}

var Version string

type Component struct{}

func (c *Component) InitializeDaemon() error {
	fmt.Println("starting dockerbox", Version)
	return nil
}
