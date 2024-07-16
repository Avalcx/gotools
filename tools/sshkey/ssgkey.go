package sshkey

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"

	"gotools/tools/ansible"
	"gotools/utils/logger"
	"gotools/utils/sshutils"

	"golang.org/x/crypto/ssh"
)

type SSHKey struct {
	Host           string
	Port           string
	User           string
	Password       string
	privateKeyPath string
	publicKeyPath  string
	privateKey     []byte
	publicKey      []byte
}

func newSSHKey() *SSHKey {
	return &SSHKey{
		Port: "22",
		User: "root",
	}
}

func (sshKey *SSHKey) generateNewSSHKeyPair() error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	sshKey.privateKey = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}

	sshKey.publicKey = ssh.MarshalAuthorizedKey(publicKey)

	return nil
}

func (sshKey *SSHKey) generatePublicKeyFromOldPrivateKey() {

	privateKeyBytes, err := os.ReadFile(sshKey.privateKeyPath)
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

	sshKey.publicKey = ssh.MarshalAuthorizedKey(publicKey)
}

func (sshKey *SSHKey) refreshPublicKeyFile() {
	if sshutils.FileIsExists(sshKey.publicKeyPath) {
		sshKey.generatePublicKeyFromOldPrivateKey()
		oldPublicKey, err := os.ReadFile(sshKey.publicKeyPath)
		if err != nil {
			logger.Fatal("failed to read id_rs.pub:%v\n", err)
		}
		if string(oldPublicKey) != string(sshKey.publicKey) {
			if err := os.WriteFile(sshKey.publicKeyPath, sshKey.publicKey, 0644); err != nil {
				logger.Fatal("failed to save public key:%v\n", err)
			}
		}
	} else {
		err := sshKey.generateNewSSHKeyPair()
		if err != nil {
			logger.Fatal("failed to generate SSH key pair:%v\n", err)
		}

		if err := os.WriteFile(sshKey.privateKeyPath, sshKey.privateKey, 0600); err != nil {
			logger.Fatal("failed to save private key:%v\n", err)
		}

		if err := os.WriteFile(sshKey.publicKeyPath, sshKey.publicKey, 0600); err != nil {
			logger.Fatal("failed to save public key:%v\n", err)
		}
	}
}

func (sshKey *SSHKey) uploadPublicKey() error {
	config := &ssh.ClientConfig{
		User: sshKey.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(sshKey.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%s", sshKey.Host, sshKey.Port)
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

	if strings.Contains(existingKeys, string(sshKey.publicKey)) {
		return nil
	}

	session, err = client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	cmd := fmt.Sprintf("mkdir -p ~/.ssh && chmod 700 ~/.ssh && echo '%s' >> ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys", sshKey.publicKey)
	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("failed to upload public key: %v", err)
	}
	return nil
}

func (sshKey *SSHKey) pushKey() {
	sshKey.refreshPublicKeyFile()
	if err := sshKey.uploadPublicKey(); err != nil {
		logger.Fatal("failed to upload public key:%v\n", err)
	}
	logger.Success("%v | User=%v | Status >> Success\n", sshKey.Host, sshKey.User)
}

func PushKeys(hostPattern, configFile, user, password string) {
	sshKeyInstance := newSSHKey()
	sshKeyInstance.privateKeyPath, sshKeyInstance.publicKeyPath = sshutils.CurrentSSHPath()
	hostsMap := ansible.ParseHostPattern(hostPattern, configFile)
	for _, hostInfo := range hostsMap {
		sshKeyInstance.Host = hostInfo.IP
		sshKeyInstance.checkPassword(password, hostInfo.Password)
		sshKeyInstance.pushKey()
	}
}

func (sshKey *SSHKey) checkPassword(cmdPassword, configPassword string) {
	if cmdPassword != "" {
		sshKey.Password = cmdPassword
		logger.Changed("使用命令行输入的密码\n")
	} else if cmdPassword == "" && configPassword == "" {
		logger.Fatal("配置文件中和命令行中配置的密码均为空\n")
	} else if cmdPassword == "" {
		sshKey.Password = configPassword
		logger.Changed("使用配置文件中的密码\n")
	}
}

func (sshKey *SSHKey) delKey() {
	key, err := os.ReadFile(sshKey.privateKeyPath)
	if err != nil {
		logger.Fatal("私钥文件读取错误: %v", err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		logger.Fatal("解析私钥失败: %v", err)
	}

	config := &ssh.ClientConfig{
		User: sshKey.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%s", sshKey.Host, sshKey.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		logger.Failed("failed to dial: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		logger.Failed("failed to create session: %v", err)
	}
	defer session.Close()
	session, err = client.NewSession()
	if err != nil {
		logger.Failed("failed to create session: %v", err)
	}
	defer session.Close()

	sshKey.generatePublicKeyFromOldPrivateKey()
	str1 := strings.ReplaceAll(string(sshKey.publicKey), "\n", "")
	str2 := strings.ReplaceAll(str1, "/", "\\/")

	cmd := fmt.Sprintf("sed -i '/%s/d' ~/.ssh/authorized_keys", str2)
	if err := session.Run(cmd); err != nil {
		logger.Failed("failed to delete public key: %v", err)
	}
	logger.Success("%v | User=%v | Status >> Success\n", sshKey.Host, sshKey.User)
}

func DelKeys(hostPattern, configFile, user string) {
	sshKeyInstance := newSSHKey()
	sshKeyInstance.privateKeyPath, sshKeyInstance.publicKeyPath = sshutils.CurrentSSHPath()
	hostsMap := ansible.ParseHostPattern(hostPattern, configFile)
	for _, hostInfo := range hostsMap {
		sshKeyInstance.Host = hostInfo.IP
		sshKeyInstance.User = user
		sshKeyInstance.delKey()
	}
}
