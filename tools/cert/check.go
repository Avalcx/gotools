package cert

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func CheckFromDomain(domain string) {
	conn, err := net.Dial("tcp", domain+":443")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	config := tls.Config{
		ServerName: domain,
	}

	tlsConn := tls.Client(conn, &config)

	err = tlsConn.Handshake()
	if err != nil {
		fmt.Println(err.Error())
	}

	cert := tlsConn.ConnectionState().PeerCertificates[0]
	certInfo := CertInfo{}
	getCertInfo(cert, &certInfo)
	printCertInfo(&certInfo)
}

func CheckFromCrtFile(file string) {
	certBytes, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("无法读取证书文件：%v", err)
		return
	}

	// 解码 PEM 格式的证书
	block, _ := pem.Decode(certBytes)
	if block == nil {
		fmt.Printf("无法解码 PEM 格式的证书：%s", err)
		return
	}
	// 解析证书
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		fmt.Printf("无法解析证书：%s", err)
		return
	}

	certInfo := CertInfo{}
	getCertInfo(cert, &certInfo)
	printCertInfo(&certInfo)
}

type CertInfo struct {
	Domain    string
	StartTime string
	EndTime   string
	Expire    string
}

func getCertInfo(cert *x509.Certificate, certInfo *CertInfo) {
	certInfo.Domain = strings.Join(cert.DNSNames, " ")
	certInfo.StartTime = cert.NotBefore.Format(time.DateTime)
	certInfo.EndTime = cert.NotAfter.Format(time.DateTime)
	expire := int(cert.NotAfter.Sub(cert.NotBefore).Hours())
	certInfo.Expire = strconv.Itoa(expire/24) + "days " + strconv.Itoa(expire%24) + "hours"
}

func printCertInfo(certInfo *CertInfo) {
	fmt.Printf("domain: %s\nstartTime: %s\nendTime: %s\nexpire: %v", certInfo.Domain, certInfo.StartTime, certInfo.EndTime, certInfo.Expire)
}
