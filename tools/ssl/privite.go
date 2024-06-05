package ssl

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"gotools/utils/logger"
	"math/big"
	"os"
	"time"
)

func GeneratePrivateCert(domainList []string, years int) {
	// 指定证书有效期
	validFrom := time.Now()
	validTo := validFrom.Add(time.Duration(years) * 365 * 24 * time.Hour)

	// 生成私钥
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		logger.Fatal("Failed to generate private key:%v\n", err)
	}

	// 准备证书模板
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"zsops"},
		},
		NotBefore:             validFrom,
		NotAfter:              validTo,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	if domainList != nil {
		// 添加主题备用名称（SAN）
		template.DNSNames = domainList
	}

	// 生成证书
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privKey.PublicKey, privKey)
	if err != nil {
		logger.Fatal("Failed to create certificate:%v\n", err)
	}

	// 将证书保存到文件
	certOut, err := os.Create("cert.pem")
	if err != nil {
		logger.Fatal("Failed to open cert.pem for writing:%v\n", err)
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()
	logger.Success("Certificate saved to cert.pem\n")

	// 将私钥保存到文件
	keyOut, err := os.Create("key.pem")
	if err != nil {
		logger.Fatal("Failed to open key.pem for writing:%v\n", err)
	}
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privKey)})
	keyOut.Close()
	logger.Success("Private key saved to key.pem\n")
}
