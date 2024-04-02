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
	// remote.Auth = []ssh.AuthMethod{ssh.Password(password)}
	return remote
}

func NewSshManagerRemoteFromData(data SshManagerRemoteData) *SshManagerRemote {
	remote := &SshManagerRemote{
		SshManagerRemoteData: data,
	}
	remote.Auth = []ssh.AuthMethod{ssh.Password(data.Password)}
	return remote
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