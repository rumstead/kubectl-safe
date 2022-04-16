package safe

const (
	// KubectlSafeCommands is a csv of commands or a path to a file
	KubectlSafeCommands = "KUBECTL_SAFE_COMMANDS"
)

var (
	defaultCommands = map[string]Void{
		"get":           Empty,
		"describe":      Empty,
		"explain":       Empty,
		"cluster-info":  Empty,
		"top":           Empty,
		"config":        Empty,
		"logs":          Empty,
		"cp":            Empty,
		"diff":          Empty,
		"completion":    Empty,
		"alpha":         Empty,
		"api-resources": Empty,
		"api-versions":  Empty,
		"plugin":        Empty,
		"version":       Empty,
	}
	DefaultSafeCommands = Commands{safeCmds: defaultCommands}
	Empty               = Void{}
)

type Void struct{}

type Commands struct {
	safeCmds map[string]Void
}

func (c *Commands) Contains(command string) bool {
	_, ok := c.safeCmds[command]
	return ok
}

func (c *Commands) Add(command string) {
	if command != "" {
		c.safeCmds[command] = Empty
	}
}
