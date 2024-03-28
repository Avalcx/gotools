package cert

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	cryptorand "crypto/rand"
	"log"
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

func Acme(domainList []string, ak string, sk string) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), cryptorand.Reader)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	cfg := alidns.NewDefaultConfig()
	cfg.APIKey = ak
	cfg.SecretKey = sk
	p, err := alidns.NewDNSProviderConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Challenge.SetDNS01Provider(p)
	if err != nil {
		log.Fatal(err)
	}

	// Registering acme account
	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		log.Fatal(err)
	}
	myUser.Registration = reg
	request := certificate.ObtainRequest{
		Domains: domainList,
		Bundle:  true,
	}
	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		log.Fatal(err)
	}
	os.MkdirAll(domainList[0], 0755)
	err = os.WriteFile(domainList[0]+"key.pem", certificates.PrivateKey, os.ModePerm)
	if err != nil {
		log.Print(err)
	}
	err = os.WriteFile(domainList[0]+"cert.pem", certificates.Certificate, os.ModePerm)
	if err != nil {
		log.Print(err)
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
