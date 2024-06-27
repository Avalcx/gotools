package ansible

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Ansible struct {
	SSH      SSH
	HostInfo HostInfo
	Command  string
	Copy     Copy
	Script   Script
}

type SSH struct {
	ClientConfig *ssh.ClientConfig
	Client       *ssh.Client
	Session      *ssh.Session
	SFTPClient   *sftp.Client
}

type HostInfo struct {
	Hostname   string
	IP         string
	Port       string
	User       string
	Password   string
	PrivateKey string
}

func NewAnsible() *Ansible {
	return &Ansible{
		HostInfo: *NewHostInfo(),
	}
}

func NewHostInfo() *HostInfo {
	return &HostInfo{
		Port: "22",
		User: "root",
	}
}

func (ansible *Ansible) newSSHClientConfig() error {
	var authMethods []ssh.AuthMethod

	key, err := os.ReadFile(ansible.HostInfo.PrivateKey)
	if err != nil {
		msg := fmt.Errorf("私钥文件读取错误: %v", err)
		return msg
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		msg := fmt.Errorf("解析私钥失败: %v", err)
		return msg
	}
	// 将私钥添加进入认证方案列表
	authMethods = append(authMethods, ssh.PublicKeys(signer))

	if ansible.HostInfo.Password != "" {
		authMethods = append(authMethods, ssh.Password(ansible.HostInfo.Password))
	}

	ansible.SSH.ClientConfig = &ssh.ClientConfig{
		User:            ansible.HostInfo.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
	return nil
}

func (ansible *Ansible) newSSHClient() error {
	var err error
	addr := fmt.Sprintf("%s:%s", ansible.HostInfo.IP, ansible.HostInfo.Port)
	ansible.SSH.Client, err = ssh.Dial("tcp", addr, ansible.SSH.ClientConfig)
	if err != nil {
		msg := fmt.Errorf("unable to connect: %s error %v", ansible.HostInfo.IP, err)
		return msg
	}
	return nil
}

func (ansible *Ansible) sshClientClose() error {
	err := ansible.SSH.Client.Close()
	return err
}

func (ansible *Ansible) newSFTPClient() error {
	var err error
	ansible.SSH.SFTPClient, err = sftp.NewClient(ansible.SSH.Client)
	if err != nil {
		return err
	}
	return nil
}

func (ansible *Ansible) sftpClientClose() error {
	err := ansible.SSH.SFTPClient.Close()
	return err
}

func (ansible *Ansible) newSSHSession() error {
	var err error
	ansible.SSH.Session, err = ansible.SSH.Client.NewSession()
	if err != nil {
		msg := fmt.Errorf("unable to connect: %s error %v", ansible.HostInfo.IP, err)
		return msg
	}
	return nil
}

func (ansible *Ansible) sshSessionClose() error {
	err := ansible.SSH.Session.Close()
	return err
}
