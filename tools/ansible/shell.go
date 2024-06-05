package ansible

import (
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

type ansible struct {
	PrivateKey string
	Config     *ssh.ClientConfig
	User       string
	Host       string
	Port       string
	Command    string
}

func (ansible *ansible) newConfig() (config *ssh.ClientConfig, err error) {
	key, err := os.ReadFile(ansible.PrivateKey)
	if err != nil {
		err = fmt.Errorf("unable to read private key: %v", err)
		return
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		err = fmt.Errorf("unable to parse private key: %v", err)
		return
	}

	config = &ssh.ClientConfig{
		User: ansible.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return
}

func (ansible *ansible) run() ([]byte, int, error) {
	addr := fmt.Sprintf("%s:22", ansible.Host)
	client, err := ssh.Dial("tcp", addr, ansible.Config)
	if err != nil {
		msg := fmt.Sprintf("unable to connect: %s error %v", ansible.Host, err)
		return nil, 99, fmt.Errorf("error: %v", msg)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		msg := fmt.Sprintf("ssh new session error %v", err)
		return nil, 99, fmt.Errorf("error: %v", msg)
	}
	defer session.Close()

	var exitStatus int
	output, err := session.CombinedOutput(ansible.Command)
	if err != nil {
		if exitErr, ok := err.(*ssh.ExitError); ok {
			exitStatus = exitErr.ExitStatus()
			msg := fmt.Sprintf(" %s\n", strings.TrimSpace(string(output)))
			return output, exitStatus, fmt.Errorf("error: %v", msg)
		}
		msg := fmt.Sprintf("run command %s on host %s error %v", ansible.Command, ansible.Host, err)
		return nil, exitStatus, fmt.Errorf("error: %v", msg)
	}
	return output, 0, nil
}

func ExecShell(host string, command string) {
	runShell := ansible{}
	runShell.PrivateKey = "/root/.ssh/id_rsa"
	runShell.User = "root"
	runShell.Host = host
	runShell.Command = command
	config, err := runShell.newConfig()
	if err != nil {
		log.Println("获取key失败")
		return
	}
	runShell.Config = config
	output, exitStatus, err := runShell.run()
	if err != nil {
		fmt.Printf("\033[31m%v | FAILED | rc=%v >>\n %v\033[0m\n", runShell.Host, exitStatus, err)
		return
	}
	// 输出正确期望值的日志
	fmt.Printf("\033[33m%v | CHANGED | rc=%v >>\n %s\033[0m\n", runShell.Host, exitStatus, output)
}
