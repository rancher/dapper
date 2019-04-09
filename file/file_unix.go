// +build linux freebsd openbsd darwin

package file

import (
	"fmt"
)

func (d *Dapperfile) vSocket() string {
	return fmt.Sprintf("%s:/var/run/docker.sock", d.env.HostSocket())
}
