package port

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/Avalcx/gotools/utils/logger"
)

func (portInfo *Port) udpClients() {
	logger.Ignore("UDP Test Host:%s\n", portInfo.Host)
	for _, port := range portInfo.Ports {
		isopen := portInfo.udpClient(port)
		if isopen {
			logger.Success("%v | Port=%v | UDP | Status >> Open\n", portInfo.Host, port)
		} else {
			logger.Failed("%v | Port=%v | UDP | Status >> Close\n", portInfo.Host, port)
		}
	}
}

func (portInfo *Port) udpClient(port int) bool {
	timeout := 3 * time.Second
	address := fmt.Sprintf("%s:%d", portInfo.Host, port)
	conn, err := net.DialTimeout("udp", address, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()

	deadline := time.Now().Add(timeout)
	conn.SetDeadline(deadline)

	_, err = conn.Write([]byte("ping"))
	if err != nil {
		return false
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return false
	}

	if n > 0 {
		return true
	}

	return false
}

func (portInfo *Port) udpServers() {
	var wg sync.WaitGroup
	for _, port := range portInfo.Ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			if isLocalUDPPortOpen(p) {
				logger.Ignore("UDP Port %d 已占用,忽略\n", p)
				return
			}
			udpServer(p)
		}(port)

	}
	wg.Wait()
}

func udpServer(port int) {
	addr := fmt.Sprintf(":%d", port)
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		logger.Failed("Failed to start UDP server: %v\n", err)
		return
	}
	defer conn.Close()
	logger.Success("UDP Port %d start success\n", port)

	buffer := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			logger.Failed("Error reading from UDP connection: %v\n", err)
			continue
		}
		logger.Ignore("Received '%s' from %s\n", string(buffer[:n]), strings.Split(addr.String(), ":")[0])

		_, err = conn.WriteTo([]byte("pong"), addr)
		if err != nil {
			logger.Failed("Error writing to UDP connection: %v\n", err)
		}
	}
}

func isLocalUDPPortOpen(port int) bool {
	addr := fmt.Sprintf(":%d", port)
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return true
	}
	defer conn.Close()
	return false
}
