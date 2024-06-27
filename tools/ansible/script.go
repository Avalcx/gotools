package ansible

import (
	"encoding/json"
	"fmt"
	"gotools/utils/logger"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Script struct {
	localScriptPath  string
	remoteScriptPath string
	fileName         string
	remoteScriptDir  string
	Changed          bool   `json:"changed"`
	Rc               int    `json:"rc"`
	Stdout           string `json:"stdout"`
	Stderr           string `json:"stderr"`
}

func (ansible *Ansible) scriptOutput(status int) {
	// 0 success
	// 1 faild
	// 2 changed
	jsonOutput, err := json.MarshalIndent(ansible.Script, "", "    ")
	if err != nil {
		logger.Failed("Failed to marshal JSON: %v\n", err)
		return
	}
	switch status {
	case 0:
		logger.Success("%v | SUCCESS => %s\n", ansible.HostInfo.IP, string(jsonOutput))
	case 1:
		logger.Failed("%v | FAILD => %v\n", ansible.HostInfo.IP, string(jsonOutput))
	case 2:
		logger.Changed("%v | CHANGED => %s\n", ansible.HostInfo.IP, string(jsonOutput))
	}
}

func (ansible *Ansible) runScriptModule() {
	output, rc, err := ansible.runScriptOverSSH()
	if err != nil {
		ansible.Script.Rc = rc
		ansible.Script.Changed = false
		ansible.Script.Stderr = err.Error()
		ansible.scriptOutput(1)
		return
	}
	ansible.Script.Rc = rc
	ansible.Script.Changed = true
	ansible.Script.Stdout = output
	ansible.scriptOutput(2)
}

func (ansible *Ansible) runScriptOverSSH() (string, int, error) {

	err := ansible.newSSHClientConfig()
	if err != nil {
		return "", -99, err
	}

	err = ansible.newSSHClient()
	if err != nil {
		return "", -99, err
	}
	defer ansible.sshClientClose()

	err = ansible.newSFTPClient()
	if err != nil {
		return "", -99, err
	}
	defer ansible.sftpClientClose()

	if err := ansible.uploadScript(); err != nil {
		return "", -99, fmt.Errorf("failed to upload script: %v", err)
	}

	output, rc, err := ansible.execRemoteScript()
	if err != nil {
		return "", -99, fmt.Errorf("failed to execute script: %v", err)
	}

	err = ansible.removeScript()
	if err != nil {
		return "", -99, err
	}

	return output, rc, nil
}

func (ansible *Ansible) uploadScript() error {
	srcFile, err := os.Open(ansible.Script.localScriptPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	ansible.generaRemoteScriptPath()
	err = ansible.SSH.SFTPClient.MkdirAll(ansible.Script.remoteScriptDir)
	if err != nil {
		return err
	}
	dstFile, err := ansible.SSH.SFTPClient.Create(ansible.Script.remoteScriptPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := dstFile.ReadFrom(srcFile); err != nil {
		return err
	}
	return nil
}

func (ansible *Ansible) removeScript() error {
	walker := ansible.SSH.SFTPClient.Walk(ansible.Script.remoteScriptDir)
	for walker.Step() {
		if err := walker.Err(); err != nil {
			return err
		}

		if walker.Stat().IsDir() {
			continue
		}

		if err := ansible.SSH.SFTPClient.Remove(walker.Path()); err != nil {
			return err
		}
	}

	if err := ansible.SSH.SFTPClient.RemoveDirectory(ansible.Script.remoteScriptDir); err != nil {
		return err
	}
	return nil
}

func (ansible *Ansible) execRemoteScript() (string, int, error) {
	ansible.Command = fmt.Sprintf("bash %s", ansible.Script.remoteScriptPath)
	output, rc, err := ansible.runRemoteShell()
	return string(output), rc, err
}

func (ansible *Ansible) generaRemoteScriptPath() {
	if strings.Contains(ansible.Script.localScriptPath, "/") {
		parts := strings.Split(ansible.Script.localScriptPath, "/")
		ansible.Script.fileName = parts[len(parts)-1]
	} else {
		ansible.Script.fileName = ansible.Script.localScriptPath
	}
	ansible.generateRandomAnsiblePath()
	ansible.Script.remoteScriptPath = ansible.Script.remoteScriptDir + "/" + ansible.Script.fileName
}

func (ansible *Ansible) generateRandomAnsiblePath() {
	timestamp := fmt.Sprintf("%.2f", float64(time.Now().Unix())+rand.Float64())
	randomNum1 := rand.Intn(10000)
	randomNum2 := rand.Int63()
	ansible.Script.remoteScriptDir = fmt.Sprintf("/tmp/.ansible/ansible-tmp-%s-%d-%d", timestamp, randomNum1, randomNum2)
}
