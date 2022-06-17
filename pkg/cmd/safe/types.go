package safe

const (
	// KubectlSafeCommands is a csv of commands or a path to a file containing a set of safe commands
	KubectlSafeCommands = "KUBECTL_SAFE_COMMANDS"
	// KubectlUnsafeCommands is a csv of commands or a path to a file containing a set of unsafe commands
	KubectlUnsafeCommands = "KUBECTL_UNSAFE_COMMANDS"
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
		"krew":          Empty,
	}
	DefaultSafeCommands = Commands{cmds: defaultCommands}
	Empty               = Void{}
	EmptyCommands       = Commands{}
)

type Void struct{}

type Commands struct {
	cmds map[string]Void
}

func (c *Commands) Contains(command string) bool {
	_, ok := c.cmds[command]
	return ok
}

func (c *Commands) Add(command string) {
	if command != "" {
		c.cmds[command] = Empty
	}
}
