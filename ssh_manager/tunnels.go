package ssh_manager

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

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
	// Add connection deadline to prevent blocking forever
	tcpListener.SetDeadline(time.Now().Add(1 * time.Second))

	conn, err := tcpListener.Accept()
	if err != nil {
		select {
		case <-tunnel.stop:
			return false, nil
		default:
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				return true, nil // Continue loop on timeout
			}
			// Don't return error for closed listener during shutdown
			if strings.Contains(err.Error(), "use of closed network connection") {
				return false, nil
			}
			return true, fmt.Errorf("failed to accept connection: %v", err)
		}
	}

	select {
	case <-tunnel.stop:
		conn.Close()
		return false, nil
	default:
		tunnel.wg.Add(1)
		go func() {
			defer func() {
				conn.Close()
				tunnel.wg.Done()
			}()

			log.Printf("handling connection from %v", conn.RemoteAddr())
			handleConn(conn, tunnel.Remote.Client, tunnel.RemoteHost, tunnel.RemotePort, tunnel.stop, &tunnel.wg)
		}()
		return true, nil
	}
}

func (tunnel *SshManagerTunnel) Connect() (bool, error) {
	log.Printf("starting tunnel %s", tunnel.ID)

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", tunnel.LocalPort))
	if err != nil {
		return true, fmt.Errorf("failed to listen: %v", err)
	}

	tcpListener, ok := listener.(*net.TCPListener)
	if !ok {
		listener.Close()
		return true, fmt.Errorf("failed to cast to TCP listener")
	}

	defer func() {
		tcpListener.Close()
		tunnel.wg.Wait()
	}()

	for {
		select {
		case <-tunnel.stop:
			return false, nil
		default:
			cont, err := tunnel.handleIncomingConnection(tcpListener)
			if !cont || err != nil {
				return false, err
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
