package port

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

func StartServer(portSpecs string) {
	gin.SetMode(gin.ReleaseMode)
	ports := parsePortSpecs(portSpecs)
	startHTTPServers(ports)
}

func StartClient(portSpecs string, host string) {
	fmt.Printf("测试主机:%s\n", host)
	ports := parsePortSpecs(portSpecs)
	for _, port := range ports {
		isopen := get(host, port)
		if isopen {
			fmt.Printf("端口: %v open\n", port)
		} else {
			fmt.Printf("端口: %v close\n", port)
		}
	}
}

func parsePortSpecs(portSpecs string) []int {
	var ports []int

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

func startHTTPServers(ports []int) {
	var wg sync.WaitGroup

	for _, port := range ports {
		wg.Add(1)

		go func(p int) {
			defer wg.Done()

			if isLocalPortOpen(p) {
				log.Printf("Port %d 已占用,忽略", p)
				return
			}

			router := gin.Default()
			router.Any("/", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "success", "message": fmt.Sprintf("Port %d is open", p)})
			})

			err := router.Run(fmt.Sprintf(":%d", p))
			if err != nil {
				log.Printf("Failed to start server on port %d: %v", p, err)
			}
		}(port)
		log.Printf("Port %d 已启动", port)
	}

	wg.Wait()
}

func isLocalPortOpen(port int) bool {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return true // 端口被占用
	}
	defer listener.Close()
	return false // 端口可用
}

func get(host string, port int) bool {
	resp, err := http.Get("http://" + host + ":" + strconv.Itoa(port))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return true
}
