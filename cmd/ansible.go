package cmd

import (
	"gotools/tools/ansible"

	"github.com/spf13/cobra"
)

var ansibleCmd = &cobra.Command{
	Use:   "ansible",
	Short: "ansible",
	Example: `
	gotools ansible shell -H 192.168.1.1 -c w
	`,
}

func setupAnsibleCmd() {
	setupShellCmd()
}

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "shell模块",
	Example: `
	gotools ansible shell -H 192.168.1.1 -c w
	`,
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		command, _ := cmd.Flags().GetString("args")
		ansible.ExecShell(host, command)
	},
}

func setupShellCmd() {
	ansibleCmd.AddCommand(shellCmd)
	shellCmd.Flags().StringP("host", "H", "", "host")
	shellCmd.Flags().StringP("args", "a", "", "args")
}
