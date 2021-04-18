package util

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
)

func IsUsedUdpPort(port int) error {
	localaddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", localaddr)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}

func IsUsedTcpPort(port int) error {
	list, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	list.Close()
	return nil
}

func UnusedPort(network string, begin, end int) int {
	for  {
		port := (rand.Int() % (end - begin)) + begin
		if strings.ToLower(network) == "udp" {
			if IsUsedUdpPort(port) == nil {
				return port
			}
		} else {
			if IsUsedTcpPort(port) == nil {
				return port
			}
		}
	}
}

func ConnectTest(bind string, public string, ctrl int) error {
	temp := UnusedPort("tcp", 10000, 50000)
	localAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", bind, temp))
	if err != nil {
		return err
	}
	remoteAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", public, ctrl))
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", localAddr, remoteAddr)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}