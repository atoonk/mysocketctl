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
    "fmt"
    "log"
    "io/ioutil"
	"github.com/atoonk/mysocketctl/go/internal/http"
	"github.com/spf13/cobra"
)



// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use: "version",
	Short: "check version",
}

var checkLatestVersionCmd = &cobra.Command{
	Use:   "check",
	Short: "Check to see if you're running the latest version",
	Run: func(cmd *cobra.Command, args []string) {
            latest_version, err := http.GetLatestVersion()
            if err != nil {
                log.Fatalf("error while checking for latest version: %v", err)
            }
            if latest_version != version {
                binary_path  := os.Args[0]
                fmt.Printf("You're running version %s\n\n",version)
                fmt.Printf("There is a newer version available (%s)!\n", latest_version)
                fmt.Printf("Please upgrade:\n%s version upgrade\n", binary_path)
            } else {
                fmt.Printf("You are up to date!\n")
                fmt.Printf("You're running version %s\n",version)
            }
	},
}
var upgradeVersionCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "upgrade the latest version",
	Run: func(cmd *cobra.Command, args []string) {
            binary_path  := os.Args[0]
            latest_version, err := http.GetLatestVersion()
            if err != nil {
                log.Fatalf("error while checking for latest version: %v", err)
            }
            if latest_version != version {
                fmt.Printf("Upgrading %s to version %s\n",binary_path,latest_version)
            } else {
                fmt.Printf("You are up to date already :)\n")
                return
            }

            latest, err := http.GetLatestBinary()
            if latest == nil {
                log.Fatalf("Error while downloading latest version %v", err)
            }

            err = ioutil.WriteFile(binary_path, latest, 0644)
	        if err != nil {
                log.Fatalf("Error while writing new file: %v", err)
	        }
            fmt.Printf("Upgrade completed\n")
	},
}


func init() {
	versionCmd.AddCommand(checkLatestVersionCmd)
	versionCmd.AddCommand(upgradeVersionCmd)
	rootCmd.AddCommand(versionCmd)
}
