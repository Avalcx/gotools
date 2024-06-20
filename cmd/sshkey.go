package cmd

import (
	"gotools/tools/sshkey"

	"github.com/spf13/cobra"
)

var sshKeyCmd = &cobra.Command{
	Use:   "sshkey",
	Short: "ssh免密",
	Example: `
	gotools sshkey -p={password} -h=192.168.1.1
	gotools sshkey -p={password} -u=zsops -h=192.168.1.1-10 
	gotools sshkey -p={password} -u=zsops -h=192.168.1.1 -h 192.168.1.10 
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
