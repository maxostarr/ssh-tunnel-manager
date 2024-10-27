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
	log.Output(2, "Handling connection")

	remoteConn, err := client.Dial("tcp", remoteHost+":"+strconv.Itoa(remotePort))
	if err != nil {
		select {
		case <-stop:
			fmt.Println("Stopping connection")
			return
		default:
			fmt.Errorf("Failed to connect to remote host: %v", err)
		}
		return
	}
	defer remoteConn.Close()

	copyConn := func(writer, reader net.Conn, done chan struct{}) {
		defer writer.Close()
		_, err := io.Copy(writer, reader)
		if err != nil {
			select {
			case <-done:
				return
			default:
				fmt.Errorf("io.Copy failed: %v", err)
			}
		}
	}

	done := make(chan struct{})
	wg.Add(2)
	go func() {
		defer wg.Done()
		copyConn(localConn, remoteConn, done)
	}()

	go func() {
		defer wg.Done()
		copyConn(remoteConn, localConn, done)
	}()

	log.Output(2, "Waiting for stop")
	<-stop
	fmt.Println("Stopping connection")
	close(done)
}

func (tunnel *SshManagerTunnel) Connect() (bool, error) {
	fmt.Println("Connecting tunnel - " + tunnel.ID)

	listener, err := net.Listen("tcp", "localhost:"+strconv.Itoa(int(tunnel.LocalPort)))
	if err != nil {
		return true, fmt.Errorf("Failed to listen: %v", err)
	}
	defer listener.Close()
	fmt.Println("Listening on localhost:" + strconv.Itoa(int(tunnel.LocalPort)))

	for {

		log.Output(2, "Waiting for connection")
		select {
		case <-tunnel.stop:
			fmt.Println("Stopping connection")
			return false, nil
		default:
		}

		conn, err := listener.Accept()
		log.Output(2, "Accepting connection")
		if err != nil {
			select {
			case <-tunnel.stop:
				log.Output(2, "Stopping connection")
				return false, nil
			default:
				fmt.Errorf("Failed to accept connection: %v", err)
			}
		}

		go handleConn(conn, tunnel.Remote.Client, tunnel.RemoteHost, tunnel.RemotePort, tunnel.stop, &tunnel.wg)

	}

	// for {
	// 	select {
	// 	case <-tunnel.stop:
	// 		log.Output(2, "Stopping connection")
	// 		return false, nil
	// 	default:
	// 		if tcpListener, ok := listener.(*net.TCPListener); ok {
	// 			tcpListener.SetDeadline(time.Now().Add(1 * time.Second)) // Set a deadline to periodically check the stop channel
	// 		}
	// 		conn, err := listener.Accept()
	// 		if err != nil {
	// 			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
	// 				continue // Timeout error, continue to check the stop channel
	// 			}
	// 			return true, fmt.Errorf("Failed to accept connection: %v", err)
	// 		}

	// 		go handleConn(conn, tunnel.Remote.Client, tunnel.RemoteHost, tunnel.RemotePort, tunnel.stop, &tunnel.wg)
	// 	}
	// }
}

func (tunnel *SshManagerTunnel) Disconnect() {
	close(tunnel.stop)
	tunnel.wg.Wait()
}
