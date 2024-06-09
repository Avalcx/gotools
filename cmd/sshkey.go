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

		hostsSlice, _ := cmd.Flags().GetStringSlice("hosts")
		user, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")
		sshkey.PushKeys(hostsSlice, user, password)
	},
}

func setupSSHkeyCmd() {
	sshKeyCmd.Flags().StringSliceP("hosts", "h", nil, "ip地址或地址段(192.168.1.1 or 192.168.1.1-10),可以指定多个-h")
	sshKeyCmd.Flags().StringP("user", "u", "root", "username")
	sshKeyCmd.Flags().StringP("password", "p", "", "password")
	sshKeyCmd.MarkFlagRequired("hosts")
	sshKeyCmd.MarkFlagRequired("password")
}
