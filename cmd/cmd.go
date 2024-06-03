package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ssltools",
	Short: "一个简单的ssl工具包",
}

func init() {
	setupCertCmd()
	setupPortCmd()
	setupSshkeyCmd()
	rootCmd.AddCommand(certCmd)
	rootCmd.AddCommand(portCmd)
	rootCmd.AddCommand(sshKeyCmd)
}

func Run() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
