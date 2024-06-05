package cmd

import (
	"gotools/tools/port"
	"gotools/utils/logger"

	"github.com/spf13/cobra"
)

var portCmd = &cobra.Command{
	Use:   "port",
	Short: "端口工具",
	Long:  "用于测试网络策略或防火墙策略",
	Example: `
	gotools port --server --port=80,443,8080-8099
	gotools port --client --port=80,443,8080-8099 --host=127.0.0.1
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && cmd.Flags().NFlag() == 0 {
			cmd.Help()
			return
		}

		client, _ := cmd.Flags().GetBool("client")
		server, _ := cmd.Flags().GetBool("server")
		if (client && server) || (!server && !client) {
			logger.Failed("--client或--server 至少要指定一个")
			return
		}

		portSpecs, _ := cmd.Flags().GetString("ports")
		mode, _ := cmd.Flags().GetString("mode")
		host, _ := cmd.Flags().GetString("host")
		if server {
			port.StartServer(portSpecs, mode)
		} else if client {
			port.StartClient(portSpecs, host, mode)
		}
	},
}

func setupPortCmd() {
	portCmd.Flags().Bool("client", false, "client")
	portCmd.Flags().Bool("server", false, "server")
	portCmd.Flags().StringP("mode", "m", "tcp", "tcp")
	portCmd.Flags().StringP("ports", "p", "", "listen port")
	portCmd.Flags().StringP("host", "h", "127.0.0.1", "target ip")
}
