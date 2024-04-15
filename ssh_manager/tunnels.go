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

type SshManagerTunnel struct {
	SshManagerTunnelData
	Remote *SshManagerRemote `json:"-"`
	stop   chan struct{}     `json:"-"`
	wg     sync.WaitGroup    `json:"-"`
}

func NewSshManagerTunnel(localPort int, remoteHost string, remotePort int, remote *SshManagerRemote) *SshManagerTunnel {
	tunnel := &SshManagerTunnel{
		SshManagerTunnelData: SshManagerTunnelData{
			LocalPort:  localPort,
			RemoteHost: remoteHost,
			RemotePort: remotePort,
			RemoteID:   remote.ID,
		},
		stop:   make(chan struct{}),
		wg:     sync.WaitGroup{},
		Remote: remote,
	}

	return tunnel
}

func NewSshManagerTunnelFromData(data SshManagerTunnelData, remote *SshManagerRemote) *SshManagerTunnel {
	tunnel := &SshManagerTunnel{
		SshManagerTunnelData: data,
		Remote:               remote,
		stop:                 make(chan struct{}),
		wg:                   sync.WaitGroup{},
	}
	return tunnel
}

func handleConn(localConn net.Conn, client *ssh.Client, remoteHost string, remotePort int, stop chan struct{}, wg *sync.WaitGroup) {
	defer localConn.Close()
	fmt.Println("Handling connection")

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

func (tunnel *SshManagerTunnel) Connect() (bool, error) {
	fmt.Println("Connecting tunnel - " + tunnel.ID)

	listener, err := net.Listen("tcp", "localhost:"+strconv.Itoa(int(tunnel.LocalPort)))
	if err != nil {
		return true, fmt.Errorf("Failed to listen: %v", err)
	}
	defer listener.Close()
	fmt.Println("Listening on localhost:" + strconv.Itoa(int(tunnel.LocalPort)))

	stop := make(chan struct{})
	var wg sync.WaitGroup

	for {
		conn, err := listener.Accept()
		fmt.Println("Accepted connection")
		if err != nil {
			select {
			case <-stop:
				return false, nil
			default:
				return true, fmt.Errorf("Failed to accept connection: %v", err)
			}
		}

		go handleConn(conn, tunnel.Remote.Client, tunnel.RemoteHost, tunnel.RemotePort, stop, &wg)
	}
}

func (tunnel *SshManagerTunnel) Disconnect() {
	close(tunnel.stop)
	tunnel.wg.Wait()
}
