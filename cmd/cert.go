package cmd

import (
	"gotools/tools/cert"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var certCmd = &cobra.Command{
	Use:   "cert",
	Short: "证书工具",
	Long:  "用于生成自签名证书、查看证书的有效期",
	Example: `
	gotools cert privite -d=zsops.cn -y=10
	gotools cert privite -d=domain1.cn -d=domain2.cn -y=10
	gotools cert acme -d=zsops.cn -a={AliAK} -s={AliSK}
	gotools cert check -d=baidu.com
	gotools cert check -f=cert.pem
	`,
}

func setupCertCmd() {
	setupCertCheckCmd()
	setupPriviteCertCmd()
	setupAcmeCertCmd()
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "查看证书有效期",
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")
		port, _ := cmd.Flags().GetString("port")
		file, _ := cmd.Flags().GetString("file")

		if domain != "" && file != "" {
			log.Fatal("domain和file 不能同时有参数")
		}

		if domain != "" {
			checkFromDomain(domain, port)
		}

		if file != "" {
			checkFromFile(file)
		}
	},
}

func setupCertCheckCmd() {
	certCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringP("domain", "d", "", "通过域名查询")
	checkCmd.Flags().StringP("port", "p", "443", "通过域名查询的端口")
	checkCmd.Flags().StringP("file", "f", "", "通过文件查询")
}

func checkFromDomain(domain string, port string) {
	cert.CheckFromDomain(domain, port)
}

func checkFromFile(file string) {
	cert.CheckFromCrtFile(file)
}

var priviteCmd = &cobra.Command{
	Use:   "privite",
	Short: "生成私有证书",
	Run: func(cmd *cobra.Command, args []string) {
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
	cert.GeneratePrivateCert(domain, years)
}

var acmeCmd = &cobra.Command{
	Use:   "acme",
	Short: "生成Let's Encrypt证书,only by aliDNS",
	Run: func(cmd *cobra.Command, args []string) {
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
			log.Fatal("ACCESSKEY or SECRETKEY is empty")
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
	cert.Acme(domainList, ak, sk)
}
