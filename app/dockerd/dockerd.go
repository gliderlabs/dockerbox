package dockerd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/client"
	"github.com/gliderlabs/com"
	"github.com/gliderlabs/com/config"
	"github.com/gliderlabs/dockerbox/lib/subprocess"
)

func init() {
	com.DefaultRegistry.Register(&com.Object{Value: &Component{}})
}

type Observer interface {
	DockerReady(*Component)
}

type Component struct {
	Observers []Observer `com:"extpoint"`
	Options   map[string]string
	Client    *client.Client

	*subprocess.Subprocess
}

func (c *Component) InitializeConfig(config config.Settings) error {
	c.Options = make(map[string]string)
	return config.Load(&(c.Options))
}

func (c *Component) InitializeDaemon() (err error) {
	var options []string
	for k, v := range c.Options {
		options = append(options, fmt.Sprintf("--%s=%s", k, v))
	}
	c.Subprocess = subprocess.NewSubprocess("dockerd", options...)
	c.Subprocess.Stdout = os.Stdout
	c.Subprocess.Stderr = os.Stderr
	c.Client, err = client.NewClient(c.Options["host"], "", nil, nil)
	return
}

func (c *Component) Serve() {
	go func() {
		for {
			_, err := c.Client.ServerVersion(context.Background())
			if err == nil {
				break
			}
			time.Sleep(20 * time.Millisecond)
		}
		for _, observer := range c.Observers {
			observer.DockerReady(c)
		}
	}()
	c.Subprocess.Serve()
}
