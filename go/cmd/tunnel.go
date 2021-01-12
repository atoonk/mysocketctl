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
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/atoonk/mysocketctl/go/internal/http"
	"github.com/atoonk/mysocketctl/go/internal/ssh"
)

// tunnelCmd represents the tunnel command
var tunnelCmd = &cobra.Command{
	Use:   "tunnel",
	Short: "Manage your tunnels",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

var tunnelListCmd = &cobra.Command{
        Use:   "ls",
        Short: "List your tunnels",
        Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
        Run: func(cmd *cobra.Command, args []string) {
               if socketID == "" {
			log.Fatalf("error: --socket_id required")
		}
                tunnels, err := http.GetTunnels(socketID)

                if err != nil {
                        log.Fatalf("error: %v", err)
                }

                fmt.Printf("%-36s | %-16s | %-10s\n", "Tunnel ID", "Tunnel Server", "Relay Port")
                for _, t := range tunnels {
                                fmt.Printf("%-36s | %-16s | %-10d\n", t.TunnelID, t.TunnelServer, t.LocalPort)
                }
	},
}

var tunnelDeleteCmd = &cobra.Command{
        Use:   "delete",
        Short: "Delete a tunnel",
        Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
        Run: func(cmd *cobra.Command, args []string) {
                if socketID == "" {
                        log.Fatalf("error: invalid socket_id")
                }
                if tunnelID == "" {
                        log.Fatalf("error: invalid tunnel_id")
                }

                err := http.DeleteTunnel(socketID, tunnelID)
                if err != nil {
                        log.Fatalf("error: %v", err)
                }

                fmt.Println("Tunnel deleted")
        },
}

var tunnelCreateCmd = &cobra.Command{
        Use:   "create",
        Short: "Create a tunnel",
        Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
        Run: func(cmd *cobra.Command, args []string) {
                if socketID == "" {
                        log.Fatalf("error: empty socket_id not allowed")
                }

                t, err := http.CreateTunnel(socketID)
                if err != nil {
                        log.Fatalf("error: %v", err)
                }

                fmt.Printf("%-36s | %-16s | %-10s\n", "Tunnel ID", "Tunnel Server", "Relay Port")
                fmt.Printf("%-36s | %-16s | %-10d\n", t.TunnelID, t.TunnelServer, t.LocalPort)
        },
}

var tunnelConnectCmd = &cobra.Command{
        Use:   "connect",
        Short: "Connect a tunnel",
        Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
        Run: func(cmd *cobra.Command, args []string) {
                if socketID == "" {
                        log.Fatalf("error: invalid socket_id")
                }
                if tunnelID == "" {
                        log.Fatalf("error: invalid tunnel_id")
                }
                if port < 1 {
                        log.Fatalf("error: invalid port")
                }

                userID, _, err := http.GetUserID()
                if err != nil {
                        log.Fatalf("error: %v", err)
                }

		userIDStr := *userID
		ssh.SshConnect(userIDStr, socketID, tunnelID, port, identityFile)
        },
}

func init() {
	rootCmd.AddCommand(tunnelCmd)
        tunnelCmd.AddCommand(tunnelListCmd)
        tunnelCmd.AddCommand(tunnelCreateCmd)
        tunnelCmd.AddCommand(tunnelDeleteCmd)
        tunnelCmd.AddCommand(tunnelConnectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tunnelCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tunnelCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

        tunnelDeleteCmd.Flags().StringVarP(&tunnelID, "tunnel_id", "t", "", "Tunnel ID")
        tunnelDeleteCmd.Flags().StringVarP(&socketID, "socket_id", "s", "", "Socket ID")
        tunnelDeleteCmd.MarkFlagRequired("tunnel_id")
        tunnelDeleteCmd.MarkFlagRequired("socket_id")
        tunnelListCmd.Flags().StringVarP(&socketID, "socket_id", "s", "", "Socket ID")
        tunnelListCmd.MarkFlagRequired("socket_id")
        tunnelCreateCmd.Flags().StringVarP(&socketID, "socket_id", "s", "", "Socket ID")
        tunnelCreateCmd.MarkFlagRequired("socket_id")
        tunnelConnectCmd.Flags().StringVarP(&tunnelID, "tunnel_id", "t", "", "Tunnel ID")
        tunnelConnectCmd.Flags().StringVarP(&socketID, "socket_id", "s", "", "Socket ID")
        tunnelConnectCmd.Flags().StringVarP(&identityFile, "identity_file", "i", "", "Identity File")
        tunnelConnectCmd.Flags().IntVarP(&port, "port", "p", 0, "Port number")
        tunnelConnectCmd.MarkFlagRequired("tunnel_id")
        tunnelConnectCmd.MarkFlagRequired("socket_id")
        tunnelConnectCmd.MarkFlagRequired("port")
}
