package ansible

import (
	"fmt"
	"gotools/utils/logger"
	"strings"

	"golang.org/x/crypto/ssh"
)

func (ansible *Ansible) runRemoteShell() (string, int, error) {
	var err error
	var rc int
	err = ansible.newSSHClientConfig()
	if err != nil {
		return "", -99, err
	}

	err = ansible.newSSHClient()
	if err != nil {
		return "", -99, err
	}
	defer ansible.sshClientClose()

	err = ansible.newSSHSession()
	if err != nil {
		return "", -99, err
	}
	defer ansible.sshSessionClose()

	output, err := ansible.SSH.Session.CombinedOutput(ansible.Command)

	if err != nil {
		if exitErr, ok := err.(*ssh.ExitError); ok {
			rc = exitErr.ExitStatus()
			msg := fmt.Sprintf(" %s\n", strings.TrimSpace(string(output)))
			return string(output), rc, fmt.Errorf("%v", msg)
		}
		msg := fmt.Sprintf("run command %s on host %s error %v", ansible.Command, ansible.HostInfo.IP, err)
		return "", rc, fmt.Errorf("error: %v", msg)
	}
	return string(output), 0, nil
}

func (ansible *Ansible) runShellModule() {
	output, rc, err := ansible.runRemoteShell()
	if err != nil {
		logger.Failed("%v | FAILED | rc=%v >>\n %v", ansible.HostInfo.IP, rc, err)
		return
	}
	logger.Changed("%v | CHANGED | rc=%v >>\n %s", ansible.HostInfo.IP, rc, output)
}
