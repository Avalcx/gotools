package cmd

import (
	"gotools/tools/port"
	"gotools/utils/logger"

	"github.com/spf13/cobra"
)

var portInfo port.Port

var portCmd = &cobra.Command{
	Use:   "port",
	Short: "端口工具",
	Long:  "用于测试网络策略或防火墙策略",
	Example: `
	gotools port --server --ports=80,443,8080-8099 --protocol tcp
	gotools port --client --ports=80,443,8080-8099 --host=127.0.0.1 --protocol udp
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && cmd.Flags().NFlag() == 0 {
			cmd.Help()
			return
		}

		client, _ := cmd.Flags().GetBool("client")
		server, _ := cmd.Flags().GetBool("server")
		if (client && server) || (!server && !client) {
			logger.Failed("--client或--server 至少且只能要指定一个")
			return
		}
		portInfo.ParsePortSpecs()
		if server {
			portInfo.StartServer()
		} else if client {
			portInfo.StartClient()
		}
	},
}

func setupPortCmd() {
	portCmd.Flags().Bool("client", false, "client模式")
	portCmd.Flags().Bool("server", false, "server模式")
	portCmd.Flags().StringVar(&portInfo.Protocol, "protocol", "tcp", "protocol: tcp or udp")
	portCmd.Flags().StringVarP(&portInfo.PortSpecs, "ports", "p", "", "port")
	portCmd.Flags().StringVarP(&portInfo.Host, "host", "h", "127.0.0.1", "target ip")
	portCmd.MarkFlagRequired("ports")
}
