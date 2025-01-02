package cmd

import (
	"github.com/Avalcx/gotools/tools/ssl"

	"github.com/spf13/cobra"
)

var sslInfo ssl.SSLInfo

var sslCmd = &cobra.Command{
	Use:   "ssl",
	Short: "证书工具",
	Example: `
	gotools ssl check
	gotools ssl privite
	gotools ssl acme
	`,
}

func setupSSLCmd() {
	setupSSLCheckCmd()
	setupPriviteCertCmd()
	setupAcmeCertCmd()
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "通过域名或证书文件检查证书信息",
	Example: `
	gotools ssl check -d=baidu.com
	gotools ssl check -f=cert.pem
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && cmd.Flags().NFlag() == 0 {
			cmd.Help()
			return
		}
		sslInfo.Check()
	},
}

func setupSSLCheckCmd() {
	sslCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringSliceVarP(&sslInfo.Domains, "domain", "d", nil, "通过域名查询,可以指定多个-d")
	checkCmd.Flags().StringVar(&sslInfo.Port, "port", "443", "--port 指定域名的端口")
	checkCmd.Flags().StringVarP(&sslInfo.CertFile, "file", "f", "", "通过文件查询")
}

var priviteCmd = &cobra.Command{
	Use:   "privite",
	Short: "生成自签证书",
	Example: `
	gotools ssl privite -d=zsops.cn
	gotools ssl privite -d=domain1.cn -d=domain2.cn -y=10
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && cmd.Flags().NFlag() == 0 {
			cmd.Help()
			return
		}
		sslInfo.GeneratePrivateCert()
	},
}

func setupPriviteCertCmd() {
	sslCmd.AddCommand(priviteCmd)
	priviteCmd.Flags().StringSliceVarP(&sslInfo.Domains, "domain", "d", nil, "通过域名查询,可以指定多个-d")
	priviteCmd.Flags().IntVarP(&sslInfo.Years, "years", "y", 10, "有效期")
}

var acmeCmd = &cobra.Command{
	Use:   "acme",
	Short: "生成Let's Encrypt证书,only by aliDNS",
	Example: `
	gotools ssl acme -d=zsops.cn -a={AliAK} -s={AliSK}
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && cmd.Flags().NFlag() == 0 {
			cmd.Help()
			return
		}
		sslInfo.Acme()
	},
}

func setupAcmeCertCmd() {
	sslCmd.AddCommand(acmeCmd)
	acmeCmd.Flags().StringSliceVarP(&sslInfo.Domains, "domain", "d", nil, "通过域名查询,可以指定多个-d")
	acmeCmd.Flags().StringVarP(&sslInfo.AliAK, "accesskey", "a", "", "accessKey")
	acmeCmd.Flags().StringVarP(&sslInfo.AliSK, "secretkey", "s", "", "secretkey")
}
