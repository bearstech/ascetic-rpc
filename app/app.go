package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bearstech/ascetic-rpc/server"
)

type Server interface {
	Config() error
	GetServers() *server.ServerUsers
}

type Application struct {
	server Server
}

func NewApplication(server Server) *Application {
	return &Application{
		server: server,
	}
}

func (a *Application) Stop() {
	a.server.GetServers().Stop()
}

func (a *Application) Serve() {
	a.server.GetServers().Serve()
}

func (a *Application) Wait() {
	a.server.GetServers().Wait()
}

func (a *Application) Start() error {
	err := a.server.Config()
	if err != nil {
		return err
	}

	a.Serve()

	cc := make(chan os.Signal, 1)
	signal.Notify(cc, os.Interrupt, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGTERM)
	go func() {
		for {
			s := <-cc
			//log.Info("Signal : ", s)
			switch s {
			case os.Interrupt:
				a.Stop()
			case syscall.SIGTERM:
				a.Stop()
			case syscall.SIGHUP:
				err := a.server.Config()
				if err != nil {
					panic(err)
				}
				a.Serve()
			}
		}
	}()

	// block
	a.Wait()
	return nil
}
