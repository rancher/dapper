package file

import (
	"fmt"
)

func (d *Dapperfile) vSocket() string {
	return fmt.Sprintf("%s://./pipe/docker_engine", d.env.HostSocket())
}
