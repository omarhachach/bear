package bear

import (
	"fmt"
	"os"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

// Version is the current version of Bear.
var Version = "0.3.1-alpha"

// Bear is the core bot.
type Bear struct {
	Commands map[string]Command
	Modules  map[string]Module
	Session  *discordgo.Session
	Log      *logrus.Logger
	Config   *Config
	Version  string
	mutex    *sync.Mutex
	logFile  *os.File
}

// New returns a new Bear.
func New(config *Config) *Bear {
	b := &Bear{
		Commands: make(map[string]Command),
		Modules:  make(map[string]Module),
		Log:      logrus.New(),
		mutex:    &sync.Mutex{},
		logFile:  nil,
		Version:  Version,
	}

	return b.UpdateConfig(config)
}

// UpdateConfig will update the configuration of the bot.
func (b *Bear) UpdateConfig(config *Config) *Bear {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.Config = config

	if b.Config.Log.Debug {
		b.Log.SetLevel(logrus.DebugLevel)
		b.Log.SetReportCaller(true)
	}

	if b.Config.Log.File == "" {
		b.Log.SetOutput(os.Stdout)
		return b
	}

	file, err := os.OpenFile(b.Config.Log.File, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		b.Log.WithError(err).Error("Error opening log file, not changing log location.")
	} else {
		b.logFile = file
		b.Log.SetOutput(file)
	}

	return b
}

// RegisterModules will register all of the modules passed to this function.
func (b *Bear) RegisterModules(modules ...Module) *Bear {
	for _, module := range modules {
		b.RegisterModule(module)
	}

	return b
}

// RegisterModule with register a module to the bot.
func (b *Bear) RegisterModule(module Module) *Bear {
	b.mutex.Lock()

	_, exists := b.Modules[module.GetName()]

	if exists {
		b.Log.Errorf("Failed to load module %s, module already exists.", module.GetName())
		b.mutex.Unlock()
		return b
	}

	b.Modules[module.GetName()] = module
	b.Log.Infof("Registered module %s.", module.GetName())
	b.mutex.Unlock()

	b.RegisterCommands(module.GetCommands())

	return b
}

// RegisterCommands will register an array of commands to the bot.
func (b *Bear) RegisterCommands(cmds []Command) *Bear {
	for _, cmd := range cmds {
		b.RegisterCommand(cmd)
	}

	return b
}

// RegisterCommand will check if a command exists, and add it to the bot.
func (b *Bear) RegisterCommand(cmd Command) *Bear {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for _, caller := range cmd.GetCallers() {
		_, exists := b.Commands[caller]

		if exists {
			b.Log.Errorf("Couldn't register caller %s, skipping.", caller)
			continue
		}

		b.Commands[caller] = cmd
		b.Log.Debugf("Registered caller %s.", caller)
	}

	return b
}

// AddHandler will add a handler to the Discord session
func (b *Bear) AddHandler(handler interface{}) *Bear {
	b.Log.Debug("Added event handler to the session")
	b.Session.AddHandler(handler)
	return b
}

// Start will open the Discord session, and initialize the bot.
func (b *Bear) Start() *Bear {
	session, err := initDiscordSession(b.Config.DiscordToken)
	if err != nil {
		b.Log.WithError(err).Fatal("Couldn't establish Discord session.")
		return b
	}

	b.Session = session
	b.AddHandler(onMessageCreate(b))

	b.Log.Info("Initing modules.")
	b.initModules()

	b.Log.Info("You have poked the bear!")
	return b
}

// CLose will close all the sessions properly.
func (b *Bear) Close() *Bear {
	b.Log.Info("The bear is sleepy.")

	err := b.Session.Close()
	if err != nil {
		b.Log.WithError(err).Error("Error closing Discord session.")
	}

	if b.logFile != nil {
		err = b.logFile.Close()
		if err != nil {
			b.Log.WithError(err).Error("Error closing log file.")
		}
	}

	b.Log.Debug("Closing all modules.")
	b.closeModules()

	b.Log.Info("The bear is now asleep.")
	return b
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

func (b *Bear) initModules() {
	for _, module := range b.Modules {
		module.Init(b)
		b.Log.Debugf("Initialized module: %s", module.GetName())
	}
}

func (b *Bear) closeModules() {
	for _, module := range b.Modules {
		module.Init(b)
		b.Log.Debugf("Closed module: %s", module.GetName())
	}
}
