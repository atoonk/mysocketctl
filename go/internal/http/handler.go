package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"
	"errors"
	"fmt"
	"os"
	h "net/http"
)

const (
	mysocketurl = "https://api.mysocket.io"
)

var (
        tokenfile   = fmt.Sprintf("%s/.mysocketio_token",os.Getenv("HOME"))
)


// Client . . .
type Client interface {
	CreateSocket(name string) error
}

type client struct {
	token string
}

var _ Client = &client{}

// NewClientWithToken ...
func NewClientWithToken(token string) Client {
	return &client{token: token}
}

// Login ...
func Login(email, password string) (Client, error) {
	c := &client{}
	form := loginForm{Email: email, Password: password}
	buf, err := json.Marshal(form)
	if err != nil {
		return nil, err
	}

	requestReader := bytes.NewReader(buf)

	resp, err := h.Post(mysocketurl+"/login", "application/json", requestReader)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return nil, errors.New("Login failed")
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("failed to login")
	}

	res := tokenForm{}
	json.NewDecoder(resp.Body).Decode(&res)

	c.token = res.Token

	f, err := os.Create(tokenfile)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	_, err2 := f.WriteString(fmt.Sprintf("%s\n", c.token))
	if err2 != nil {
		return nil, err2
	}

	return c, nil
}

// Register ...
func Register(name, email, password, sshkey string) error {
	form := registerForm{Name: name, Email: email, Password: password, Sshkey: sshkey}
	buf, err := json.Marshal(form)
	if err != nil {
		return err
	}
	requestReader := bytes.NewReader(buf)
	resp, err := h.Post(mysocketurl+"/user", "application/json", requestReader)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("failed to register user %d", resp.StatusCode))
	}
	return nil
}

func GetToken() (string, error) {
	content, err := ioutil.ReadFile(tokenfile)
	if err != nil {
		return "", err
	}

	token := strings.TrimRight(string(content), "\n")

	return token, nil
}

func DeleteSocket(socketID string) (error) {
	token, err := GetToken()
	if err != nil {
		return err
	}

	client := &h.Client{}
	req, err := h.NewRequest("DELETE",mysocketurl+"/socket/"+socketID, nil)
	req.Header.Add("x-access-token", token)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		responseData, _ := ioutil.ReadAll(resp.Body)
                return errors.New(fmt.Sprintf("Failed to delete socket (%d) %v", resp.StatusCode, string(responseData)))
	}

	return nil
}

func CreateSocket(name string, protected bool, username string, password string, socketType string) (*Socket, error) {
	s := &Socket{
		Name:                   name,
		ProtectedSocket:        protected,
		SocketType:             socketType,
		ProtectedUsername:      username,
		ProtectedPassword:      password,
	}

	jv, _ := json.Marshal(s)
	body := bytes.NewBuffer(jv)

	token, err := GetToken()
	if err != nil {
		return nil, err
	}

	client := &h.Client{}
	req, err := h.NewRequest("POST",mysocketurl+"/socket",  body)
	req.Header.Add("x-access-token", token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		responseData, _ := ioutil.ReadAll(resp.Body)
                return nil, errors.New(fmt.Sprintf("Failed to create socket (%d) %v", resp.StatusCode, string(responseData)))
	}

	err = json.NewDecoder(resp.Body).Decode(&s)
	if err != nil {
                return nil, errors.New("Failed to decode create socket response")
	}
	return s, nil
}

func GetSockets() ([]Socket, error) {
	sockets := []Socket{}
	token, err := GetToken()
	if err != nil {
		return nil, err
	}

	client := &h.Client{}
	req, err := h.NewRequest("GET",mysocketurl+"/connect", nil)
	req.Header.Add("x-access-token", token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
                return nil, errors.New(fmt.Sprintf("Failed to get sockets (%d)", resp.StatusCode))
	}

	err = json.NewDecoder(resp.Body).Decode(&sockets)
	if err != nil {
                return nil, errors.New("Failed to decode sockets response")
	}
	return sockets, nil
}

func GetTunnels(socketID string) ([]Tunnel, error) {
	tunnels := []Tunnel{}
	token, err := GetToken()
	if err != nil {
		return nil, err
	}

	client := &h.Client{}
	req, err := h.NewRequest("GET",mysocketurl+"/socket/"+socketID+"/tunnel", nil)
	req.Header.Add("x-access-token", token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
                return nil, errors.New(fmt.Sprintf("Failed to get tunnels (%d)", resp.StatusCode))
	}

	err = json.NewDecoder(resp.Body).Decode(&tunnels)
	if err != nil {
                return nil, errors.New("Failed to decode tunnels response")
	}
	return tunnels, nil
}

func DeleteTunnel(socketID string, tunnelID string) (error) {
	token, err := GetToken()
	if err != nil {
		return err
	}

	client := &h.Client{}
	req, err := h.NewRequest("DELETE",mysocketurl+"/socket/"+socketID+"/tunnel/"+tunnelID, nil)
	req.Header.Add("x-access-token", token)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		responseData, _ := ioutil.ReadAll(resp.Body)
                return errors.New(fmt.Sprintf("Failed to delete tunnel (%d) %v", resp.StatusCode, string(responseData)))
	}

	return nil
}

func CreateTunnel(socketID string) (*Tunnel, error) {
	t := &Tunnel{}

	jv, _ := json.Marshal(t)
	body := bytes.NewBuffer(jv)

	token, err := GetToken()
	if err != nil {
		return nil, err
	}

	client := &h.Client{}
	req, err := h.NewRequest("POST",mysocketurl+"/socket/"+socketID+"/tunnel",  body)
	req.Header.Add("x-access-token", token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		responseData, _ := ioutil.ReadAll(resp.Body)
                return nil, errors.New(fmt.Sprintf("Failed to create tunnel (%d) %v", resp.StatusCode, string(responseData)))
	}

	err = json.NewDecoder(resp.Body).Decode(&t)
	if err != nil {
                return nil, errors.New("Failed to decode create tunnel response")
	}
	return t, nil
}

func GetTunnel(socketID string, tunnelID string) (*Tunnel, error) {
	tunnel := Tunnel{}
	token, err := GetToken()
	if err != nil {
		return nil, err
	}

	client := &h.Client{}
	req, err := h.NewRequest("GET",mysocketurl+"/socket/"+socketID+"/tunnel/"+tunnelID, nil)
	req.Header.Add("x-access-token", token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
                return nil, errors.New(fmt.Sprintf("Failed to get tunnel (%d)", resp.StatusCode))
	}

	err = json.NewDecoder(resp.Body).Decode(&tunnel)
	if err != nil {
                return nil, errors.New("Failed to decode tunnel response")
	}
	return &tunnel, nil
}

func (c *client) CreateSocket(name string) error {

	return nil
}
