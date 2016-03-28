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
	"log"
	"os"
	"time"

	"github.com/minio/cli"
)

// ACME CA url. -- TODO make this configurable.
const (
	acmeStagingServer = "https://acme-staging.api.letsencrypt.org/directory"
	acmeServer        = "https://acme-v01.api.letsencrypt.org/directory"
)

func checkFolder(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0700)
	}
	return nil
}

// main for gen command.
func genMain(c *cli.Context) {
	if !c.Args().Present() || c.Args().First() == "help" {
		cli.ShowCommandHelpAndExit(c, "gen", 1) // last argument is exit code
	}

	// Create certs dir.
	certsDir := c.String("dir")
	if err := checkFolder(certsDir); err != nil {
		log.Fatalln(err)
	}

	// Get email and domain.
	email := c.Args().Get(0)
	domain := c.Args().Get(1)

	newCertificates, err := genCerts(email, domain)
	if err != nil {
		log.Fatalln(err)
	}

	// Each certificate comes back with the cert bytes, the bytes of
	// the client's private key, and a certificate URL. This is where
	// you should save them to files!
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
