package ssh

import (
	"fmt"
	"github.com/atoonk/mysocketctl/go/internal/http"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
)

const (
	mySocketSSHServer = "ssh.mysocket.io"
)

func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

func SshConnect(userID string, socketID string, tunnelID string, port int, identityFile string) error {
	tunnel, err := http.GetTunnel(socketID, tunnelID)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	sshConfig := &ssh.ClientConfig{
		User: userID,
		Auth: []ssh.AuthMethod{
			SSHAgent(),
			publicKeyFile(identityFile),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	serverConn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", mySocketSSHServer, 22), sshConfig)
	if err != nil {
		log.Fatalf("Dial INTO remote server error: %s", err)
	}

	listener, err := serverConn.Listen("tcp", fmt.Sprintf("localhost:%d", tunnel.LocalPort))
	if err != nil {
		log.Fatalf("Listen open port ON remote server on port %d error: %s", tunnel.LocalPort, err)
	}
	defer listener.Close()

	log.Printf("ssh tunnel started to localhost:%d", tunnel.LocalPort)

	session, err := serverConn.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	defer session.Close()

	log.Printf("ssh session started")

	session.Stdout = os.Stdout

	if err := session.Shell(); err != nil {
		log.Fatal(err)
	}

	for {
		client, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		local, err := net.Dial("tcp", fmt.Sprintf("%s:%d", "localhost", port))
		if err != nil {
			log.Printf("Dial INTO local service error: %s", err)
			continue
		}
		handleClient(client, local)
	}
}

func handleClient(client net.Conn, remote net.Conn) {
	defer client.Close()
	chDone := make(chan bool)

	// Start remote -> local data transfer
	go func() {
		_, err := io.Copy(client, remote)
		if err != nil {
			log.Println(fmt.Sprintf("error while copy remote->local: %s", err))
		}
		chDone <- true
	}()

	// Start local -> remote data transfer
	go func() {
		_, err := io.Copy(remote, client)
		if err != nil {
			log.Println(fmt.Sprintf("error while copy local->remote: %s", err))
		}
		chDone <- true
	}()

	<-chDone
}

func publicKeyFile(file string) ssh.AuthMethod {
	if file == "" {
		return nil
	}

	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Cannot read SSH public key file %s", file))
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Cannot parse SSH public key file %s", file))
		return nil
	}
	return ssh.PublicKeys(key)
}
