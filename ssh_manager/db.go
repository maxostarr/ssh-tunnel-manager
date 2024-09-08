package ssh_manager

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var connection *sql.DB

type SshManagerRemoteData struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
}

type SshManagerTunnelData struct {
	ID         string `json:"id"`
	LocalPort  int    `json:"local_port"`
	RemoteHost string `json:"remote_host"`
	RemotePort int    `json:"remote_port"`
	RemoteID   string `json:"remote_id"`
}

func ConnectDB() {
	dbName := "local.db"
	dbPath := filepath.Join(".", dbName)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("Error opening database:", err)
		os.Exit(1)
	}

	connection = db
}

func GetConnection() *sql.DB {
	return connection
}

func CloseConnection() {
	connection.Close()
}

func CreateTables() {
	_, err := connection.Exec(`CREATE TABLE IF NOT EXISTS remotes (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		host TEXT NOT NULL,
		port INTEGER NOT NULL,
		username TEXT NOT NULL
	);`)
	if err != nil {
		fmt.Println("Error creating remotes table:", err)
		os.Exit(1)
	}

	_, err = connection.Exec(`CREATE TABLE IF NOT EXISTS tunnels (
		id TEXT PRIMARY KEY,
		local_port INTEGER NOT NULL,
		remote_host TEXT NOT NULL,
		remote_port INTEGER NOT NULL,
		remote_id TEXT NOT NULL,
		FOREIGN KEY(remote_id) REFERENCES remotes(id) ON DELETE CASCADE
	);`)
	if err != nil {
		fmt.Println("Error creating tunnels table:", err)
		os.Exit(1)
	}
}

func GenerateUUID() string {
	return uuid.New().String()
}

func InsertRemote(remoteData *SshManagerRemoteData) (string, error) {
	id := GenerateUUID()
	fmt.Println("Inserting remote", remoteData)
	_, err := connection.Exec(`INSERT INTO remotes (id, name, host, port, username)
		VALUES (?, ?, ?, ?, ?);`, id, remoteData.Name, remoteData.Host, remoteData.Port, remoteData.Username)
	if err != nil {
		fmt.Println("Error inserting remote:", err)
		return "", err
	}
	fmt.Println("Inserted remote", remoteData)
	return id, nil
}

func UpdateRemote(remoteData *SshManagerRemoteData) error {
	_, err := connection.Exec(`UPDATE remotes SET name = ?, host = ?, port = ?, username = ? WHERE id = ?;`, remoteData.Name, remoteData.Host, remoteData.Port, remoteData.Username, remoteData.ID)
	if err != nil {
		fmt.Println("Error updating remote:", err)
		return err
	}
	return nil
}

func InsertTunnel(tunnelData *SshManagerTunnelData) (string, error) {
	id := GenerateUUID()
	_, err := connection.Exec(`INSERT INTO tunnels (id, local_port, remote_host, remote_port, remote_id)
		VALUES (?, ?, ?, ?, ?);`, id, tunnelData.LocalPort, tunnelData.RemoteHost, tunnelData.RemotePort, tunnelData.RemoteID)
	if err != nil {
		return "", err
	}
	return id, nil
}

func GetRemote(id string) (SshManagerRemoteData, error) {
	var remoteData SshManagerRemoteData
	err := connection.QueryRow(`SELECT name, host, port, username FROM remotes WHERE id = ?;`, id).Scan(&remoteData.Name, &remoteData.Host, &remoteData.Port, &remoteData.Username)
	if err != nil {
		return SshManagerRemoteData{}, err
	}
	remoteData.ID = id
	return remoteData, nil
}

// func GetTunnel(id string) (*SshManagerTunnelData, error) {
// 	var tunnelData SshManagerTunnelData
// 	err := connection.QueryRow(`SELECT local_port, remote_host, remote_port, remote_id FROM tunnels WHERE id = ?;`, id).Scan(&tunnelData.LocalPort, &tunnelData.RemoteHost, &tunnelData.RemotePort, &tunnelData.RemoteID)
// 	if err != nil {
// 		return &SshManagerTunnelData{}, err
// 	}
// 	tunnelData.ID = id
// 	return &tunnelData, nil
// }

func GetRemotes() ([]*SshManagerRemoteData, error) {
	rows, err := connection.Query(`SELECT id, name, host, port, username FROM remotes;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var remotes []*SshManagerRemoteData
	for rows.Next() {
		var remoteData SshManagerRemoteData
		err = rows.Scan(&remoteData.ID, &remoteData.Name, &remoteData.Host, &remoteData.Port, &remoteData.Username)
		if err != nil {
			return nil, err
		}
		fmt.Println(remoteData)
		remotes = append(remotes, &remoteData)
	}
	return remotes, nil
}

func GetTunnels() ([]*SshManagerTunnelData, error) {
	rows, err := connection.Query(`SELECT id, local_port, remote_host, remote_port, remote_id FROM tunnels;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tunnels []*SshManagerTunnelData
	for rows.Next() {
		var tunnelData SshManagerTunnelData
		err = rows.Scan(&tunnelData.ID, &tunnelData.LocalPort, &tunnelData.RemoteHost, &tunnelData.RemotePort, &tunnelData.RemoteID)
		if err != nil {
			return nil, err
		}
		tunnels = append(tunnels, &tunnelData)
	}
	return tunnels, nil
}

func GetTunnelsByRemote(remoteID string) ([]*SshManagerTunnelData, error) {
	rows, err := connection.Query(`SELECT id, local_port, remote_host, remote_port FROM tunnels WHERE remote_id = ?;`, remoteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tunnels []*SshManagerTunnelData
	for rows.Next() {
		var tunnelData SshManagerTunnelData
		err = rows.Scan(&tunnelData.ID, &tunnelData.LocalPort, &tunnelData.RemoteHost, &tunnelData.RemotePort)
		if err != nil {
			return nil, err
		}
		tunnelData.RemoteID = remoteID
		tunnels = append(tunnels, &tunnelData)
	}
	return tunnels, nil
}

func DeleteRemote(id string) error {
	_, err := connection.Exec(`DELETE FROM remotes WHERE id = ?;`, id)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTunnel(id string) error {
	_, err := connection.Exec(`DELETE FROM tunnels WHERE id = ?;`, id)
	if err != nil {
		return err
	}
	return nil
}
