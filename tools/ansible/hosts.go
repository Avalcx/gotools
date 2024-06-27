package ansible

import (
	"bufio"
	"gotools/utils/logger"
	"gotools/utils/netutils"
	"os"
	"strings"
)

// 设置默认值
func (h *HostInfo) SetDefaults() {
	if h.Port == "" {
		h.Port = "22"
	}
	if h.User == "" {
		h.User = "root"
	}
	if h.Password == "" {
		h.Password = ""
	}
}

func parseGroupFromFile(configFile string) (map[string][]*HostInfo, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	hosts := make(map[string][]*HostInfo)
	var currentGroup string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentGroup = strings.TrimPrefix(strings.TrimSuffix(line, "]"), "[")
			continue
		}

		parts := strings.Fields(line)
		// 格式A
		// [group]
		// 192.168.1.1
		if len(parts) == 1 {
			if netutils.IsIPAddress(parts[0]) {
				hostInfo := HostInfo{
					Hostname: parts[0],
					IP:       parts[0],
				}
				hostInfo.SetDefaults()
				hosts[currentGroup] = append(hosts[currentGroup], &hostInfo)
			} else {
				logger.Fatal("%v: 不是正确的ip格式\n", parts[0])
			}
			// 格式B
			// host1 ansible_host=172.168.101.71 ansible_port=22 ansible_user=root ansible_ssh_pass=317210
		} else if len(parts) >= 2 {
			hostInfo := HostInfo{}
			for _, item := range parts {
				if !strings.Contains(item, "=") {
					hostInfo = HostInfo{Hostname: parts[0]}
				} else {
					keyValue := strings.Split(item, "=")
					if len(keyValue) != 2 {
						continue
					}
					key := strings.TrimSpace(keyValue[0])
					value := strings.TrimSpace(keyValue[1])
					switch key {
					case "ansible_host":
						hostInfo.IP = value
					case "ansible_port":
						hostInfo.Port = value
					case "ansible_user":
						hostInfo.User = value
					case "ansible_ssh_pass":
						hostInfo.Password = value
					}
				}
			}
			hostInfo.SetDefaults()
			hosts[currentGroup] = append(hosts[currentGroup], &hostInfo)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return hosts, nil
}

// 解析/etc/ansible/hosts文件
func ParseHostPattern(hostPattern, configFile string) []*HostInfo {
	hostInfoInstanceList := make([]*HostInfo, 0, 10)

	// IP地址段
	if netutils.IsIPRange(hostPattern) {
		ipRange, err := netutils.ParseIPRange(hostPattern)
		if err != nil {
			logger.Fatal("%v", err)
		}
		for _, ip := range ipRange {
			hostInfoInstance := NewHostInfo()
			hostInfoInstance.IP = ip.String()
			hostInfoInstanceList = append(hostInfoInstanceList, hostInfoInstance)
		}
		return hostInfoInstanceList
		// IP
	} else if netutils.IsIPAddress(hostPattern) {
		hostInfoInstance := NewHostInfo()
		hostInfoInstance.IP = hostPattern
		hostInfoInstanceList = append(hostInfoInstanceList, hostInfoInstance)
		return hostInfoInstanceList
		// 组名
	} else {
		hostInfoInstanceMap, err := parseGroupFromFile(configFile)
		if err != nil {
			logger.Fatal("%v", err)
		}
		for group, infos := range hostInfoInstanceMap {
			if group == hostPattern {
				hostInfoInstanceList = append(hostInfoInstanceList, infos...)
			}
		}
		return hostInfoInstanceList
	}
}
