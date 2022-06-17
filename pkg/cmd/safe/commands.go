package safe

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
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

func isVerbSafe(verb string) (bool, error) {
	if verb == "" {
		return true, nil
	}

	// first check if the verb is configured as unsafe
	unsafeCommands, err := getUnSafeCommands()
	if err != nil {
		return false, err
	}
	if len(unsafeCommands.set) > 0 {
		// if unsafeCommands contains the verb, it is unsafe
		return !unsafeCommands.Contains(verb), nil
	}
	// second check if it configured as safe
	safeCommands, err := getSafeCommands()
	if err != nil {
		return false, err
	}
	return safeCommands.Contains(verb), nil
}

func getSafeContexts() (*KubeCtlSafeMap, error) {
	commands, err := parseCommands(KubectlSafeContexts)
	if err != nil {
		return nil, err
	}
	if len(commands.set) == 0 {
		return &DefaultSafeContexts, nil
	}
	return commands, nil
}

func getUnSafeCommands() (*KubeCtlSafeMap, error) {
	commands, err := parseCommands(KubectlUnsafeCommands)
	if err != nil {
		return nil, err
	}
	return commands, nil
}

func getSafeCommands() (*KubeCtlSafeMap, error) {
	commands, err := parseCommands(KubectlSafeCommands)
	if err != nil {
		return nil, err
	}
	if len(commands.set) == 0 {
		return &DefaultSafeCommands, nil
	}
	return commands, nil
}

func isDryRun(cmd []string) bool {
	for _, c := range cmd {
		if strings.Contains(c, "--dry-run") {
			return true
		}
	}
	return false
}

func parseCommands(env string) (*KubeCtlSafeMap, error) {
	commands := os.Getenv(env)
	if commands == "" {
		return &EmptyCommands, nil
	}
	// commands can be a csv of strings or a file path
	if _, err := os.Stat(commands); errors.Is(err, fs.ErrNotExist) {
		cmds := strings.Split(commands, ",")
		if len(cmds) >= 1 {
			return getConfigFromSlice(cmds)
		}
		return nil, fmt.Errorf("%s was provided by %s but its neither a word or existing file", commands, env)
	}
	return getCommandsFromFile(commands)
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
