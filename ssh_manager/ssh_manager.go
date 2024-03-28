package ssh_manager

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"

	"golang.org/x/crypto/ssh"
)

type SshManager struct {
	Remotes []SshManagerRemote
}

type SshManagerRemote struct {
	Host    string
	Port    int
	Username string
	Auth    []ssh.AuthMethod
	Client  *ssh.Client
	Tunnels []SshManagerTunnel
}

type SshManagerTunnel struct {
	LocalPort int
	RemoteHost string
	RemotePort int
	stop      chan struct{}
	wg        sync.WaitGroup
}

func (tunnel *SshManagerTunnel) Disconnect() {
	close(tunnel.stop)
	tunnel.wg.Wait()
}

func handleConn(localConn net.Conn, client *ssh.Client, remoteHost string, remotePort int, stop chan struct{}, wg *sync.WaitGroup) {
	defer localConn.Close()

	remoteConn, err := client.Dial("tcp", remoteHost+":"+strconv.Itoa(remotePort))
	if err != nil {
		log.Fatalf("Failed to connect to remote host: %v", err)
	}
	defer remoteConn.Close()

	copyConn := func(writer, reader net.Conn) {
		defer writer.Close()
		_, err := io.Copy(writer, reader)
		if err != nil {
			log.Fatalf("io.Copy failed: %v", err)
		}
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		copyConn(localConn, remoteConn)
	}()

	go func() {
		defer wg.Done()
		copyConn(remoteConn, localConn)
	}()

	<-stop
}

func (tunnel *SshManagerRemote) Connect(localPort int, remoteHost string, remotePort int) (bool, error) {
	config := &ssh.ClientConfig{
		User: tunnel.Username,
		Auth: tunnel.Auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", tunnel.Host+":"+strconv.Itoa(tunnel.Port), config)
	if err != nil {
		return true, fmt.Errorf("Failed to dial: %v", err)
	}
	defer client.Close()

	listener, err := net.Listen("tcp", "localhost:"+strconv.Itoa(int(localPort)))
	if err != nil {
		return true, fmt.Errorf("Failed to listen: %v", err)
	}
	defer listener.Close()

	stop := make(chan struct{})
	var wg sync.WaitGroup

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-stop:
				return false, nil
			default:
				return true, fmt.Errorf("Failed to accept connection: %v", err)
			}
		}

		go handleConn(conn, client, remoteHost, remotePort, stop, &wg)
	}
}