package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rumstead/kubectl-safe/pkg/exec"
	"k8s.io/klog/v2"
)

func Confirm(verb string, args []string) bool {
	context, err := exec.ExecCmd("kubectl", "config", "current-context")
	if err != nil {
		klog.Error(err)
		context = "<unknown>"
	}
	context = strings.TrimSpace(context)

	ns := extractNamespace(args)

	reader := bufio.NewReader(os.Stdin)
	if ns != "" {
		fmt.Printf("You are running a %s against context %s (namespace %s), continue? [yY] ", verb, context, ns)
	} else {
		fmt.Printf("You are running a %s against context %s, continue? [yY] ", verb, context)
	}
	input, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	return strings.ToLower(strings.TrimSpace(input)) == "y"
}

// extractNamespace pulls the namespace from -n or --namespace flags in args.
func extractNamespace(args []string) string {
	for i, arg := range args {
		if arg == "-n" || arg == "--namespace" {
			if i+1 < len(args) {
				return args[i+1]
			}
		}
		if strings.HasPrefix(arg, "--namespace=") {
			return strings.TrimPrefix(arg, "--namespace=")
		}
		if strings.HasPrefix(arg, "-n=") {
			return strings.TrimPrefix(arg, "-n=")
		}
	}
	return ""
}
