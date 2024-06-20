package cmd

import (
	"gotools/tools/sshkey"

	"github.com/spf13/cobra"
)

var sshKeyCmd = &cobra.Command{
	Use:   "sshkey",
	Short: "ssh免密",
	Example: `
	gotools sshkey 192.168.1.1-10 -p={password}
	gotools sshkey group1 -p={password}
	gotools sshkey [HOST-PATTERN] -p={password} -u=zsops --config-file=/path/to/file 
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && cmd.Flags().NFlag() == 0 {
			cmd.Help()
			return
		}
		hostPattern := args[0]
		configFile, _ := cmd.Flags().GetString("config-file")
		user, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")
		isDelete, _ := cmd.Flags().GetBool("delete")
		if isDelete {
			sshkey.DelKeys(hostPattern, configFile, user)
		} else {
			sshkey.PushKeys(hostPattern, configFile, user, password)
		}
	},
}

func setupSSHkeyCmd() {
	sshKeyCmd.Flags().StringP("user", "u", "root", "username")
	sshKeyCmd.Flags().StringP("password", "p", "", "password")
	sshKeyCmd.Flags().String("config-file", "/etc/ansible/hosts", "兼容ansible-hosts文件,默认读取/etc/ansible/hosts,手动指定其他路径")
	sshKeyCmd.Flags().Bool("delete", false, "删除指定主机的sshkey")
}
