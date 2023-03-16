package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/omarhachach/bear"
)

type Module struct{}

func (m *Module) Name() string {
	return "Test"
}

func (m *Module) Desc() string {
	return "Test"
}

func (m *Module) Commands() []bear.Command {
	return []bear.Command{}
}

func (m *Module) Version() string {
	return "v0.1.0"
}

func (m *Module) Init(*bear.Bear) error {
	return nil
}

func (m *Module) Close(*bear.Bear) error {
	return nil
}

func main() {
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	b := bear.New(&bear.Config{
		DiscordToken: "your-token-goes-here",
	})

	_ = b.RegisterModule(&Module{})

	_ = b.Start()

	<-c

	b.Close()
}
