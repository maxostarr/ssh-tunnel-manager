package main

// func (a *App) PromptUser(prompt string) ssh_manager.PromptResponse {
// 	fmt.Println("Prompting user with: " + prompt)
// 	// runtime.EventsEmit(a.ctx, "prompt", prompt)
// 	// // Wait for the response
// 	// responseChannel := make(chan ssh_manager.PromptResponse)
// 	// runtime.EventsOnce(a.ctx, "prompt-response", func(data ...interface{}) {
// 	// 	promptResponse := ssh_manager.PromptResponse{
// 	// 		Status:   ssh_manager.PromptResponseStatus(data[0].(string)),
// 	// 		Response: data[1].(string),
// 	// 	}

// 	// 	responseChannel <- promptResponse
// 	// })

// 	// return <-responseChannel
// }

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
