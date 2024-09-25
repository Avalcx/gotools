package sshkey

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"gotools/tools/ansible"
	"gotools/utils/logger"
	"gotools/utils/sshutils"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

const (
	passwordLength = 16
	characters     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+"
	passwordFile   = "new_password.txt"
)

type SSHKey struct {
	Host           string
	Port           string
	User           string
	Password       string
	NewPassword    string
	sshPath        string
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

		if err := os.MkdirAll(sshKey.sshPath, 0700); err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
			return
		}

		if err := os.WriteFile(sshKey.privateKeyPath, sshKey.privateKey, 0600); err != nil {
			logger.Fatal("failed to save private key:%v\n", err)
		}

		if err := os.WriteFile(sshKey.publicKeyPath, sshKey.publicKey, 0600); err != nil {
			logger.Fatal("failed to save public key:%v\n", err)
		}
	}
}

func (sshKey *SSHKey) runCMD(cmd string) error {
	key, err := os.ReadFile(sshKey.privateKeyPath)
	if err != nil {
		return fmt.Errorf("私钥文件读取错误: %v", err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return fmt.Errorf("解析私钥失败: %v", err)
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
		if strings.Contains(err.Error(), "[none publickey]") {
			return fmt.Errorf("没有这个主机的公钥: %v", sshKey.Host)
		} else {
			return fmt.Errorf("failed to dial: %v", err)
		}
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()
	session, err = client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	if err := session.Run(cmd); err != nil {
		return err
	}
	return nil
}

func (sshKey *SSHKey) uploadPublicKey() (int, error) {
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
		return 0, fmt.Errorf("failed to dial: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return 0, fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	var buf bytes.Buffer
	session.Stdout = &buf
	if err := session.Run("cat ~/.ssh/authorized_keys 2>/dev/null || true"); err != nil {
		return 0, fmt.Errorf("failed to read authorized_keys: %v", err)
	}

	existingKeys := buf.String()

	if strings.Contains(existingKeys, string(sshKey.publicKey)) {
		return 1, nil
	}

	session, err = client.NewSession()
	if err != nil {
		return 0, fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	cmd := fmt.Sprintf("mkdir -p ~/.ssh && chmod 700 ~/.ssh && echo '%s' >> ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys", sshKey.publicKey)
	if err := session.Run(cmd); err != nil {
		return 0, fmt.Errorf("failed to upload public key: %v", err)
	}
	return 0, nil
}

func (sshKey *SSHKey) pushKey() {
	sshKey.refreshPublicKeyFile()
	if status, err := sshKey.uploadPublicKey(); err != nil {
		logger.Failed("%v | User=%v | Status >> Failed\n%v\n", sshKey.Host, sshKey.User, err)
	} else {
		if status == 1 {
			logger.Success("%v | User=%v | Status >> Success\n", sshKey.Host, sshKey.User)
		} else {
			logger.Changed("%v | User=%v | Status >> CHANGED\n", sshKey.Host, sshKey.User)
		}
	}
}

func PushKeys(hostPattern, configFile, user, password string) {
	sshKeyInstance := newSSHKey()
	sshKeyInstance.privateKeyPath, sshKeyInstance.publicKeyPath, sshKeyInstance.sshPath = sshutils.CurrentSSHPath()
	hostsMap := ansible.ParseHostPattern(hostPattern, configFile)
	isloged := false
	for _, hostInfo := range hostsMap {
		sshKeyInstance.selectPassword(password, hostInfo.Password, isloged)
		sshKeyInstance.Host = hostInfo.IP
		sshKeyInstance.pushKey()
		isloged = true
	}
	logger.Success("已完成免密\n")
}

func (sshKey *SSHKey) selectPassword(cmdPassword, configPassword string, isloged bool) {
	if cmdPassword != "" {
		sshKey.Password = cmdPassword
		if !isloged {
			logger.Changed("使用命令行输入的密码\n")
		}
	} else if cmdPassword == "" && configPassword == "" {
		if !isloged {
			logger.Fatal("配置文件中和命令行中配置的密码均为空\n")
		}
	} else if cmdPassword == "" {
		sshKey.Password = configPassword
		if !isloged {
			logger.Changed("使用配置文件中的密码\n")
		}
	}
}

func (sshKey *SSHKey) delKey() error {
	cmd := fmt.Sprintf("sed -i '/%s/d' ~/.ssh/authorized_keys", strings.ReplaceAll(strings.ReplaceAll(string(sshKey.publicKey), "\n", ""), "/", "\\/"))
	if err := sshKey.runCMD(cmd); err != nil {
		return fmt.Errorf("failed to delete public key: %v", err)
	}
	logger.Success("%v | User=%v | Status >> Success\n", sshKey.Host, sshKey.User)
	return nil
}

func DelKeys(hostPattern, configFile, user string) {
	sshKeyInstance := newSSHKey()
	sshKeyInstance.privateKeyPath, sshKeyInstance.publicKeyPath, sshKeyInstance.sshPath = sshutils.CurrentSSHPath()
	sshKeyInstance.generatePublicKeyFromOldPrivateKey()
	hostsMap := ansible.ParseHostPattern(hostPattern, configFile)
	for _, hostInfo := range hostsMap {
		sshKeyInstance.Host = hostInfo.IP
		sshKeyInstance.User = user
		err := sshKeyInstance.delKey()
		if err != nil {
			logger.Failed(err.Error())
		}
	}
	logger.Success("已完成删除免密\n")
}

func (sshKey *SSHKey) generatePassword() string {
	password := make([]byte, passwordLength)
	rand.Read(password)
	for i := range password {
		password[i] = characters[int(password[i])%len(characters)]
	}
	return string(password)
}

func (sshKey *SSHKey) clearPasswordFile() error {
	err := os.WriteFile(passwordFile, []byte(""), 0600)
	if err != nil {
		return fmt.Errorf("failed to clear password file: %v", err)
	}
	return nil
}

func (sshKey *SSHKey) savePasswordToFile() error {
	data := []byte(fmt.Sprintf("ansible_host=%s ansible_user=%s ansible_ssh_pass=%s\n", sshKey.Host, sshKey.User, sshKey.NewPassword))

	file, err := os.OpenFile(passwordFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open password file: %v", err)
	}
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("failed to append password to file: %v", err)
	}

	return nil
}

func (sshKey *SSHKey) chpasswd() error {
	cmd := fmt.Sprintf("echo 'root:%s' | chpasswd", sshKey.NewPassword)
	if err := sshKey.runCMD(cmd); err != nil {
		return fmt.Errorf("failed to change user: %s password: %v", sshKey.User, err)
	}
	logger.Success("%v | User=%v | Status >> Success\n", sshKey.Host, sshKey.User)
	return nil
}

func Chpasswd(hostPattern, configFile, password string) {
	sshKeyInstance := newSSHKey()
	sshKeyInstance.privateKeyPath, sshKeyInstance.publicKeyPath, sshKeyInstance.sshPath = sshutils.CurrentSSHPath()
	sshKeyInstance.generatePublicKeyFromOldPrivateKey()
	hostsMap := ansible.ParseHostPattern(hostPattern, configFile)
	sshKeyInstance.clearPasswordFile()
	for _, hostInfo := range hostsMap {
		sshKeyInstance.Host = hostInfo.IP
		sshKeyInstance.User = hostInfo.User
		if password == "" {
			sshKeyInstance.NewPassword = sshKeyInstance.generatePassword()
		} else {
			sshKeyInstance.NewPassword = password
		}
		err := sshKeyInstance.chpasswd()
		if err != nil {
			logger.Failed(err.Error())
		}
		sshKeyInstance.savePasswordToFile()
	}
	logger.Success("新密码已保存在: %s\n", passwordFile)
}
