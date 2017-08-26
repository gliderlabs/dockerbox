package main

import (
	"github.com/gliderlabs/com/daemon"

	"github.com/gliderlabs/dockerbox/app/dockerbox"
	_ "github.com/gliderlabs/dockerbox/app/dockerd"
	_ "github.com/gliderlabs/dockerbox/app/iptables"
)

var Version string

func main() {
	dockerbox.Version = Version
	daemon.Run("dockerbox")
}
