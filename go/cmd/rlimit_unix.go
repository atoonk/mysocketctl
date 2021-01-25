// +build !windows

package cmd

import (
	"syscall"
	"fmt"
	"log"
)

func SetRlimit() {
	// check open file limit
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
	        fmt.Println("Error Getting Rlimit ", err)
	}

	if rLimit.Cur < rLimit.Max {
	        rLimit.Cur = rLimit.Max
	        err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	        if err != nil {
	                log.Println("Error Setting Rlimit ", err)
	        }
	}
}
