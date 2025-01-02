package sshutils

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"os/user"
	"path/filepath"

	"github.com/Avalcx/gotools/utils/logger"
)

func CurrentSSHPath() (privateKeyPath, publicKeyPath, sshPath string) {
	currentUser, err := user.Current()
	if err != nil {
		logger.Fatal("failed to get current user: %v\n", err)
	}
	sshDir := filepath.Join(currentUser.HomeDir, ".ssh")
	privateKeyPath = filepath.Join(sshDir, "id_rsa")
	publicKeyPath = filepath.Join(sshDir, "id_rsa.pub")
	return privateKeyPath, publicKeyPath, sshDir
}

func GetLocalFileMd5(file string) (string, error) {
	fileByte, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	m := md5.New()
	m.Write([]byte(fileByte))
	return hex.EncodeToString(m.Sum(nil)), nil
}

func FileIsExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
