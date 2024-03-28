package cmd

import (
	"gotools/tools/cert"

	"github.com/spf13/cobra"
)

var certCmd = &cobra.Command{
	Use:   "cert",
	Short: "证书工具",
	Long:  "用于生成自签名证书、查看证书的有效期",
	Example: `
	gotools cert create --domain=zsops.cn --years=10
	gotools cert check --domain=baidu.com
	gotools cert check --file=cert.pem
	`,
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "查看证书有效期",
	Run: func(cmd *cobra.Command, args []string) {
		domain, _ := cmd.Flags().GetString("domain")

		if domain != "" {
			checkFromDomain(domain)
		}

		file, _ := cmd.Flags().GetString("file")
		if file != "" {
			checkFromFile(file)
		}
	},
}

func stepCertCmd() {
	stepCertCheckCmd()
	stepCertCreateCmd()
}

func stepCertCheckCmd() {
	certCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringP("domain", "d", "", "通过域名查询")
	checkCmd.Flags().StringP("file", "f", "", "通过文件查询")
}

func checkFromDomain(domain string) {
	cert.CheckFromDomain(domain)
}

func checkFromFile(file string) {
	cert.CheckFromCrtFile(file)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "生成私有证书",
	Run: func(cmd *cobra.Command, args []string) {
		years, _ := cmd.Flags().GetInt("years")
		domain, _ := cmd.Flags().GetStringSlice("domain")
		createPriviteCert(domain, years)
	},
}

func stepCertCreateCmd() {
	certCmd.AddCommand(createCmd)
	createCmd.Flags().IntP("years", "y", 10, "有效期")
	createCmd.Flags().StringSliceP("domain", "d", nil, "域名,默认为空,可指定多个")
}

func createPriviteCert(domain []string, years int) {
	cert.GeneratePrivateCert(domain, years)
}
