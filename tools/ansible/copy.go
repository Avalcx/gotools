package ansible

import (
	"encoding/json"
	"fmt"
	"gotools/utils/logger"
	"os/exec"
)

type Copy struct {
	fileName     string
	src          string
	dest         string
	DestFullPath string `json:"dest"`
	Changed      bool   `json:"changed"`
	Status       string `json:"status"`
	Checksum     string `json:"checksum"`
	Msg          string `json:"msg"`
}

func (ansible *Ansible) copyOutput(status int) {
	// 0 success
	// 1 faild
	// 2 changed
	jsonOutput, err := json.MarshalIndent(ansible.Copy, "", "    ")
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

func (ansible *Ansible) isNeedCopy() (bool, error) {
	ansible.Command = "[[ -f " + ansible.Copy.DestFullPath + " ]]"
	_, exitStatus, err := ansible.runRemoteShell()
	// 排除test -f 返回值问题
	if err != nil && exitStatus == -99 {
		return false, err
	}
	// 不存在文件
	if exitStatus != 0 {
		return true, nil
	}

	ansible.Command = "md5sum " + ansible.Copy.DestFullPath
	remoteMD5, _, err := ansible.runRemoteShell()
	if err != nil {
		return false, err
	}

	localMD5, err := getLocalFileMd5(ansible.Copy.src)
	if err != nil {
		return false, err
	}
	//文件一致,无需发送文件
	if removeOtherString(remoteMD5) == removeOtherString(localMD5) {
		ansible.Copy.Checksum = removeOtherString(remoteMD5)
		return false, nil
	}

	return true, nil
}

func (ansible *Ansible) runSCP() error {
	cmd := exec.Command("scp", "-r", ansible.Copy.src, fmt.Sprintf("%s@%s:%s", ansible.HostInfo.User, ansible.HostInfo.IP, ansible.Copy.dest))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run scp command: %v, output: %s", err, output)
	}
	return nil
}

func (ansible *Ansible) execCopy() {
	isNeedCopy, err := ansible.isNeedCopy()
	if err != nil {
		ansible.Copy.Changed = false
		ansible.Copy.Status = "faild"
		ansible.Copy.Msg = err.Error()
		ansible.copyOutput(1)
		return
	}

	if !isNeedCopy {
		ansible.Copy.Changed = false
		ansible.Copy.Status = "success"
		ansible.copyOutput(0)
		return
	}
	err = ansible.runSCP()
	if err != nil {
		ansible.Copy.Changed = false
		ansible.Copy.Status = "faild"
		ansible.Copy.Msg = err.Error()
		ansible.copyOutput(1)
		return
	}

	ansible.Command = "md5sum " + ansible.Copy.DestFullPath
	remoteMD5, _, err := ansible.runRemoteShell()
	if err != nil {
		logger.Failed("%v", err)
		return
	}

	ansible.Copy.Checksum = removeOtherString(remoteMD5)
	ansible.Copy.Changed = true
	ansible.Copy.Status = "success"
	ansible.copyOutput(2)
}
