package http

type registerForm struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Sshkey   string `json:"sshkey"`
}
type loginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type tokenForm struct {
	Token string `json:"token"`
}

type Socket struct {
  Tunnels               []Tunnel `json:"tunnels"`
  Username              string `json:"user_name"`
  SocketID              string `json:"socket_id"`
  SocketTcpPorts        []int  `json:"socket_tcp_ports"`
  Dnsname               string `json:"dnsname"`
  Name                  string `json:"name"`
  SocketType            string `json:"socket_type"`
  ProtectedSocket       bool   `json:"protected_socket"`
  ProtectedUsername     string `json:"protected_username"`
  ProtectedPassword     string `json:"protected_password"`
}

type Tunnel struct {
  TunnelID      string `json:"tunnel_id"`
  LocalPort     int    `json:"local_port"`
  TunnelServer  string `json:"tunnel_server"`
}

