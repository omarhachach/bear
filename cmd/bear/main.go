package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/omarhachach/bear"
)

type Module struct {
	Name     string
	Desc     string
	Commands []bear.Command
	Version  string
}

func (m *Module) GetName() string {
	return m.Name
}

func (m *Module) GetDesc() string {
	return m.Desc
}

func (m *Module) GetCommands() []bear.Command {
	return m.Commands
}

func (m *Module) GetVersion() string {
	return m.Version
}

func (m *Module) Init(*bear.Bear) {
	return
}

func (m *Module) Close(*bear.Bear) {
	return
}

func main() {
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	b := bear.New(&bear.Config{
		Log: &bear.LogConfig{
			Debug: true,
			File:  "",
		},
		DiscordToken: "your-token-goes-here",
	}).RegisterModules(&Module{}).Start()

	<-c

	b.Close()
}
