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

	"github.com/atoonk/mysocketctl/go/internal/http"
	"github.com/spf13/cobra"
	//"github.com/davecgh/go-spew/spew"
)

// socketCmd represents the socket command
var socketCmd = &cobra.Command{
	Use:   "socket",
	Short: "Manage your global sockets",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

// socketsListCmd represents the socket ls command
var socketsListCmd = &cobra.Command{
	Use:   "ls",
	Short: "List your global sockets",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		sockets, err := http.GetSockets()

		if err != nil {
			log.Fatalf("error: %v", err)
		}

		fmt.Printf("%-36s | %-40s | %-40s\n", "Socket ID", "DNS Name", "Name")
		for _, s := range sockets {
			fmt.Printf("%-36s | %-40s | %-40s\n", s.SocketID, s.Dnsname, s.Name)
		}

	},
}

// socketCreateCmd represents the socket create command
var socketCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new global socket",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

		s, err := http.CreateSocket(name, protected, username, password, socketType)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		fmt.Printf("%-36s | %-40s | %-40s\n", "Socket ID", "DNS Name", "Name")
		fmt.Printf("%-36s | %-40s | %-40s\n", s.SocketID, s.Dnsname, s.Name)
	},
}

// socketDeleteCmd represents the socket delete command
var socketDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a global socket",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if socketID == "" {
			log.Fatalf("error: invalid socketid")
		}

		err := http.DeleteSocket(socketID)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		fmt.Println("Socket deleted")
	},
}

func init() {
	rootCmd.AddCommand(socketCmd)
	socketCmd.AddCommand(socketsListCmd)
	socketCmd.AddCommand(socketCreateCmd)
	socketCmd.AddCommand(socketDeleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// socketCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// socketCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	socketCreateCmd.Flags().StringVarP(&name, "name", "n", "", "Socket name")
	socketCreateCmd.Flags().BoolVarP(&protected, "protected", "p", false, "Protected, default no")
	socketCreateCmd.Flags().StringVarP(&username, "username", "u", "", "Username, required when protected set to true")
	socketCreateCmd.Flags().StringVarP(&password, "password", "", "", "Password, required when protected set to true")
	socketCreateCmd.Flags().StringVarP(&socketType, "type", "t", "http", "Socket type, defaults to http")
	socketCreateCmd.MarkFlagRequired("name")
	socketDeleteCmd.Flags().StringVarP(&socketID, "id", "i", "", "Socket ID")
	socketDeleteCmd.MarkFlagRequired("id")
}
