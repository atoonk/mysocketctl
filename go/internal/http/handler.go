package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	h "net/http"
)

const mysocketurl = "https://api.mysocket.io"

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
	return c, nil
}

// Register ...
func Register(name, email, password, sshkey string) error {
	form := registerForm{Name: name, Email: email, Password: password, Sshkey: sshkey}
	fmt.Printf("%s %s %s %s", name, email, password, sshkey)
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

func (c *client) CreateSocket(name string) error {

	return nil
}
