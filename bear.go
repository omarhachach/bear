package bear

import (
	"errors"
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
)

// The errors returned from Bear commands.
var (
	ErrModuleAlreadyExists  = errors.New("module does not have a unique name")
	ErrCommandAlreadyExists = errors.New("command does not have a unique caller")
)

// Version is the current version of Bear.
var Version = "0.4.0-alpha"

// Bear is the core bot.
type Bear struct {
	Commands map[string]Command
	Modules  map[string]Module
	Session  *discordgo.Session
	Config   *Config
	Version  string
	mutex    *sync.Mutex
}

// New returns a new Bear.
func New(config *Config) *Bear {
	b := &Bear{
		Commands: make(map[string]Command),
		Modules:  make(map[string]Module),
		mutex:    &sync.Mutex{},
		Version:  Version,
	}

	return b.UpdateConfig(config)
}

// UpdateConfig will update the configuration of the bot.
func (b *Bear) UpdateConfig(config *Config) *Bear {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.Config = config

	return b
}

// RegisterModules will register all of the modules passed to this function.
func (b *Bear) RegisterModules(modules ...Module) error {
	var errs []error
	for _, module := range modules {
		err := b.RegisterModule(module)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// RegisterModule with register a module to the bot.
func (b *Bear) RegisterModule(module Module) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	_, exists := b.Modules[module.Name()]

	if exists {
		return fmt.Errorf("module with name %s alreayd exists: %w", module.Name(), ErrModuleAlreadyExists)
	}

	b.Modules[module.Name()] = module

	return b.RegisterCommands(module.Commands())
}

// RegisterCommands will register an array of commands to the bot.
func (b *Bear) RegisterCommands(cmds []Command) error {
	var errs []error
	for _, cmd := range cmds {
		err := b.RegisterCommand(cmd)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// RegisterCommand will check if a command exists, and add it to the bot.
func (b *Bear) RegisterCommand(cmd Command) error {
	var errs []error

	b.mutex.Lock()
	defer b.mutex.Unlock()

	for _, caller := range cmd.Callers() {
		_, exists := b.Commands[caller]

		if exists {
			errs = append(errs, fmt.Errorf("could not register command %s: %w", caller, ErrCommandAlreadyExists))
			continue
		}

		b.Commands[caller] = cmd
	}

	return errors.Join(errs...)
}

// AddHandler will add a handler to the Discord session
func (b *Bear) AddHandler(handler interface{}) {
	b.Session.AddHandler(handler)
}

// Start will open the Discord session, and initialize the bot.
func (b *Bear) Start() error {
	session, err := initDiscordSession(b.Config.DiscordToken)
	if err != nil {
		return fmt.Errorf("could not start discord session: %w", err)
	}

	b.Session = session
	b.AddHandler(onMessageCreate(b))

	return b.initModules()
}

// Close will close all the sessions properly.
func (b *Bear) Close() error {
	err := b.Session.Close()
	if err != nil {
		return fmt.Errorf("could not close discord session: %w", err)
	}

	return b.closeModules()
}

func initDiscordSession(token string) (*discordgo.Session, error) {
	session, err := discordgo.New(fmt.Sprintf("Bot %s", token))
	if err != nil {
		return nil, err
	}

	err = session.Open()
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (b *Bear) initModules() error {
	var errs []error
	for _, module := range b.Modules {
		errs = append(errs, module.Init(b))
	}
	return errors.Join(errs...)
}

func (b *Bear) closeModules() error {
	var errs []error
	for _, module := range b.Modules {
		errs = append(errs, module.Close(b))
	}
	return errors.Join(errs...)
}
