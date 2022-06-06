package safe

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"k8s.io/klog/v2"

	"github.com/rumstead/kubectl-safe/pkg/exec"
)

func IsSafe(verb string, args []string) (bool, error) {
	isContextSafe, err := isContextSafe()
	if err != nil {
		return false, err
	}
	if isContextSafe {
		return true, nil
	}
	isVerbSafe, err := isVerbSafe(verb)
	if err != nil {
		return false, err
	}
	isDryRun := isDryRun(args)
	return isVerbSafe || isDryRun, nil
}

func isContextSafe() (bool, error) {
	context, err := exec.ExecCmd("kubectl", []string{"config", "current-context"}...)
	if err != nil {
		return false, err
	}

	safeContexts, err := getSafeContexts()
	if err != nil {
		return false, err
	}
	return safeContexts.Contains(context), nil
}

func getSafeContexts() (*KubeCtlSafeMap, error) {
	if contexts := os.Getenv(KubectlSafeContexts); contexts != "" {
		return parseSafeConfig(contexts)
	}
	return &DefaultSafeContexts, nil
}

func isVerbSafe(verb string) (bool, error) {
	if verb == "" {
		return true, nil
	}
	commands, err := getSafeCommands()
	if err != nil {
		return false, err
	}
	return commands.Contains(verb), nil
}

func isDryRun(cmd []string) bool {
	for _, c := range cmd {
		if strings.Contains(c, "--dry-run") {
			return true
		}
	}
	return false
}

func getSafeCommands() (*KubeCtlSafeMap, error) {
	if commands := os.Getenv(KubectlSafeCommands); commands != "" {
		return parseSafeConfig(commands)
	}
	return &DefaultSafeCommands, nil
}

func parseSafeConfig(safeConfig string) (*KubeCtlSafeMap, error) {
	// safeConfig can be a csv of strings or a file path
	_, err := os.Stat(safeConfig)
	switch {
	case os.IsNotExist(err):
		cmds := strings.Split(safeConfig, ",")
		if len(cmds) >= 1 {
			return getConfigFromSlice(cmds)
		}
		return nil, fmt.Errorf("%s was provided but its neither a word or existing file", safeConfig)
	case err == nil:
		return getCommandsFromFile(safeConfig)
	default:
		return nil, err
	}
}

func getCommandsFromFile(commands string) (*KubeCtlSafeMap, error) {
	readCommands := NewCommands()
	file, err := os.Open(commands)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	klog.V(3).Infof("reading commands from %s.\n", commands)
	scanner := bufio.NewScanner(file)
	// scanner has a 64k limit on lines... i hope that is never reached
	for scanner.Scan() {
		token := scanner.Text()
		klog.V(3).Infof("adding %s command to the safe list.\n", token)
		readCommands.Add(token)
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}
	return &readCommands, nil
}

func getConfigFromSlice(cmds []string) (*KubeCtlSafeMap, error) {
	commands := NewCommands()
	for i := range cmds {
		commands.Add(cmds[i])
	}
	return &commands, nil
}

func NewCommands() KubeCtlSafeMap {
	return KubeCtlSafeMap{set: make(map[string]Void)}
}
