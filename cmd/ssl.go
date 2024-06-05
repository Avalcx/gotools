package cmd

import (
	"gotools/tools/ssl"
	"gotools/utils/logger"
	"os"

	"github.com/spf13/cobra"
)

var certCmd = &cobra.Command{
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

		domain, _ := cmd.Flags().GetString("domain")
		port, _ := cmd.Flags().GetString("port")
		file, _ := cmd.Flags().GetString("file")

		if domain != "" && file != "" {
			logger.Fatal("domain和file 不能同时有参数")
		}

		if domain != "" {
			checkFromDomain(domain, port)
		}

		if file != "" {
			checkFromFile(file)
		}
	},
}

func setupSSLCheckCmd() {
	certCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringP("domain", "d", "", "通过域名查询")
	checkCmd.Flags().StringP("port", "p", "443", "通过域名查询的端口")
	checkCmd.Flags().StringP("file", "f", "", "通过文件查询")
}

func checkFromDomain(domain string, port string) {
	ssl.CheckFromDomain(domain, port)
}

func checkFromFile(file string) {
	ssl.CheckFromCrtFile(file)
}

var priviteCmd = &cobra.Command{
	Use:   "privite",
	Short: "生成自签证书",
	Example: `
	gotools ssl privite -d=zsops.cn -y=10
	gotools ssl privite -d=domain1.cn -d=domain2.cn -y=10
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && cmd.Flags().NFlag() == 0 {
			cmd.Help()
			return
		}

		years, _ := cmd.Flags().GetInt("years")
		domain, _ := cmd.Flags().GetStringSlice("domain")
		createPriviteCert(domain, years)
	},
}

func setupPriviteCertCmd() {
	certCmd.AddCommand(priviteCmd)
	priviteCmd.Flags().IntP("years", "y", 10, "有效期")
	priviteCmd.Flags().StringSliceP("domain", "d", nil, "domian list")
}

func createPriviteCert(domain []string, years int) {
	ssl.GeneratePrivateCert(domain, years)
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

		domainList, _ := cmd.Flags().GetStringSlice("domain")
		ak, _ := cmd.Flags().GetString("accesskey")
		if ak == "" {
			ak = os.Getenv("ALI_ACCESSKEY")
		}
		sk, _ := cmd.Flags().GetString("secretkey")
		if sk == "" {
			sk = os.Getenv("ALI_SECRETKEY")
		}
		if ak == "" || sk == "" {
			logger.Fatal("ACCESSKEY or SECRETKEY is empty")
		}
		cmd.Flags()
		createAcmeCert(domainList, ak, sk)
	},
}

func setupAcmeCertCmd() {
	certCmd.AddCommand(acmeCmd)
	acmeCmd.Flags().StringSliceP("domain", "d", nil, "domain")
	acmeCmd.Flags().StringP("accesskey", "a", "", "accessKey,default from env")
	acmeCmd.Flags().StringP("secretkey", "s", "", "secretkey,default from env")
}

func createAcmeCert(domainList []string, ak string, sk string) {
	ssl.Acme(domainList, ak, sk)
}
