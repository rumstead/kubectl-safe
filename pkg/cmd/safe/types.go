package safe

const (
	// KubectlSafeCommands is a csv of commands or a path to a file containing a set of safe commands
	KubectlSafeCommands = "KUBECTL_SAFE_COMMANDS"
	// KubectlUnsafeCommands is a csv of commands or a path to a file containing a set of unsafe commands
	KubectlUnsafeCommands = "KUBECTL_UNSAFE_COMMANDS"
	// KubectlSafeContexts is a csv of contexts or a path to a file
	KubectlSafeContexts = "KUBECTL_SAFE_CONTEXTS"
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
	defaultContexts = map[string]Void{
		"minikube":        Empty,
		"docker-desktop":  Empty,
		"rancher-desktop": Empty,
	}
	DefaultSafeCommands = KubeCtlSafeMap{set: defaultCommands}
	DefaultSafeContexts = KubeCtlSafeMap{set: defaultContexts}
	Empty               = Void{}
	EmptyCommands       = KubeCtlSafeMap{}
)

type Void struct{}

type KubeCtlSafeMap struct {
	set map[string]Void
}

func (c *KubeCtlSafeMap) Contains(command string) bool {
	_, ok := c.set[command]
	return ok
}

func (c *KubeCtlSafeMap) Add(command string) {
	if command != "" {
		c.set[command] = Empty
	}
}
