package port

// func httpServer(ports []int) {
// 	var wg sync.WaitGroup

// 	for _, port := range ports {
// 		wg.Add(1)

// 		go func(p int) {
// 			defer wg.Done()

// 			if isLocalTCPPortOpen(p) {
// 				log.Printf("TCP Port %d 已占用,忽略", p)
// 				return
// 			}

// 			router := gin.Default()
// 			router.Any("/", func(c *gin.Context) {
// 				c.JSON(http.StatusOK, gin.H{"status": "success", "message": fmt.Sprintf("Port %d is open", p)})
// 			})

// 			err := router.Run(fmt.Sprintf(":%d", p))
// 			if err != nil {
// 				log.Printf("Failed to start server on port %d: %v", p, err)
// 			}
// 		}(port)
// 		log.Printf("Port %d 已启动", port)
// 	}
// 	wg.Wait()
// }

// func startHTTPClient(host string, portSpecs string) {
// 	fmt.Printf("TCP 测试主机:%s\n", host)
// 	ports := parsePortSpecs(portSpecs)
// 	for _, port := range ports {
// 		isopen := httpClient(host, port)
// 		if isopen {
// 			fmt.Printf("TCP 端口: %v open\n", port)
// 		} else {
// 			fmt.Printf("TCP 端口: %v close\n", port)
// 		}
// 	}
// }

// func httpClient(host string, port int) bool {
// 	resp, err := http.Get("http://" + host + ":" + strconv.Itoa(port))
// 	if err != nil {
// 		return false
// 	}
// 	defer resp.Body.Close()
// 	return true
// }
