package port

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

func tcpClients(host string, portSpecs string) {
	ports := parsePortSpecs(portSpecs)
	for _, port := range ports {
		isopen := tcpClient(host, port)
		if isopen {
			log.Printf("Host: %v TCP Port: %v open\n", host, port)
		} else {
			log.Printf("Host: %v TCP Port: %v close\n", host, port)
		}
	}
}

func tcpClient(host string, port int) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return false
	}
	defer conn.Close()

	_, err = conn.Write([]byte("ping"))
	if err != nil {
		log.Printf("Error writing to TCP connection: %v\n", err)
		return false
	}

	buffer := make([]byte, 1024)
	_, err = conn.Read(buffer)
	if err != nil {
		log.Printf("Error reading from TCP connection: %v\n", err)
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
				log.Printf("TCP Port %d 已占用,忽略", p)
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
		log.Printf("Failed to start TCP server: %v\n", err)
		return
	}
	defer listener.Close()
	log.Printf("TCP Port %d start success", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting TCP connection: %v\n", err)
			continue
		}
		go func(conn net.Conn) {
			defer conn.Close()

			buffer := make([]byte, 1024)
			n, err := conn.Read(buffer)
			if err != nil {
				log.Printf("Error reading from TCP connection: %v\n", err)
				return
			}
			log.Printf("TCP Port: %d Received '%s' from %s\n", port, string(buffer[:n]), strings.Split(conn.RemoteAddr().String(), ":")[0])

			_, err = conn.Write([]byte("pong"))
			if err != nil {
				log.Printf("Error writing to TCP connection: %v\n", err)
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
