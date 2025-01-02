package port

import (
	"strconv"
	"strings"

	"github.com/Avalcx/gotools/utils/logger"
)

type Port struct {
	PortSpecs string
	Ports     []int
	Protocol  string
	Host      string
}

func (portInfo *Port) StartServer() {
	switch portInfo.Protocol {
	case "tcp":
		portInfo.tcpServers()
	case "udp":
		portInfo.udpServers()
	default:
		portInfo.tcpServers()
	}
}

func (portInfo *Port) StartClient() {
	switch portInfo.Protocol {
	case "tcp":
		portInfo.tcpClients()
	case "udp":
		portInfo.udpClients()
	default:
		portInfo.tcpClients()
	}
}

func (portInfo *Port) ParsePortSpecs() {
	var ports []int
	if portInfo.PortSpecs == "" {
		logger.Fatal("ports参数不能为空\n")
	}
	portSpecList := strings.Split(portInfo.PortSpecs, ",")
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
	portInfo.Ports = ports
}
