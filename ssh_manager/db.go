package ssh_manager

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"github.com/tursodatabase/go-libsql"
)

var connection *sql.DB

type SshManagerRemoteData struct {
	ID       string
	Name     string
	Host     string
	Port     int
	Username string
	Password string
}

type SshManagerTunnelData struct {
	ID         string
	LocalPort  int
	RemoteHost string
	RemotePort int
	RemoteID   string
}

func Connect() {
    dbName := "local.db"
    primaryUrl := "libsql://ssh-tunnel-manager-maxostarr.turso.io"
    authToken := "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3MTIwMTUwNDYsImlkIjoiOWZlY2MwZmUtN2RjOS00YjJhLTliN2EtYzVkZjlkNDk3ZjIwIn0.AoK1G2XL87t8EwuUBZ_Nws4tK6Tm4nw1S8hSUTYp4NkJuSrIJjvwGUM_2oTfj5fvSQsICcyCcf6SeMymXJWNDA"

    dir, err := os.MkdirTemp("", "libsql-*")
    if err != nil {
        fmt.Println("Error creating temporary directory:", err)
        os.Exit(1)
    }
    defer os.RemoveAll(dir)

    dbPath := filepath.Join(dir, dbName)

    connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, primaryUrl,
        libsql.WithAuthToken(authToken),
    )
    if err != nil {
        fmt.Println("Error creating connector:", err)
        os.Exit(1)
    }
    defer connector.Close()

    connection = sql.OpenDB(connector)
    // defer connection.Close()
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
		username TEXT NOT NULL,
		password TEXT NOT NULL
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
	_, err := connection.Exec(`INSERT INTO remotes (id, name, host, port, username, password)
		VALUES (?, ?, ?, ?, ?, ?);`, id, remoteData.Name, remoteData.Host, remoteData.Port, remoteData.Username, remoteData.Password)
	if err != nil {
		return "", err
	}
	return id, nil
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
	err := connection.QueryRow(`SELECT name, host, port, username, password FROM remotes WHERE id = ?;`, id).Scan(&remoteData.Name, &remoteData.Host, &remoteData.Port, &remoteData.Username, &remoteData.Password)
	if err != nil {
		return SshManagerRemoteData{}, err
	}
	remoteData.ID = id
	return remoteData, nil
}

func GetTunnel(id string) (*SshManagerTunnelData, error) {
	var tunnelData SshManagerTunnelData
	err := connection.QueryRow(`SELECT local_port, remote_host, remote_port, remote_id FROM tunnels WHERE id = ?;`, id).Scan(&tunnelData.LocalPort, &tunnelData.RemoteHost, &tunnelData.RemotePort, &tunnelData.RemoteID)
	if err != nil {
		return &SshManagerTunnelData{}, err
	}
	tunnelData.ID = id
	return &tunnelData, nil
}

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

