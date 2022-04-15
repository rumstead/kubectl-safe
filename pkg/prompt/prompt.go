package prompt

import (
	"bufio"
	"fmt"
	"github.com/rumstead/kubectl-safe/pkg/exec"
	"k8s.io/klog/v2"
	"os"
	"strings"
)

func Confirm(verb string) bool {
	output, err := exec.ExecCmd("kubectl", "config", "current-context")
	if err != nil {
		klog.Error(err)
	}
	output = strings.ReplaceAll(output, "\n", "")
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("You are running a %s against context %s, continue? [yY] ", verb, output)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	return strings.ToLower(input) == "y"
}
