package bear

// Module interface defines signature for a module.
type Module interface {
	GetName() string
	GetDesc() string
	GetCommands() []Command
	GetVersion() string
	Init(*Bear)
}
