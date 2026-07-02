//go:build windows

package exec

import (
	"os"
	"os/exec"
)

func KubeCtl(args []string) error {
	path, err := exec.LookPath("kubectl")
	if err != nil {
		return err
	}
	cmd := exec.Command(path, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
