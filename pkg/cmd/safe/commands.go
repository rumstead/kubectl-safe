package safe

import (
	"bufio"
	"fmt"
	"k8s.io/klog/v2"
	"os"
	"strings"
)

func IsVerbSafe(verb string) (bool, error) {
	commands, err := getSafeCommands()
	if err != nil {
		return false, err
	}
	return commands.Contains(verb), nil
}

func getSafeCommands() (*Commands, error) {
	if commands := os.Getenv(KubectlSafeCommands); commands != "" {
		return paresSafeCommands(commands)
	}
	return &DefaultSafeCommands, nil
}

func paresSafeCommands(commands string) (*Commands, error) {
	// commands can be a csv of strings or a file path
	_, err := os.Stat(commands)
	switch {
	case os.IsNotExist(err):
		cmds := strings.Split(commands, ",")
		if len(cmds) >= 1 {
			return getCommandsFromSlice(cmds)
		}
		return nil, fmt.Errorf("%s was provided by %s but its neither a word or existing file", commands, KubectlSafeCommands)
	case err == nil:
		return getCommandsFromFile(commands)
	default:
		return nil, err
	}
}

func getCommandsFromFile(commands string) (*Commands, error) {
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

func getCommandsFromSlice(cmds []string) (*Commands, error) {
	commands := NewCommands()
	for i := range cmds {
		commands.Add(cmds[i])
	}
	return &commands, nil
}

func NewCommands() Commands {
	return Commands{safeCmds: make(map[string]Void)}
}
