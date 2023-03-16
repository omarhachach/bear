package bear

// Command is the interface for a given command.
type Command interface {
	Callers() []string
	Handler(*Context)
}
