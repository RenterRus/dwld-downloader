package sensors

import (
	"fmt"
	"os/exec"
)

// Install and configure
// sudo apt install lm-sensors
// sudo sensors-detect

func Sensors() string {
	var out []byte
	var err error

	if out, err = exec.Command("/bin/sh", "-c", "sensors").Output(); err != nil {
		return fmt.Sprintf("sensors error: %s", err.Error())
	}

	return string(out)
}
