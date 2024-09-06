package ssh_manager

import (
	"fmt"
	"ssh-tunnel-manager/utils"
	"strconv"

	"golang.org/x/crypto/ssh"
)

type SshManagerRemote struct {
	SshManagerRemoteData
	Auth    []ssh.AuthMethod    `json:"-"`
	Client  *ssh.Client         `json:"-"`
	Tunnels []*SshManagerTunnel `json:"tunnels"`
}

func (manager *SshManager) NewSshManagerRemote(name string, host string, port int, username string) *SshManagerRemote {
	remote := manager.NewSshManagerRemoteFromData(SshManagerRemoteData{
		Name:     name,
		Host:     host,
		Port:     port,
		Username: username,
	})

	return remote
}

func (manager *SshManager) NewSshManagerRemoteFromData(data SshManagerRemoteData) *SshManagerRemote {
	remote := &SshManagerRemote{
		SshManagerRemoteData: SshManagerRemoteData{
			Name:     data.Name,
			Host:     data.Host,
			Port:     data.Port,
			Username: data.Username,
			ID:       data.ID,
		},
	}

	remote.Auth = []ssh.AuthMethod{
		ssh.PasswordCallback(manager.promptPasswordChallenge),
		ssh.KeyboardInteractive(manager.promptKeyboardChallenge),
	}

	return remote
}

func (manager *SshManager) promptKeyboardChallenge(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
	fmt.Println(instruction)
	// answers = make([]string, len(questions))
	// for i := range questions {
	// 	answer := manager.PromptUser(questions[i])
	// 	if answer.Status == PromptResponseStatusCancelled {
	// 		return nil, fmt.Errorf("keyboard challenge cancelled")
	// 	}
	// 	answers[i] = answer.Response
	// }

	inputs := []utils.PromptInput{}

	for i, question := range questions {
		inputs = append(inputs, utils.PromptInput{
			Label: question,
			Key:   strconv.Itoa(i),
			Type:  utils.PromptInputTypeText,
		})
	}

	result, err := manager.PromptUser(utils.PromptOptions{
		Inputs: inputs,
	})

	if err != nil {
		return nil, err
	}

	resultArray := []string{}

	for _, response := range result.Response {
		resultArray = append(resultArray, response)
	}

	answers = resultArray

	return answers, nil
}

func (manager *SshManager) promptPasswordChallenge() (string, error) {
	// response := manager.PromptUser("Password: ")
	// fmt.Println("Password response: " + response.Response)
	// if response.Status == PromptResponseStatusCancelled {
	// 	return "", fmt.Errorf("password prompt cancelled")
	// }
	// return response.Response, nil

	result, err := manager.PromptUser(utils.PromptOptions{
		Inputs: []utils.PromptInput{
			{
				Label: "Password",
				Key:   "password",
				Type:  utils.PromptInputTypePassword,
			},
		},
	})

	if err != nil {
		return "", err
	}

	return result.Response["password"], nil
}

func (remote *SshManagerRemote) Initialize() {
	tunnelsData, err := GetTunnelsByRemote(remote.ID)
	remote.Tunnels = []*SshManagerTunnel{}
	if err != nil {
		panic(err)
	}
	for _, tunnelData := range tunnelsData {
		tunnel := NewSshManagerTunnelFromData(*tunnelData, remote)
		remote.Tunnels = append(remote.Tunnels, tunnel)
	}
}

func (remote *SshManagerRemote) Save() (string, error) {
	return InsertRemote(&remote.SshManagerRemoteData)
}

func (remote *SshManagerRemote) Update() error {
	return UpdateRemote(&remote.SshManagerRemoteData)
}

func (remote *SshManagerRemote) Connect() (bool, error) {
	config := &ssh.ClientConfig{
		User:            remote.Username,
		Auth:            remote.Auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	connectionString := remote.Host + ":" + strconv.Itoa(remote.Port)
	fmt.Println("Connecting to " + connectionString)
	client, err := ssh.Dial("tcp", connectionString, config)
	if err != nil {
		return false, err
	}
	fmt.Println("Connected to " + connectionString)
	remote.Client = client

	for _, tunnel := range remote.Tunnels {
		tunnel.Connect()
	}

	return true, nil
}

func (remote *SshManagerRemote) Disconnect() {
	for _, tunnel := range remote.Tunnels {
		tunnel.Disconnect()
	}
	remote.Client.Close()
}

func (remote *SshManagerRemote) AddTunnel(localPort int, remoteHost string, remotePort int) (bool, error) {
	tunnel := NewSshManagerTunnel(localPort, remoteHost, remotePort, remote)
	id, err := InsertTunnel(&tunnel.SshManagerTunnelData)
	if err != nil {
		return false, err
	}
	tunnel.ID = id
	remote.Tunnels = append(remote.Tunnels, tunnel)
	return true, nil
}

func (remote *SshManagerRemote) RemoveTunnel(localPort int) (bool, error) {
	for i, tunnel := range remote.Tunnels {
		if tunnel.LocalPort == localPort {
			tunnel.Disconnect()
			remote.Tunnels = append(remote.Tunnels[:i], remote.Tunnels[i+1:]...)
			return true, nil
		}
	}
	return false, nil
}
