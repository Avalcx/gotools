package port

import (
	"fmt"
	"gotools/utils/logger"
	"net"
	"strings"
	"sync"
)

func tcpClients(host string, portSpecs string) {
	ports := parsePortSpecs(portSpecs)
	for _, port := range ports {
		isopen := tcpClient(host, port)
		if isopen {
			logger.Success("%v | Port=%v | TCP | Status >> Open\n", host, port)
		} else {
			logger.Failed("%v | Port=%v | TCP | Status >> Close\n", host, port)
		}
	}
}

func tcpClient(host string, port int) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		logger.Failed("Error: %v\n", err)
		return false
	}
	defer conn.Close()

	_, err = conn.Write([]byte("ping"))
	if err != nil {
		logger.Failed("Error writing to TCP connection: %v\n", err)
		return false
	}

	buffer := make([]byte, 1024)
	_, err = conn.Read(buffer)
	if err != nil {
		logger.Failed("Error reading from TCP connection: %v\n", err)
		return false
	}
	return true
}

func tcpServers(portSpecs string) {
	ports := parsePortSpecs(portSpecs)
	var wg sync.WaitGroup
	for _, port := range ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			if isLocalTCPPortOpen(p) {
				logger.Ignore("TCP Port %d 已占用,忽略", p)
				return
			}
			tcpServer(p)
		}(port)

	}
	wg.Wait()
}

func tcpServer(port int) {
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
