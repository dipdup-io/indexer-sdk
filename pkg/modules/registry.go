package modules

var (
	globalRegistry *registry
)

type modules map[string]modulePorts

type registry struct {
	modules modules
}

func newRegistry() *registry {
	return &registry{
		modules: make(modules),
	}
}

// Register - call this function to be able to connect modules to each other
func Register(modules ...Module) error {
	if globalRegistry == nil {
		globalRegistry = newRegistry()
	}

	for i := range modules {
		ports, err := getModulePorts(modules[i])
		if err != nil {
			return err
		}
		globalRegistry.modules[modules[i].Name()] = ports
	}
	return nil
}
