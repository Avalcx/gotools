package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gotools",
	Short: "一个简单的gotools工具包",
}

func init() {
	help()
	setupSSLCmd()
	setupPortCmd()
	setupSshkeyCmd()
	setupAnsibleCmd()
	rootCmd.AddCommand(certCmd)
	rootCmd.AddCommand(portCmd)
	rootCmd.AddCommand(sshKeyCmd)
	rootCmd.AddCommand(ansibleCmd)
}

func Run() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func help() {
	//去除默认的--help的-h 的flag
	rootCmd.PersistentFlags().BoolP("help", "", false, "Help for this command")

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		help, err := cmd.Flags().GetBool("help")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if help {
			cmd.Help()
			os.Exit(0)
		}
	}
}
