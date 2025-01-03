package cmd

import (
	"github.com/Avalcx/gotools/tools/sshkey"

	"github.com/spf13/cobra"
)

var sshKeyCmd = &cobra.Command{
	Use:   "sshkey",
	Short: "ssh免密",
	Long:  "批量免密,批量删除免密,批量改密码",
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
		isChpasswd, _ := cmd.Flags().GetBool("chpasswd")
		if isDelete {
			sshkey.DelKeys(hostPattern, configFile, user)
		} else if isChpasswd {
			sshkey.Chpasswd(hostPattern, configFile, password)
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
	sshKeyCmd.Flags().Bool("chpasswd", false, "修改指定主机的密码,如果没有指定-p参数,将使用随机密码")
}
