/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"os"
	"os/signal"
	"syscall"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/atoonk/mysocketctl/go/internal/http"
	"github.com/atoonk/mysocketctl/go/internal/ssh"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Quickly connect, wrapper around sockets and tunnels",
	Run: func(cmd *cobra.Command, args []string) {
                if protected {
                        if username == "" {
                                log.Fatalf("error: --username required when using --protected")
                        }
                        if password == "" {
                                log.Fatalf("error: --password required when using --protected")
                        }
                }

                if name == "" {
                        log.Fatalf("error: empty name not allowed")
                }

                if socketType != "http" && socketType != "https" && socketType != "tcp" && socketType != "tls" {
                        log.Fatalf("error: --type should be either http, https, tcp or tls")
                }

                c, err := http.CreateConnection(name, protected, username, password, socketType)
                if err != nil {
                        log.Fatalf("error: %v", err)
                }

		fmt.Printf("%-36s | %-40s | %-40s\n", "Socket ID", "DNS Name", "Name")
		fmt.Printf("%-36s | %-40s | %-40s\n", c.SocketID, c.Dnsname, c.Name)

                userID, _, err2 := http.GetUserID()
                if err2 != nil {
                        log.Fatalf("error: %v", err2)
                }

                userIDStr := *userID
		time.Sleep(2 * time.Second)
		ch := make(chan os.Signal)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-ch
			fmt.Println("cleaning up...")
			err = http.DeleteSocket(c.SocketID)
			if err != nil {
				log.Fatalf("error: %v", err)
			}
			os.Exit(0)
		}()

                ssh.SshConnect(userIDStr, c.SocketID, c.Tunnels[0].TunnelID, port, identityFile)
		fmt.Println("cleaning up...")
                err = http.DeleteSocket(c.SocketID)
                if err != nil {
                        log.Fatalf("error: %v", err)
                }
	},
}

func init() {
	connectCmd.Flags().IntVarP(&port, "port", "p", 0, "Port")
	connectCmd.Flags().StringVarP(&name, "name", "n", "", "Service name")
        connectCmd.Flags().BoolVarP(&protected, "protected", "", false, "Protected, default no")
        connectCmd.Flags().StringVarP(&username, "username", "u", "", "Username, required when protected set to true")
        connectCmd.Flags().StringVarP(&password, "password", "", "", "Password, required when protected set to true")
	connectCmd.Flags().StringVarP(&socketType, "type", "t", "http", "Socket type: http, https, tcp, tls")
        connectCmd.Flags().StringVarP(&identityFile, "identity_file", "i", "", "Identity File")
	connectCmd.MarkFlagRequired("port")
        connectCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(connectCmd)
}
