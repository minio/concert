package main

import (
	"log"

	"github.com/minio/cli"
)

func renewMain(c *cli.Context) {
	if !c.Args().Present() || c.Args().First() == "help" {
		cli.ShowCommandHelpAndExit(c, "renew", 1) // last argument is exit code
	}

	// Renew keys from this dir.
	certsDir := c.String("dir")

	// Get email and domain.
	email := c.Args().Get(0)

	// Renew a certificate.
	newCertificates, err := renewCerts(email, certsDir)
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
