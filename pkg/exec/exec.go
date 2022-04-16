package exec

import (
	"os"
	"os/exec"
	"syscall"
)

func KubeCtl(args []string) error {
	path, err := exec.LookPath("kubectl")
	if err != nil {
		return err
	}
	// https://man7.org/linux/man-pages/man2/execve.2.html
	// first arg in the argv list is the binary path
	argv := append([]string{path}, args...)
	return syscall.Exec(path, argv, os.Environ())
}

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
