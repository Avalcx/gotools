package netutils

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

func IsIPAddress(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil
}

func IsIPRange(ipRange string) bool {
	re := regexp.MustCompile(`^(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})-(\d{1,3})$`)
	matches := re.FindStringSubmatch(ipRange)
	if len(matches) != 3 {
		return false
	}
	startIP := matches[1]
	endSuffix := matches[2]

	if !IsIPAddress(startIP) {
		return false
	}

	end, err := strconv.Atoi(endSuffix)
	if err != nil || end < 0 || end > 255 {
		return false
	}

	startParts := strings.Split(startIP, ".")
	startLastOctet, _ := strconv.Atoi(startParts[3])
	return startLastOctet <= end
}

func ParseIPRange(ipRange string) ([]net.IP, error) {
	if strings.Contains(ipRange, "-") {
		parts := strings.Split(ipRange, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid range format")
		}

		startIP := net.ParseIP(parts[0])
		if startIP == nil {
			return nil, fmt.Errorf("invalid start IP")
		}

		startParts := strings.Split(parts[0], ".")
		if len(startParts) != 4 {
			return nil, fmt.Errorf("invalid IP format")
		}

		startLastOctet, err := strconv.Atoi(startParts[3])
		if err != nil {
			return nil, fmt.Errorf("invalid last octet in start IP")
		}

		endLastOctet, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid end of range")
		}

		if endLastOctet < startLastOctet || endLastOctet > 255 {
			return nil, fmt.Errorf("invalid range")
		}

		var ips []net.IP
		for i := startLastOctet; i <= endLastOctet; i++ {
			ip := fmt.Sprintf("%s.%s.%s.%d", startParts[0], startParts[1], startParts[2], i)
			ips = append(ips, net.ParseIP(ip))
		}
		return ips, nil
	} else {
		var ips []net.IP
		ips = append(ips, net.ParseIP(ipRange))
		return ips, nil
	}
}
