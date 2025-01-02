package port

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/Avalcx/gotools/utils/logger"
)

func (portInfo *Port) tcpClients() {
	for _, port := range portInfo.Ports {
		isopen := portInfo.tcpClient(port)
		if isopen {
			logger.Success("%v | Port=%v | TCP | Status >> Open\n", portInfo.Host, port)
		} else {
			logger.Failed("%v | Port=%v | TCP | Status >> Close\n", portInfo.Host, port)
		}
	}
}

func (portInfo *Port) tcpClient(port int) bool {
	address := fmt.Sprintf("%s:%d", portInfo.Host, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return false
	}
	defer conn.Close()

	_, err = conn.Write([]byte("ping"))
	if err != nil {
		return false
	}

	buffer := make([]byte, 1024)
	_, err = conn.Read(buffer)
	return err == nil
}

func (portInfo *Port) tcpServers() {
	var wg sync.WaitGroup
	for _, port := range portInfo.Ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			if isLocalTCPPortOpen(p) {
				logger.Ignore("TCP Port %d 已占用,忽略", p)
				return
			}
			portInfo.tcpServer(p)
		}(port)

	}
	wg.Wait()
}

func (portInfo *Port) tcpServer(port int) {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Failed("Failed to start TCP server: %v\n", err)
		return
	}
	defer listener.Close()
	logger.Success("TCP Port %d start success\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Failed("Error accepting TCP connection: %v\n", err)
			continue
		}
		go func(conn net.Conn) {
			defer conn.Close()

			buffer := make([]byte, 1024)
			n, err := conn.Read(buffer)
			if err != nil {
				logger.Failed("Error reading from TCP connection: %v\n", err)
				return
			}
			logger.Ignore("TCP Port: %d Received '%s' from %s\n", port, string(buffer[:n]), strings.Split(conn.RemoteAddr().String(), ":")[0])

			_, err = conn.Write([]byte("pong"))
			if err != nil {
				logger.Failed("Error writing to TCP connection: %v\n", err)
			}
		}(conn)
	}
}

func isLocalTCPPortOpen(port int) bool {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return true
	}
	defer listener.Close()
	return false
}
