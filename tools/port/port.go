package port

import (
	"gotools/utils/logger"
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
		logger.Fatal("ports参数不能为空\n")
	}
	portSpecList := strings.Split(portSpecs, ",")
	for _, portSpec := range portSpecList {
		if strings.Contains(portSpec, "-") {
			rangeParts := strings.Split(portSpec, "-")
			if len(rangeParts) != 2 {
				logger.Fatal("错误的参数格式. example: '8080-8090'\n")
			}

			startPort, err := strconv.Atoi(rangeParts[0])
			if err != nil {
				logger.Fatal("端口参数错误: %v\n", err)
			}

			endPort, err := strconv.Atoi(rangeParts[1])
			if err != nil {
				logger.Fatal("端口参数错误: %v\n", err)
			}

			for port := startPort; port <= endPort; port++ {
				ports = append(ports, port)
			}
		} else {
			port, err := strconv.Atoi(portSpec)
			if err != nil {
				logger.Fatal("端口参数错误: %v\n", err)
			}
			ports = append(ports, port)
		}
	}
	for _, port := range ports {
		if port >= 65535 || port <= 0 {
			logger.Fatal("端口超出范围(0~65535)\n")
		}
	}
	return ports
}
