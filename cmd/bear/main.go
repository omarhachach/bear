package main

import (
	"github.com/omarhachach/bear"
	"os"
	"os/signal"
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

func main() {
	c := make(chan os.Signal, 1)
	
	signal.Notify(c, os.Interrupt, os.Kill)
	b := bear.New(&bear.Config{
		Debug: true,
		DiscordToken: "your-token-goes-here",
	}).RegisterModules(&Module{}).Start()

	<-c

	b.Close()
}
