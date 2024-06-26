package cmd

import (
	"gotools/tools/port"

	"github.com/spf13/cobra"
)

var portCmd = &cobra.Command{
	Use:   "port",
	Short: "端口工具",
	Long:  "用于测试网络策略或防火墙策略",
	Example: `
	gotools port server --port=80,443,8080-8099
	gotools port client --port=80,443,8080-8099 --host=127.0.0.1
	`,
}

func setupPortCmd() {
	setupPortServerCmd()
	setupPortClientCmd()
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "server mode",
	Run: func(cmd *cobra.Command, args []string) {
		portSpecs, _ := cmd.Flags().GetString("ports")
		mode, _ := cmd.Flags().GetString("mode")
		startPortServer(portSpecs, mode)
	},
}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "client mode",
	Run: func(cmd *cobra.Command, args []string) {
		portSpecs, _ := cmd.Flags().GetString("ports")
		host, _ := cmd.Flags().GetString("host")
		mode, _ := cmd.Flags().GetString("mode")
		startPortClient(portSpecs, host, mode)
	},
}

func setupPortServerCmd() {
	portCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringP("ports", "p", "", "listen port")
	serverCmd.Flags().StringP("mode", "m", "", "tcp or udp")
}

func setupPortClientCmd() {
	portCmd.AddCommand(clientCmd)
	clientCmd.Flags().StringP("ports", "p", "", "test port")
	clientCmd.Flags().StringP("host", "H", "127.0.0.1", "test host")
	clientCmd.Flags().StringP("mode", "m", "", "tcp or udp")
}

func startPortServer(portSpecs string, mode string) {
	port.StartServer(portSpecs, mode)
}

func startPortClient(portSpecs string, host string, mode string) {
	port.StartClient(portSpecs, host, mode)
}
