package ssh_manager

import (
	"strconv"

	"golang.org/x/crypto/ssh"
)


type SshManagerRemote struct {
	SshManagerRemoteData	
	Auth    []ssh.AuthMethod
	Client  *ssh.Client
	Tunnels []*SshManagerTunnel
}

func NewSshManagerRemote(name string, host string, port int, username string, password string) *SshManagerRemote {
	remote := &SshManagerRemote{
		SshManagerRemoteData: SshManagerRemoteData{
			Name:     name,
			Host:     host,
			Port:     port,
			Username: username,
			Password: password,
		},
	}
	remote.Auth = []ssh.AuthMethod{ssh.Password(password)}
	return remote
}

func (remote *SshManagerRemote) Save() (string, error) {
	return InsertRemote(&remote.SshManagerRemoteData)
}

func (remote *SshManagerRemote) Connect() (bool, error) {
	config := &ssh.ClientConfig{
		User: remote.Username,
		Auth: remote.Auth,
	}
	client, err := ssh.Dial("tcp", remote.Host+":"+strconv.Itoa(remote.Port), config)
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

