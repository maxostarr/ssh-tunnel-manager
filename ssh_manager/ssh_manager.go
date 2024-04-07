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

func (manager *SshManager) AddRemote(name string, host string, port int, username string) (bool, error) {
	remote := manager.NewSshManagerRemote(name, host, port, username)
	id, err := InsertRemote(&remote.SshManagerRemoteData) 
	if err != nil {
		return false, err
	}
	remote.ID = id
	manager.Remotes = append(manager.Remotes, remote)
	return true, nil
}

func (manager *SshManager) GetRemote(id string) (*SshManagerRemote, error) {
	for _, remote := range manager.Remotes {
		if remote.ID == id {
			return remote, nil
		}
	}
	return nil, nil
}

func (manager *SshManager) RemoveRemote(id string) (bool, error) {
	for i, remote := range manager.Remotes {
		if remote.ID == id {
			manager.Remotes = append(manager.Remotes[:i], manager.Remotes[i+1:]...)
			DeleteRemote(id)
			return true, nil
		}
	}
	return false, nil
}

func (manager *SshManager) GetRemotes() []*SshManagerRemoteData {
	var remotesData []*SshManagerRemoteData
	for _, remote := range manager.Remotes {
		remotesData = append(remotesData, &remote.SshManagerRemoteData)
	}
	return remotesData
}

