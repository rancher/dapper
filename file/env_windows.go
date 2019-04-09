package file

import (
	"os"
	"strings"
)

func (c Context) HostSocket() string {
	s := os.Getenv("DOCKER_HOST")
	if strings.HasPrefix(s, "npipe://") {
		return strings.TrimPrefix(s, "npipe://")
	}
	return "//./pipe/docker_engine"
}
