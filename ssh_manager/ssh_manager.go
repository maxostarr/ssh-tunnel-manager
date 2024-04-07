package ssh_manager

type SshManager struct {
	Remotes []*SshManagerRemote
	PromptUser func(prompt string) string
}


func (manager *SshManager) Initialize() {
	// Get all remotesData from the database
	remotesData, err := GetRemotes()
	if err != nil {
		panic(err)
	}

	for _, remoteData := range remotesData {
		remote := manager.NewSshManagerRemoteFromData(*remoteData)
		manager.Remotes = append(manager.Remotes, remote)
	}
}

func (manager *SshManager) AddRemote(name string, host string, port int, username string, password string) (bool, error) {
	remote := manager.NewSshManagerRemote(name, host, port, username, password)
	manager.Remotes = append(manager.Remotes, remote)
	return true, nil
}

func (manager *SshManager) GetRemote(name string) (*SshManagerRemote, error) {
	for _, remote := range manager.Remotes {
		if remote.Name == name {
			return remote, nil
		}
	}
	return nil, nil
}

func (manager *SshManager) RemoveRemote(name string) (bool, error) {
	for i, remote := range manager.Remotes {
		if remote.Name == name {
			manager.Remotes = append(manager.Remotes[:i], manager.Remotes[i+1:]...)
			return true, nil
		}
	}
	return false, nil
}

func (manager *SshManager) GetRemotes() []*SshManagerRemote {
	return manager.Remotes
}

