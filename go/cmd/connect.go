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
)

func validSocketType() bool {
	// TODO: figure out the right way to do this w/ cobra
	switch socketType {
	case "http",
		"https",
		"tcp",
		"tls":
		return true
	default:
		return false
	}
}

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Quickly connect, wrapper around sockets and tunnels",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("connect called")
		if len(name) == 0 {
			name = fmt.Sprintf("Local port %s", port)
		}
		if !validSocketType() {
			log.Fatalf("--type should one of: http, https, tcp or tls")
		}

	},
}

func init() {
	connectCmd.Flags().StringVarP(&port, "port", "p", "", "Port")
	connectCmd.Flags().StringVarP(&name, "name", "n", "", "Service name")
	connectCmd.Flags().StringVarP(&socketType, "type", "t", "http", "Socket type: http, https, tcp, tls")
	connectCmd.MarkFlagRequired("port")

	rootCmd.AddCommand(connectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// connectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// connectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
