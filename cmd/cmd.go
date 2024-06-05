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
	setupCertCmd()
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
