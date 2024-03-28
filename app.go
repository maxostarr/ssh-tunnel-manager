package main

import (
	"context"
	"fmt"
	"ssh-tunnel-manager/ssh_manager"

	"golang.org/x/crypto/ssh"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) Connect(localPort float64, remoteHost string, remotePort float64) error {

	// Create a new SSH Manager
	sshManager := ssh_manager.SshManagerRemote{
		Host:     remoteHost,
		Port:     22,
		Username: "rezo",
		Auth: []ssh.AuthMethod{
			ssh.Password("W1ndyC1tyMS"),
		},
	}

	// Connect to the remote host
	shouldReturn, returnValue := sshManager.Connect(int(localPort), remoteHost, int(remotePort))


	// shouldReturn, returnValue := Connect(localPort, remoteHost, remotePort)
	if shouldReturn {
		return returnValue
	}

	return nil
}


