package ssh_manager

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
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
		stop:   make(chan struct{}, 1),
		wg:     sync.WaitGroup{},
		Remote: remote,
	}

	return tunnel
}

func NewSshManagerTunnelFromData(data SshManagerTunnelData, remote *SshManagerRemote) *SshManagerTunnel {
	tunnel := &SshManagerTunnel{
		SshManagerTunnelData: data,
		Remote:               remote,
		stop:                 make(chan struct{}, 1),
		wg:                   sync.WaitGroup{},
	}
	return tunnel
}

func handleConn(localConn net.Conn, client *ssh.Client, remoteHost string, remotePort int, stop chan struct{}, wg *sync.WaitGroup) {
	defer localConn.Close()

	remoteConn, err := client.Dial("tcp", remoteHost+":"+strconv.Itoa(remotePort))
	if err != nil {
		select {
		case <-stop:
			log.Printf("Connection stopped while dialing")
			return
		default:
			log.Printf("Failed to connect to remote host: %v", err)
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
				log.Printf("Copy stopped via done channel")
				return
			default:
				log.Printf("Copy failed: %v", err)
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

	<-stop
	log.Printf("Stopping connection and cleanup")
	close(done)
}

func (tunnel *SshManagerTunnel) handleIncomingConnection(tcpListener *net.TCPListener) (bool, error) {
	// tcpListener.SetDeadline(time.Now().Add(100 * time.Millisecond))

	conn, err := tcpListener.Accept()
	if err != nil {
		log.Printf("handleIncomingConnection: Failed to accept connection: %v", err)
		select {
		case <-tunnel.stop:
			log.Printf("handleIncomingConnection: Stop signal received")
			return false, nil
		default:
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Printf("handleIncomingConnection: Timeout")

				log.Printf("handleIncomingConnection: Continue loop")
				return true, nil // Continue loop
			}
			if !strings.Contains(err.Error(), "use of closed network connection") {
				return true, fmt.Errorf("handleIncomingConnection: Failed to accept connection: %v", err)
			}
			return false, nil
		}
	}
	defer conn.Close()

	select {
	case <-tunnel.stop:
		conn.Close()
		return false, nil
	default:
		tunnel.wg.Add(1)
		go func() {
			defer tunnel.wg.Done()
			log.Printf("handleIncomingConnection: Handling connection")
			handleConn(conn, tunnel.Remote.Client, tunnel.RemoteHost, tunnel.RemotePort, tunnel.stop, &tunnel.wg)
			log.Printf("handleIncomingConnection: Connection handled")
		}()
		return true, nil
	}
}

func (tunnel *SshManagerTunnel) Connect() (bool, error) {
	fmt.Println("Connect: Starting tunnel - " + tunnel.ID)
	defer log.Output(2, "Connect: Connected tunnel - "+tunnel.ID)

	listener, err := net.Listen("tcp", "localhost:"+strconv.Itoa(int(tunnel.LocalPort)))
	if err != nil {
		return true, fmt.Errorf("Connect: Failed to listen: %v", err)
	}
	defer listener.Close()

	tcpListener, ok := listener.(*net.TCPListener)
	if !ok {
		return true, fmt.Errorf("Connect: Failed to cast to TCP listener")
	}

	// Add logging before select
	fmt.Println("Connect: Entering main loop")

	for {
		select {
		case <-tunnel.stop:
			fmt.Println("Connect: Stop signal received") // Add this
			fmt.Println("Connect: Stopping connection")
			tcpListener.Close()
			tunnel.wg.Wait()
			fmt.Println("Connect: Clean shutdown complete") // Add this
			return false, nil
		default:
			shouldContinue, err := tunnel.handleIncomingConnection(tcpListener)
			if err != nil {
				return true, err
			}
			if !shouldContinue {
				return false, nil
			}
		}

	}
}

func (tunnel *SshManagerTunnel) Disconnect() {
	fmt.Printf("Disconnect: Starting shutdown - %f, %p\n ", tunnel.LocalPort, tunnel.stop)
	close(tunnel.stop)
	fmt.Println("Disconnect: Stop channel closed")
	tunnel.wg.Wait()
	fmt.Println("Disconnect: All goroutines finished")
}
