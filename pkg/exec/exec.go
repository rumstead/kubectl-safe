package exec

import (
	"os/exec"
)

func ExecCmd(binary string, args ...string) (string, error) {
	path, err := exec.LookPath(binary)
	if err != nil {
		return "", err
	}
	cmd := exec.Command(path, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
