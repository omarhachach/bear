package bear

// Command is the interface for a given command.
type Command interface {
	GetCallers() []string
	GetHandler() func(*Context)
}
