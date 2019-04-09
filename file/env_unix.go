// +build linux freebsd openbsd darwin

package file

import (
	"os"
	"strings"
)

func (c Context) HostSocket() string {
	s := os.Getenv("DOCKER_HOST")
	if strings.HasPrefix(s, "unix://") {
		return strings.TrimPrefix(s, "unix://")
	}
	return "/var/run/docker.sock"
}
