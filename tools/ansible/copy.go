package ansible

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/Avalcx/gotools/utils/logger"
	"github.com/Avalcx/gotools/utils/sshutils"
)

type Copy struct {
	fileName     string // 文件名
	src          string // 源路径，包含文件名
	dest         string // 目的路径，不包含文件名
	DestFullPath string `json:"dest"` // 目的路径，包含文件名
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
	_, rc, err := ansible.runRemoteShell()
	// 排除test -f 返回值问题
	if err != nil && rc == -99 {
		return false, err
	}
	// 不存在文件
	if rc != 0 {
		return true, nil
	}

	ansible.Command = "md5sum " + ansible.Copy.DestFullPath
	remoteMD5, _, err := ansible.runRemoteShell()
	if err != nil {
		return false, err
	}

	localMD5, err := sshutils.GetLocalFileMd5(ansible.Copy.src)
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

type ProgressWriter struct {
	Total       int64
	Transferred int64
	StartTime   time.Time
	LastTime    time.Time
}

func (pw *ProgressWriter) Update(n int64) {
	pw.Transferred += n
	now := time.Now()
	percentage := float64(pw.Transferred) / float64(pw.Total) * 100
	elapsed := now.Sub(pw.LastTime).Seconds()
	speed := float64(n) / 1024 / elapsed
	pw.LastTime = now
	fmt.Printf("\rProgress: %.2f%%, Speed: %.2f KB/s", percentage, speed)
}

// func (ansible *Ansible) runSCP() error {
// 	cmd := exec.Command("scp", "-r", ansible.Copy.src, fmt.Sprintf("%s@%s:%s", ansible.HostInfo.User, ansible.HostInfo.IP, ansible.Copy.dest))
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("failed to run scp command: %v, output: %s", err, output)
// 	}
// 	return nil
// }

func (ansible *Ansible) runSFTP() error {
	err := ansible.newSSHClient()
	if err != nil {
		return err
	}
	defer ansible.sshClientClose()

	err = ansible.newSFTPClient()
	if err != nil {
		return err
	}
	defer ansible.sftpClientClose()

	srcFile, err := os.Open(ansible.Copy.src)
	if err != nil {
		return fmt.Errorf("failed to open local file: %v", err)
	}
	defer srcFile.Close()

	srcFileInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	dstFile, err := ansible.SSH.SFTPClient.Create(ansible.Copy.DestFullPath)
	if err != nil {
		return fmt.Errorf("failed to create remote file: %v", err)
	}
	defer dstFile.Close()

	// 创建进度写入器
	progressWriter := &ProgressWriter{
		Total:     srcFileInfo.Size(),
		StartTime: time.Now(),
		LastTime:  time.Now(),
	}

	// 使用较大的缓冲区提高传输效率
	const bufferSize = 32 * 1024 // 32KB
	buf := make([]byte, bufferSize)

	// 复制本地文件内容到远程文件
	for {
		n, err := srcFile.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read from source file: %v", err)
		}
		if n == 0 {
			break
		}

		// 写入远程文件
		if _, err := dstFile.Write(buf[:n]); err != nil {
			return fmt.Errorf("failed to write to destination file: %v", err)
		}

		progressWriter.Update(int64(n))
	}

	fmt.Println("\nFile transfer completed.")
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
	// err = ansible.runSCP()
	err = ansible.runSFTP()
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

func removeOtherString(input string) string {
	noFileName := strings.Split(string(input), " ")[0]
	noSpaces := strings.ReplaceAll(noFileName, " ", "")
	noNewlines := strings.ReplaceAll(noSpaces, "\n", "")
	return noNewlines
}
