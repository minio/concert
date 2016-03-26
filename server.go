package main

import (
	"log"
	"time"

	"github.com/minio/cli"
)

func serverMain(c *cli.Context) {
	if !c.Args().Present() || c.Args().First() == "help" {
		cli.ShowCommandHelpAndExit(c, "renew", 1) // last argument is
		// exit code
	}

	// Renew keys from this dir.
	certsDir := c.String("dir")

	// Get email and domain.
	email := c.Args().Get(0)
	domain := c.Args().Get(1)

	if !isCertAvailable(certsDir) {
		newCertificates, err := genCerts(email, domain)
		if err != nil {
			log.Fatalln(err)
		}
		// Each certificate comes back with the cert bytes, the bytes
		// of the client's private key, and a certificate URL. This is
		// where you should save them to files!
		err = saveCerts(certsDir, newCertificates)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Generated certificates for %s under %s\n", domain, certsDir)
	}
	log.Printf("Starting timer thread waiting for %d\n", renewDaysLimit)
	// Initialize a new timer, ticks every 45 days.
	ticker := time.NewTicker((renewDaysLimit - 1) * 24 * time.Hour)
	for {
		select {
		case <-ticker.C:
			newCertificates, err := renewCerts(certsDir, email)
			if err != nil {
				log.Fatalln(err)
			}
			// Each certificate comes back with the cert bytes, the bytes
			// of the client's private key, and a certificate URL. This is
			// where you should save them to files!
			err = saveCerts(certsDir, newCertificates)
			if err != nil {
				log.Fatalln(err)
			}
			log.Printf("Renewed certificates for user %s under %s\n", email, certsDir)
		}
	}
}
