package port

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

func udpClients(host string, portSpecs string) {
	log.Printf("UDP Test Host:%s\n", host)
	ports := parsePortSpecs(portSpecs)
	for _, port := range ports {
		isopen := udpClient(host, port)
		if isopen {
			log.Printf("Host: %v UDP Port: %v open\n", host, port)
		} else {
			log.Printf("Host: %v UDP Port: %v close\n", host, port)
		}
	}
}

func udpClient(host string, port int) bool {
	timeout := 3 * time.Second
	address := fmt.Sprintf("%s:%d", host, port)
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

func udpServers(portSpecs string) {
	ports := parsePortSpecs(portSpecs)
	var wg sync.WaitGroup
	for _, port := range ports {
		wg.Add(1)

		go func(p int) {
			defer wg.Done()
			if isLocalUDPPortOpen(p) {
				log.Printf("TCP Port %d 已占用,忽略", p)
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
		log.Printf("Failed to start UDP server: %v\n", err)
		return
	}
	defer conn.Close()
	log.Printf("UDP Port %d start success", port)

	buffer := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			log.Printf("Error reading from UDP connection: %v\n", err)
			continue
		}
		log.Printf("Received '%s' from %s\n", string(buffer[:n]), strings.Split(addr.String(), ":")[0])

		_, err = conn.WriteTo([]byte("pong"), addr)
		if err != nil {
			log.Printf("Error writing to UDP connection: %v\n", err)
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
