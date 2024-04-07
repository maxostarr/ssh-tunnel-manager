package main

import (
	"context"
	"ssh-tunnel-manager/ssh_manager"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// var sshManager *ssh_manager.SshManager = ssh_manager.NewSshManager()

// App struct
type App struct {
	ctx      context.Context
	manager  *ssh_manager.SshManager
}

// NewApp creates a new App application struct
func NewApp() *App {
	ssh_manager.ConnectDB()
	ssh_manager.CreateTables()
	manager := &ssh_manager.SshManager{}
	manager.Initialize()
	return &App{
		manager: manager,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.manager.PromptUser = a.PromptUser
}

func (a *App) GetRemotes() []*ssh_manager.SshManagerRemote {
	return a.manager.GetRemotes()
}

func (a *App) AddRemote(name string, host string, port int, username string, password string) (bool, error) {
	return a.manager.AddRemote(name, host, port, username, password)
}

func (a *App) RemoveRemote(name string) (bool, error) {
	return a.manager.RemoveRemote(name)
}

func (a *App) GetRemote(name string) (*ssh_manager.SshManagerRemote, error) {
	return a.manager.GetRemote(name)
}

func (a *App) GetTunnels(remoteName string) []*ssh_manager.SshManagerTunnel {
	remote, err := a.manager.GetRemote(remoteName)
	if err != nil {
		return nil
	}
	return remote.Tunnels
}

func (a *App) AddTunnel(remoteName string, localPort int, remoteHost string, remotePort int) (bool, error) {
	remote, err := a.manager.GetRemote(remoteName)
	if err != nil {
		return false, err
	}
	return remote.AddTunnel(localPort, remoteHost, remotePort)
}

func (a *App) RemoveTunnel(remoteName string, localPort int) (bool, error) {
	remote, err := a.manager.GetRemote(remoteName)
	if err != nil {
		return false, err
	}
	return remote.RemoveTunnel(localPort)
}

func (a *App) Connect(remoteName string) (bool, error) {
	remote, err := a.manager.GetRemote(remoteName)
	if err != nil {
		return false, err
	}
	return remote.Connect()
}

func (a *App) Disconnect(remoteName string) {
	remote, err := a.manager.GetRemote(remoteName)
	if err != nil {
		return
	}
	remote.Disconnect()
}

func (a *App) PromptUser(prompt string) string {
	runtime.EventsEmit(a.ctx, "prompt", prompt)
	// Wait for the response
	responseChannel := make(chan string)
	runtime.EventsOnce(a.ctx, "prompt-response", func(data ...interface{}) {
		responseChannel <- data[0].(string)
	})

	return <-responseChannel
}

// func (a *App) Connect(localPort float64, remoteHost string, remotePort float64) error {

// 	// Create a new SSH Manager
// 	sshManager := ssh_manager.SshManagerRemote{
// 		Host:     remoteHost,
// 		Port:     22,
// 		Username: "rezo",
// 		Auth: []ssh.AuthMethod{
// 			ssh.Password("W1ndyC1tyMS"),
// 		},
// 	}

// 	// Connect to the remote host
// 	shouldReturn, returnValue := sshManager.Connect(int(localPort), remoteHost, int(remotePort))


// 	// shouldReturn, returnValue := Connect(localPort, remoteHost, remotePort)
// 	if shouldReturn {
// 		return returnValue
// 	}

// 	return nil
// }


