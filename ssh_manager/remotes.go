package ssh_manager

import (
	"fmt"
	"strconv"

	"golang.org/x/crypto/ssh"
)


type SshManagerRemote struct {
	SshManagerRemoteData	
	Auth    []ssh.AuthMethod
	Client  *ssh.Client
	Tunnels []*SshManagerTunnel
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
			ID:				data.ID,
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
	answers = make([]string, len(questions))
	for i := range questions {
		answers[i] = manager.PromptUser(questions[i])
	}
	
	return answers, nil
}

func (manager *SshManager) promptPasswordChallenge() (string, error) {
	fmt.Println("Password: ")
	return manager.PromptUser("Password: "), nil
}

func (remote *SshManagerRemote) Initialize() {
	tunnelsData, err := GetTunnelsByRemote(remote.ID)
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

func (remote *SshManagerRemote) Connect() (bool, error) {
	config := &ssh.ClientConfig{
		User: remote.Username,
		Auth: remote.Auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	connectionString := remote.Host + ":" + strconv.Itoa(remote.Port)
	fmt.Println("Connecting to " + connectionString)
	client, err := ssh.Dial("tcp", connectionString, config)
	if err != nil {
		return false, err
	}
	remote.Client = client
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