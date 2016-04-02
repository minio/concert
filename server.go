/*
 * Concert (C) 2016 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"errors"
	"log"
	"net"
	"strings"
	"time"

	"github.com/minio/cli"
)

// main for server command.
func serverMain(c *cli.Context) {
	if !c.Args().Present() || c.Args().First() == "help" {
		cli.ShowCommandHelpAndExit(c, "renew", 1) // last argument is
		// exit code
	}

	// Renew keys from this dir.
	certsDir := c.String("dir")

	// Sub domains if any.
	subDomainsValue := c.String("sub-domains")
	subDomains := strings.Split(subDomainsValue, ",")

	// Get email and domain.
	email := c.Args().Get(0)
	domain := c.Args().Get(1)

	// Validate if its a valid domain.
	if !isValidDomain(domain) {
		log.Fatalln(&net.DNSError{Err: "Invalid domain name", Name: domain})
	}
	// Validate if domain is already a sub domain, error if
	// sub-domains option is specified.
	if isSubDomain(domain) && len(subDomainsValue) > 0 {
		log.Fatalln(errors.New("Domain is already a subdomain, SSL certs not allowed for sub sub domains."))
	}

	if !isCertAvailable(certsDir) {
		newCertificates, err := genCerts(email, domain, subDomains)
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
		expTime, err := getCertExpTime(certsDir)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Generated certificates for %s under %s will expire in %d days.\n", domain, certsDir, int(expTime.Sub(time.Now()).Hours()/24.0))
	}
	log.Printf("Starting timer thread waiting for %d\n", renewDaysLimit)
	// Initialize a new timer, ticks every 45 days.
	for range time.Tick((renewDaysLimit - 1) * 24 * time.Hour) {
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
		expTime, err := getCertExpTime(certsDir)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Renewed certificates for %s under %s will expire in %d days.\n", domain, certsDir, int(expTime.Sub(time.Now()).Hours()/24.0))
	}
}
