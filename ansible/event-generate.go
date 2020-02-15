package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	s := bufio.NewScanner(os.Stdin)
	if len(os.Args) != 1+1 {
		fmt.Fprintf(os.Stderr, "not enough arguments\nusage: %s [ssh user]\n", os.Args[0])
		return
	}
	user := os.Args[1]

	fmt.Printf("[recorders]\n")

	for s.Scan() {
		host := s.Text()
		if strings.HasPrefix(host, "#") {
			// Is a comment
			continue
		}
		if strings.HasSuffix(host, ":0") {
			// Is a "system" host
			continue
		}
		hostname, roomId, roomType, err := discover(host, user)
		if err != nil {
			fmt.Fprintf(os.Stderr, "host %s discovery failed: %v\n", host, err)
		}
		fmt.Printf("%s ansible_host=%s ansible_user=%s room_id=%s room_type=%s\n", hostname, host, user, roomId, roomType)
	}
}

func discover(host, user string) (hostname, roomId, roomType string, err error) {
	cmd := exec.Command("/usr/bin/ssh", "-l", user, host, "sh", "-c", "hostname; cat room_id; cat room_type; exit 0")
	fmt.Fprintf(os.Stderr, "command: %s\n", cmd)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return
	}
	hostname, err = out.ReadString('\n')
	if err != nil {
		return
	}
	hostname = strings.TrimSpace(hostname)
	roomId, _ = out.ReadString('\n')
	roomId = strings.TrimSpace(roomId)
	roomType, _ = out.ReadString('\n')
	roomType = strings.TrimSpace(roomType)
	return
}
