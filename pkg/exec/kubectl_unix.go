//go:build !windows

package exec

import (
	"os"
	"os/exec"
	"syscall"
)

// KubeCtl replaces the current process with kubectl via execve(2).
// We use syscall.Exec instead of os/exec so that kubectl fully replaces this process:
// signals (Ctrl+C), exit codes, and the process name all pass through transparently
// to the caller without needing manual forwarding or leaving a parent process resident.
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
