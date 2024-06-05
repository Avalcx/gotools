package cmd

import (
	"gotools/tools/ansible"

	"github.com/spf13/cobra"
)

var ansibleCmd = &cobra.Command{
	Use:   "ansible",
	Short: "ansible",
	Example: `
	gotools ansible shell
	gotools ansible script
	`,
}

func setupAnsibleCmd() {
	setupShellCmd()
}

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "shell模块",
	Example: `
	gotools ansible shell -h 192.168.1.1 -c [shell]
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && cmd.Flags().NFlag() == 0 {
			cmd.Help()
			return
		}

		host, _ := cmd.Flags().GetString("host")
		command, _ := cmd.Flags().GetString("cmd")
		ansible.ExecShell(host, command)
	},
}

func setupShellCmd() {
	ansibleCmd.AddCommand(shellCmd)
	shellCmd.Flags().StringP("host", "h", "", "host")
	shellCmd.Flags().StringP("cmd", "c", "", "cmd")
}
