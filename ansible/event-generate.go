package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
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
		ip := net.ParseIP(host).To16()
		if ip == nil {
			// Not a valid IPv6, skip
			continue
		}
		if !bytes.Equal(ip[8:14], []byte{0, 0, 0, 0, 0, 0}) {
			// First subnet is the recorders subnet
			// Other subnets are "system" hosts
			continue
		}
		hostname, roomType, roomId, webcamName, err := discover(host, user)
		if err != nil {
			fmt.Fprintf(os.Stderr, "host %s discovery failed: %v\n", host, err)
		}
		fmt.Printf("%s ansible_host=%s ansible_user=%s room_type=%s room_id=%s webcam_name=%s\n", hostname, host, user, roomType, roomId, webcamName)
	}
}

const script = `
hostname;
(cat room_type || echo) | head -n 1;
(cat room_id || echo) | head -n 1;
(ls /dev/v4l/by-id/ | grep C920 | grep index0 || echo) | cut -d - -f 1,2 | head -n 1;
(cat /sys/class/net/en*/address || echo) | head -n 1;
`

func discover(host, user string) (hostname, roomType, roomId, webcamName string, err error) {
	cmd := exec.Command("/usr/bin/ssh", "-o", "ConnectTimeout=3", "-l", user, host, "sh", "-c", script)
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
	roomType, _ = out.ReadString('\n')
	roomType = strings.TrimSpace(roomType)
	roomId, _ = out.ReadString('\n')
	roomId = strings.TrimSpace(roomId)
	webcamName, _ = out.ReadString('\n')
	webcamName = strings.TrimSpace(webcamName)
	return
}
