package cmd

import (
	"gotools/tools/sshkey"

	"github.com/spf13/cobra"
)

var sshKeyCmd = &cobra.Command{
	Use:   "sshkey",
	Short: "生成并传输sshkey",
	Example: `
	gotools sshkey --hosts=192.168.1.1 --password={password} [default user is root]
	gotools sshkey --hosts=192.168.1.1-10 --user=zsops --password={password}
	`,
	Run: func(cmd *cobra.Command, args []string) {
		hosts, _ := cmd.Flags().GetString("hosts")
		user, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")
		sshkey.PushKeys(hosts, user, password)
	},
}

func setupSshkeyCmd() {
	sshKeyCmd.Flags().StringP("hosts", "H", "", "ip(192.168.1.1 or 192.168.1.1-10)")
	sshKeyCmd.Flags().StringP("user", "u", "root", "username")
	sshKeyCmd.Flags().StringP("password", "p", "", "password")
}
