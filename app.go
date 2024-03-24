package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"

	"golang.org/x/crypto/ssh"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}


// func (a *App) Connect(localPort float64, remotePort float64, remoteHost string) {
	
// 	fmt.Printf("Connecting to macmini.local:22: %f:%s:%f", localPort, remoteHost, remotePort)

// 	// publicKeyFile, err := os.ReadFile("~/.ssh/id_rsa.pub")

// 	// if err != nil {
// 	// 	fmt.Println("Failed to read public key: ", err)
// 	// 	return
// 	// }

// 	// signer, err := ssh.ParsePrivateKey(publicKeyFile)

// 	// if err != nil {
// 	// 	fmt.Println("Failed to parse private key: ", err)
// 	// 	return
// 	// }

// 	config := &ssh.ClientConfig{
// 		User: "rezo",
// 		Auth: []ssh.AuthMethod{
// 			// ssh.PublicKeys(signer),
// 			ssh.Password("W1ndyC1tyMS"),
// 		},
// 		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
// 	}

// 	client, err := ssh.Dial("tcp", "macmini.local:22", config)

// 	if err != nil {
// 		fmt.Println("Failed to dial: ", err)
// 		return
// 	}


// 	listener, err := net.Listen("tcp", "localhost:"+strconv.Itoa(int(localPort)))
//   if err != nil {
//     // log.Fatal(err)
// 		fmt.Println("Failed to listen: ", err)
//   }
//   defer listener.Close()

// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		go handleConn(conn, client, remoteHost, remotePort)
// 	}
// }

// func handleConn(conn net.Conn, client *ssh.Client, remoteHost string, remotePort float64) {
// 	remote, err := client.Dial("tcp", remoteHost+":"+strconv.Itoa(int(remotePort)))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("Connected to macmini.local:9443")

// 	var wg sync.WaitGroup
// 	wg.Add(2)

// 	go func() {
// 		defer wg.Done()
// 		_, err := io.Copy(remote, conn)
// 		if err != nil {
// 			if err == io.EOF {
// 				fmt.Println("Connection closed by remote")
// 			} else {
// 				fmt.Println("remote => local Error: ", err)
// 			}
// 		}
// 	}()

// 	go func() {
// 		defer wg.Done()
// 		_, err := io.Copy(conn, remote)
// 		if err != nil {
// 			if err == io.EOF {
// 				fmt.Println("Connection closed by local")
// 			} else {
// 				fmt.Println("local => remote Error: ", err)
// 			}
// 		}
// 	}()

// 	wg.Wait()

// 	fmt.Println("Closing connection")

// 	conn.Close()
// 	remote.Close()
// }

func (a *App) Connect(localPort float64, remoteHost string, remotePort float64) error {
	config := &ssh.ClientConfig{
		User: "rezo",
		Auth: []ssh.AuthMethod{
			ssh.Password("W1ndyC1tyMS"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", "macmini.local:22", config)
	if err != nil {
		return fmt.Errorf("Failed to dial: %v", err)
	}
	defer client.Close()

	listener, err := net.Listen("tcp", "localhost:"+strconv.Itoa(int(localPort)))
	if err != nil {
		return fmt.Errorf("Failed to listen: %v", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("Failed to accept connection: %v", err)
		}

		go handleConn(conn, client, remoteHost, remotePort)
	}
}

func handleConn(localConn net.Conn, client *ssh.Client, remoteHost string, remotePort float64) {
    defer localConn.Close()

    remoteConn, err := client.Dial("tcp", remoteHost+":"+strconv.Itoa(int(remotePort)))
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

    go copyConn(localConn, remoteConn)
    copyConn(remoteConn, localConn)
}