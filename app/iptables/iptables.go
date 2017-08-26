package iptables

import (
	api "github.com/coreos/go-iptables/iptables"
	"github.com/gliderlabs/com"
	"github.com/gliderlabs/com/config"
	"github.com/gliderlabs/dockerbox/app/dockerd"
)

func init() {
	com.DefaultRegistry.Register(&com.Object{Value: &Component{}})
}

type Component struct {
	Drop []string

	iptables *api.IPTables
}

func (c *Component) InitializeConfig(settings config.Settings) error {
	return settings.Load(&c)
}

func (c *Component) InitializeDaemon() (err error) {
	c.iptables, err = api.New()
	return
}

func (c *Component) DockerReady(docker *dockerd.Component) {
	c.iptables.Insert("filter", "INPUT", 1, "-i", "docker0", "-j", "DROP")
	for _, block := range c.Drop {
		c.iptables.Insert("filter", "FORWARD", 1, "-d", block, "!", "-o", "docker0", "-j", "DROP")
	}
	c.iptables.Insert("filter", "FORWARD", 1, "-d", docker.Options["dns"], "-j", "ACCEPT")
}
