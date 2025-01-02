package cmd

import (
	"github.com/Avalcx/gotools/tools/ansible"

	"github.com/spf13/cobra"
)

var ansibleArgs ansible.AnsibleArgs

var ansibleCmd = &cobra.Command{
	Use:   "ansible",
	Short: "ansible",
	Long:  "Host-Pattern为ip地址或地址段时,会直接执行,为组名时,会读取配置文件中的配置组,默认为/etc/ansible/hosts",
	Args:  cobra.ExactArgs(1),
	Example: `
	gotools ansible [HOST-PATTERN] -m [MODULE_NAME] -a [ARGS]
	gotools ansible 192.168.1.1 -m shell -a "w"
	gotools ansible 192.168.1.1-10 -m copy -a "src=xxx dest=xxx"
	gotools ansible group1 -m script -a "/path/to/script.sh"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		hostPattern := args[0]
		ansible.RunModules(hostPattern, ansibleArgs)
	},
}

func setupAnsibleCmd() {
	ansibleCmd.Flags().StringVarP(&ansibleArgs.ModuleName, "module-name", "m", "", "--module-name [ shell | copy | script ]")
	ansibleCmd.Flags().StringVarP(&ansibleArgs.Args, "args", "a", "", "--args")
	ansibleCmd.Flags().StringVar(&ansibleArgs.Password, "password", "", "当PrivateKey无效时,使用密码登录")
	ansibleCmd.Flags().StringVar(&ansibleArgs.ConfigFile, "config-file", "/etc/ansible/hosts", "兼容ansible-hosts文件,默认读取/etc/ansible/hosts,手动指定其他路径")
}
