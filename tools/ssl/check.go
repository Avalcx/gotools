package ssl

import (
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"gotools/utils/logger"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type CertInfo struct {
	Domain        string
	StartTime     string
	EndTime       string
	Expire        string
	IsPriviteCert bool
	Issuer        pkix.Name
}

func (sslInfo *SSLInfo) Check() {
	if (len(sslInfo.Domains) == 0 && sslInfo.CertFile == "") || (len(sslInfo.Domains) != 0 && sslInfo.CertFile != "") {
		logger.Fatal("-d和-f至少且只能有一个参数")
	}
	if len(sslInfo.Domains) >= 1 {
		sslInfo.checkFromDomains()
	} else if sslInfo.CertFile != "" {
		sslInfo.checkFromCrtFile()
	}
}

func (sslInfo *SSLInfo) checkFromDomains() {
	for _, domain := range sslInfo.Domains {
		conn, err := net.Dial("tcp", domain+":"+sslInfo.Port)
		if err != nil {
			logger.Failed("%v\n", err)
			return
		}
		defer conn.Close()

		config := tls.Config{
			ServerName: domain,
		}

		tlsConn := tls.Client(conn, &config)

		err = tlsConn.Handshake()
		if errors.Is(err, io.EOF) {
			logger.Failed("目标域名错误或网络异常\n")
			return
		} else if !errors.Is(err, nil) {
			logger.Failed("%v\n", err)
		}

		if len(tlsConn.ConnectionState().PeerCertificates) < 1 {
			logger.Failed("证书读取异常\n")
			return
		}

		cert := tlsConn.ConnectionState().PeerCertificates[0]
		c := CertInfo{}
		c.getCertInfo(cert)
		c.printCertInfo()
	}

}

func (sslInfo *SSLInfo) checkFromCrtFile() {
	certBytes, err := os.ReadFile(sslInfo.CertFile)
	if err != nil {
		logger.Failed("无法读取证书文件：%v\n", err)
		return
	}

	// 解码 PEM 格式的证书
	block, _ := pem.Decode(certBytes)
	if block == nil {
		logger.Failed("无法解码 PEM 格式的证书：%v\n", err)
		return
	}
	// 解析证书
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		logger.Failed("无法解析证书：%v\n", err)
		return
	}

	c := CertInfo{}
	c.getCertInfo(cert)
	c.printCertInfo()
}

func (c *CertInfo) getCertInfo(cert *x509.Certificate) {
	c.Domain = strings.Join(cert.DNSNames, " ")
	c.StartTime = cert.NotBefore.Local().Format(time.DateTime)
	c.EndTime = cert.NotAfter.Local().Format(time.DateTime)
	now := time.Now()
	expire := int(cert.NotAfter.Sub(now).Hours())
	c.Expire = strconv.Itoa(expire / 24)
	if len(cert.IssuingCertificateURL) == 0 {
		c.IsPriviteCert = true
	} else {
		c.IsPriviteCert = false
	}
	c.Issuer = cert.Issuer
}

func (c *CertInfo) printCertInfo() {
	logger.Printf("domain: %s\nstartTime: %s\nendTime: %s\nexpire: %v\nisPriviteCert: %v\nIssuer: %v\n", c.Domain, c.StartTime, c.EndTime, c.Expire, c.IsPriviteCert, c.Issuer)
}
