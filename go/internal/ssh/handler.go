package ssh

import (
        "fmt"
        "log"
        "strings"
        "os"
        "net"
        "io"
	"io/ioutil"
	"golang.org/x/crypto/ssh"
	"github.com/atoonk/mysocketctl/go/internal/http"
	"github.com/dgrijalva/jwt-go"
)

const (
	mySocketSSHServer = "ssh.mysocket.io"
)

func SshConnect(socketID string, tunnelID string, port int, identityFile string) (error) {
	tunnel, err := http.GetTunnel(socketID, tunnelID)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	tokenStr, err := http.GetToken()
	if err != nil {
		return err
	}

	token, err := jwt.Parse(tokenStr, nil)
	if token == nil {
		log.Fatalf("error: %v", err)
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	tokenUserId := fmt.Sprintf("%v", claims["user_id"])
	userID := strings.ReplaceAll(tokenUserId, "-", "")

	sshConfig := &ssh.ClientConfig{
		User: userID,
		Auth: []ssh.AuthMethod{
			publicKeyFile(identityFile),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	serverConn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d",mySocketSSHServer, 22), sshConfig)
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
		log.Fatalf("Failed to create session: ", err)
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

		local, err := net.Dial("tcp", fmt.Sprintf("%s:%d","localhost",port))
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
