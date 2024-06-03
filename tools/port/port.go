package port

import (
	"log"
	"strconv"
	"strings"
)

func StartServer(portSpecs string, mode string) {
	switch mode {
	case "tcp":
		tcpServers(portSpecs)
	case "udp":
		udpServers(portSpecs)
	default:
		tcpServers(portSpecs)
	}
}

func StartClient(portSpecs string, host string, mode string) {
	switch mode {
	case "tcp":
		tcpClients(host, portSpecs)
	case "udp":
		udpClients(host, portSpecs)
	default:
		tcpClients(host, portSpecs)
	}
}

func parsePortSpecs(portSpecs string) []int {
	var ports []int
	if portSpecs == "" {
		log.Fatal("ports参数不能为空")
	}
	portSpecList := strings.Split(portSpecs, ",")
	for _, portSpec := range portSpecList {
		if strings.Contains(portSpec, "-") {
			rangeParts := strings.Split(portSpec, "-")
			if len(rangeParts) != 2 {
				log.Fatal("错误的参数格式. example: '8080-8090'")
			}

			startPort, err := strconv.Atoi(rangeParts[0])
			if err != nil {
				log.Fatal("端口参数错误: ", err)
			}

			endPort, err := strconv.Atoi(rangeParts[1])
			if err != nil {
				log.Fatal("端口参数错误: ", err)
			}

			for port := startPort; port <= endPort; port++ {
				ports = append(ports, port)
			}
		} else {
			port, err := strconv.Atoi(portSpec)
			if err != nil {
				log.Fatal("端口参数错误: ", err)
			}
			ports = append(ports, port)
		}
	}
	for _, port := range ports {
		if port >= 65535 || port <= 0 {
			log.Fatal("端口超出范围(0~65535)")
		}
	}
	return ports
}
