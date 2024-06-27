package ssl

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	cryptorand "crypto/rand"
	"gotools/utils/logger"
	"math/rand"
	"os"
	"time"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns/alidns"
	"github.com/go-acme/lego/v4/registration"
)

type MyUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *MyUser) GetEmail() string {
	return u.Email
}
func (u MyUser) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *MyUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func (sslInfo *SSLInfo) Acme() {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), cryptorand.Reader)
	if err != nil {
		logger.Fatal("%v\n", err)
	}
	email := generateEmail()
	myUser := MyUser{
		Email: email,
		key:   privateKey,
	}
	config := lego.NewConfig(&myUser)
	config.Certificate.KeyType = certcrypto.RSA2048
	client, err := lego.NewClient(config)
	if err != nil {
		logger.Fatal("%v\n", err)
	}

	cfg := alidns.NewDefaultConfig()
	cfg.APIKey = sslInfo.AliAK
	cfg.SecretKey = sslInfo.AliSK
	p, err := alidns.NewDNSProviderConfig(cfg)
	if err != nil {
		logger.Fatal("%v\n", err)
	}
	err = client.Challenge.SetDNS01Provider(p)
	if err != nil {
		logger.Fatal("%v\n", err)
	}

	// Registering acme account
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		logger.Fatal("%v\n", err)
	}
	myUser.Registration = reg
	request := certificate.ObtainRequest{
		Domains: sslInfo.Domains,
		Bundle:  true,
	}
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		logger.Fatal("%v\n", err)
	}
	os.MkdirAll(sslInfo.Domains[0], 0755)
	err = os.WriteFile("acme/key.pem", certificates.PrivateKey, os.ModePerm)
	if err != nil {
		logger.Failed("%v\n", err)
	}
	err = os.WriteFile("acme/cert.pem", certificates.Certificate, os.ModePerm)
	if err != nil {
		logger.Failed("%v\n", err)
	}
}

func generateEmail() string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	digits := make([]byte, 10)
	for i := range digits {
		digits[i] = byte(r.Intn(10) + '0')
	}
	return string(digits) + "@qq.com"
}
