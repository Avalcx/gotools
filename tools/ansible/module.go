package ansible

import (
	"gotools/utils/logger"
	"gotools/utils/sshutils"
	"strings"
)

type AnsibleArgs struct {
	Password   string
	Args       string
	ModuleName string
	ConfigFile string
}

func RunModules(hostPattern string, ansibleArgs AnsibleArgs) {
	ansible := NewAnsible()
	hostsMap := ParseHostPattern(hostPattern, ansibleArgs.ConfigFile)
	for _, hostInfo := range hostsMap {
		hostInfo.PrivateKey, _ = sshutils.CurrentSSHPath()
		if ansibleArgs.Password != "" {
			hostInfo.Password = ansibleArgs.Password
		}
		ansible.HostInfo = *hostInfo
		switch ansibleArgs.ModuleName {
		case "shell":
			ansible.Command = ansibleArgs.Args
			ansible.runShellModule()
		case "copy":
			ansible.parseCopyModuleArgs(ansibleArgs.Args)
			ansible.execCopy()
		case "script":
			ansible.Script.localScriptPath = ansibleArgs.Args
			ansible.runScriptModule()
		default:
			logger.Failed("模块错误\n")
		}
	}
}

func (ansible *Ansible) parseCopyModuleArgs(args string) {
	result := make(map[string]string)
	pairs := strings.Split(args, " ")
	for _, pair := range pairs {
		keyValue := strings.SplitN(pair, "=", 2)
		if len(keyValue) != 2 {
			logger.Fatal("参数格式错误: src=xxx dest=xxx")
		}
		key := strings.TrimSpace(keyValue[0])
		value := strings.TrimSpace(keyValue[1])

		if key != "src" && key != "dest" {
			logger.Fatal("错误的参数key")
		}

		if _, exists := result[key]; exists {
			logger.Fatal("错误的参数")
		}

		result[key] = value
	}

	if _, srcExists := result["src"]; !srcExists {
		logger.Fatal("missing key: src")
	}
	if _, destExists := result["dest"]; !destExists {
		logger.Fatal("missing key: dest")
	}
	if len(result) != 2 {
		logger.Fatal("unexpected number of keys")
	}
	ansible.Copy.src = result["src"]
	ansible.Copy.dest = result["dest"]
	if strings.Contains(ansible.Copy.src, "/") {
		parts := strings.Split(ansible.Copy.src, "/")
		ansible.Copy.fileName = parts[len(parts)-1]
	} else {
		ansible.Copy.fileName = ansible.Copy.src
	}
	ansible.Copy.DestFullPath = ansible.Copy.dest + "/" + ansible.Copy.fileName
}
