package main

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
	"time"

	"github.com/minio/cli"
	"github.com/xenolf/lego/acme"
)

func renewMain(c *cli.Context) {
	if !c.Args().Present() || c.Args().First() == "help" {
		cli.ShowCommandHelpAndExit(c, "renew", 1) // last argument is exit code
	}

	// Renew keys from this dir.
	certsDir := c.String("dir")

	// Get email and domain.
	email := c.Args().Get(0)

	// Create a user. New accounts need an email and private key to start with.
	const rsaKeySize = 2048
	privateKey, err := rsa.GenerateKey(rand.Reader, rsaKeySize)
	if err != nil {
		log.Fatalln(err)
	}

	// Initialize user.
	user := User{
		Email: email,
		key:   privateKey,
	}

	// A client facilitates communication with the CA server. This CA
	// URL is configured for a local dev instance of Boulder running
	// in Docker in a VM.
	client, err := acme.NewClient(acmeServer, &user, acme.RSA2048)
	if err != nil {
		log.Fatalln(err)
	}

	client.ExcludeChallenges([]acme.Challenge{acme.DNS01})

	certBytes, err := loadCert(certsDir)
	if err != nil {
		log.Fatalln(err)
	}
	expTime, err := acme.GetPEMCertExpiration(certBytes)
	if int(expTime.Sub(time.Now()).Hours()/24.0) > 45 {
		log.Println("Keys have not expired yet, will not renew.")
		return
	}
	certMeta, err := loadCertMeta(certsDir)
	if err != nil {
		log.Fatalln(err)
	}

	certMeta.Certificate = certBytes

	isBundle := false
	newCertificates, err := client.RenewCertificate(certMeta, isBundle)
	if err != nil {
		log.Fatalln(err)
	}

	// New certificate comes back with the cert bytes, the bytes of
	// the client's private key, and a certificate URL. This is where
	// you should save them to files!
	err = saveCerts(certsDir, newCertificates)
	if err != nil {
		log.Fatalln(err)
	}

}
