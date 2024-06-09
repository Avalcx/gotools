package sshkey

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"gotools/utils/logger"

	"golang.org/x/crypto/ssh"
)

func currentSSHPath() (string, string) {
	currentUser, err := user.Current()
	if err != nil {
		logger.Fatal("failed to get current user: %v\n", err)
	}
	sshDir := filepath.Join(currentUser.HomeDir, ".ssh")
	privateKeyPath := filepath.Join(sshDir, "id_rsa")
	publicKeyPath := filepath.Join(sshDir, "id_rsa.pub")
	return privateKeyPath, publicKeyPath
}

func generateNewSSHKeyPair() ([]byte, []byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	publicKeyBytes := ssh.MarshalAuthorizedKey(publicKey)

	return privateKeyPEM, publicKeyBytes, nil
}

func generateOldPublicKey() []byte {
	privateKeyPath, _ := currentSSHPath()
	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		logger.Fatal("failed to read private key file:%v\n", err)
	}

	block, _ := pem.Decode(privateKeyBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		logger.Fatal("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		logger.Fatal("failed to parse private key:%v\n", err)
	}

	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		logger.Fatal("failed to generate public key:%v\n", err)
	}

	publicKeyBytes := ssh.MarshalAuthorizedKey(publicKey)
	return publicKeyBytes
}

func uploadPublicKey(user, host, password, publicKey string) error {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:22", host)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("failed to dial: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	var buf bytes.Buffer
	session.Stdout = &buf
	if err := session.Run("cat ~/.ssh/authorized_keys 2>/dev/null || true"); err != nil {
		return fmt.Errorf("failed to read authorized_keys: %v", err)
	}

	existingKeys := buf.String()

	if strings.Contains(existingKeys, publicKey) {
		return nil
	}

	session, err = client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	cmd := fmt.Sprintf("mkdir -p ~/.ssh && chmod 700 ~/.ssh && echo '%s' >> ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys", publicKey)
	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("failed to upload public key: %v", err)
	}
	return nil
}

func fileIsExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func parseHostSpecs(rangeStr string) ([]net.IP, error) {
	if strings.Contains(rangeStr, "-") {
		parts := strings.Split(rangeStr, "-")
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
		ips = append(ips, net.ParseIP(rangeStr))
		return ips, nil
	}
}

func pushKey(host string, user string, password string) {
	privateKeyPath, publicKeyPath := currentSSHPath()

	var privateKey, publicKey []byte
	if fileIsExists(publicKeyPath) {
		newPublicKey := generateOldPublicKey()
		oldPublicKey, err := os.ReadFile(publicKeyPath)
		if err != nil {
			logger.Fatal("failed to read id_rs.pub:%v\n", err)
		}
		if string(oldPublicKey) == string(newPublicKey) {
			publicKey = oldPublicKey
		} else {
			publicKey = newPublicKey
			if err := os.WriteFile(publicKeyPath, publicKey, 0644); err != nil {
				logger.Fatal("failed to save public key:%v\n", err)
			}
		}
	} else {
		var err error
		privateKey, publicKey, err = generateNewSSHKeyPair()
		if err != nil {
			logger.Fatal("failed to generate SSH key pair:%v\n", err)
		}

		if err := os.WriteFile(privateKeyPath, privateKey, 0600); err != nil {
			logger.Fatal("failed to save private key:%v\n", err)
		}

		if err := os.WriteFile(publicKeyPath, publicKey, 0600); err != nil {
			logger.Fatal("failed to save public key:%v\n", err)
		}
	}

	if err := uploadPublicKey(user, host, password, string(publicKey)); err != nil {
		logger.Fatal("failed to upload public key:%v\n", err)
	}
	logger.Success("%v | User=%v |Status >> Success\n", host, user)
}

func PushKeys(hostsSlice []string, user string, password string) {
	for _, hosts := range hostsSlice {
		ips, err := parseHostSpecs(hosts)
		if err != nil {
			logger.Fatal("hosts format error:%v\n", err)
		}
		for _, ip := range ips {
			pushKey(ip.String(), user, password)
		}
	}
}
