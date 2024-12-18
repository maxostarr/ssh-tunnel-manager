package app

import (
	"context"
	"fmt"
	"ssh-tunnel-manager/ssh_manager"
	"ssh-tunnel-manager/utils"
)

// var sshManager *ssh_manager.SshManager = ssh_manager.NewSshManager()

// App struct
type App struct {
	ctx           context.Context
	manager       *ssh_manager.SshManager
	eventsManager utils.EventManager
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
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.WithEvents()
}

func (a *App) WithEvents() {
	a.eventsManager = utils.NewEventManager(a.ctx)

	a.manager.PromptUser = a.eventsManager.Prompt
}

func (a *App) GetRemotes() []*ssh_manager.SshManagerRemoteData {
	return a.manager.GetRemotes()
}

func (a *App) AddRemote(name string, host string, port int, username string) (bool, error) {
	return a.manager.AddRemote(name, host, port, username)
}

func (a *App) UpdateRemote(id string, name string, host string, port int, username string) (bool, error) {
	fmt.Println("Updating remote with ID", id)
	return a.manager.UpdateRemote(id, name, host, port, username)
}

func (a *App) DeleteRemote(id string) (bool, error) {
	return a.manager.DeleteRemote(id)
}

func (a *App) GetRemote(id string) (*ssh_manager.SshManagerRemote, error) {
	return a.manager.GetRemote(id)
}

func (a *App) GetTunnels(remoteId string) []*ssh_manager.SshManagerTunnel {
	remote, err := a.manager.GetRemote(remoteId)
	if err != nil {
		return nil
	}
	return remote.Tunnels
}

func (a *App) AddTunnel(remoteId string, localPort int, remoteHost string, remotePort int) (bool, error) {
	remote, err := a.manager.GetRemote(remoteId)
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

// func (a *App) Connect(id string) (bool, error) {
// 	fmt.Println("Initiating connection to remote with ID" + id)
// 	remote, err := a.manager.GetRemote(id)
// 	if err != nil {
// 		return false, err
// 	}

// 	update := make(chan struct{})
// 	go remote.Connect(update)

// 	a.eventsManager.Emit("remotes-updated", nil)

// 	return true, nil
// }

func (a *App) Connect(id string) (bool, error) {
	fmt.Println("Initiating connection to remote with ID" + id)
	remote, err := a.manager.GetRemote(id)
	if err != nil {
		return false, err
	}

	update := make(chan struct{})
	done := make(chan struct{})

	go func() {
		remote.Connect(update)
		close(done)
	}()

	go func() {
		for {
			select {
			case <-update:
				a.eventsManager.Emit("remotes-updated", nil)
			case <-done:
				return
			}
		}
	}()

	<-done
	return true, nil
}

func (a *App) Disconnect(remoteName string) {
	remote, err := a.manager.GetRemote(remoteName)
	if err != nil {
		return
	}
	remote.Disconnect()
	a.eventsManager.Emit("remotes-updated", nil)
}

// func (a *App) TestPrompt() {
// 	promptOptions := utils.NewPromptOptions("Test prompt", "Confirm", "Cancel", []utils.PromptInput{
// 		{
// 			Label: "Name",
// 			Key:   "name",
// 			Type:  utils.PromptInputTypeText,
// 		},
// 		{
// 			Label: "Password",
// 			Key:   "password",
// 			Type:  utils.PromptInputTypePassword,
// 		},
// 	})

// 	response, err := a.manager.PromptUser(promptOptions)
// 	if err != nil {
// 		fmt.Println("Error getting prompt response", err)
// 		return
// 	}
// 	fmt.Println("Prompt response: ", response.Response["name"])
// }
