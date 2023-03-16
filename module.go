package bear

// Module interface defines signature for a module.
type Module interface {
	Name() string
	Desc() string
	Commands() []Command
	Version() string
	Init(*Bear) error
	Close(*Bear) error
}
