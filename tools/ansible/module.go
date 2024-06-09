package ansible

import (
	"gotools/utils/logger"
	"strings"
)

type AnsibleArgs struct {
	Password   string
	Args       string
	ModuleName string
	ConfigFile string
}

func parseHostPattern(hostPattern, configFile string) []*Ansible {
	ansibleInstanceList := make([]*Ansible, 0, 10)

	// IP地址段
	if isIPRange(hostPattern) {
		ipRange, err := ParseIPRange(hostPattern)
		if err != nil {
			logger.Fatal("%v", err)
		}
		for _, ip := range ipRange {
			ansibleInstance := NewAnsible()
			ansibleInstance.HostInfo.IP = ip.String()
			ansibleInstanceList = append(ansibleInstanceList, ansibleInstance)
		}
		return ansibleInstanceList
		// IP
	} else if isIPAddress(hostPattern) {
		ansibleInstance := NewAnsible()
		ansibleInstance.HostInfo.IP = hostPattern
		ansibleInstanceList = append(ansibleInstanceList, ansibleInstance)
		return ansibleInstanceList
		// 组名
	} else {
		ansibleInstanceMap, err := parseGroupFromFile(configFile)
		if err != nil {
			logger.Fatal("%v", err)
		}
		for group, infos := range ansibleInstanceMap {
			if group == hostPattern {
				ansibleInstanceList = append(ansibleInstanceList, infos...)
			}
		}
		return ansibleInstanceList
	}
}

func RunModules(hostPattern string, ansibleArgs AnsibleArgs) {
	ansibleMap := parseHostPattern(hostPattern, ansibleArgs.ConfigFile)
	for _, ansible := range ansibleMap {
		ansible.HostInfo.PrivateKey, _ = currentSSHPath()
		if ansibleArgs.Password != "" {
			ansible.HostInfo.Password = ansibleArgs.Password
		}
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
