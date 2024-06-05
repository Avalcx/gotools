package cmd

import (
	"gotools/tools/sshkey"

	"github.com/spf13/cobra"
)

var sshKeyCmd = &cobra.Command{
	Use:   "sshkey",
	Short: "ssh免密",
	Example: `
	gotools sshkey -h=192.168.1.1 -p={password} [default user is root]
	gotools sshkey -h=192.168.1.1-10 -u=zsops -p={password}
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && cmd.Flags().NFlag() == 0 {
			cmd.Help()
			return
		}

		hosts, _ := cmd.Flags().GetString("hosts")
		user, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")
		sshkey.PushKeys(hosts, user, password)
	},
}

func setupSshkeyCmd() {
	sshKeyCmd.Flags().StringP("hosts", "h", "", "ip(192.168.1.1 or 192.168.1.1-10)")
	sshKeyCmd.Flags().StringP("user", "u", "root", "username")
	sshKeyCmd.Flags().StringP("password", "p", "", "password")
}
