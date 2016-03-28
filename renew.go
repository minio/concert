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

	"github.com/minio/cli"
)

// main for renew command.
func renewMain(c *cli.Context) {
	if !c.Args().Present() || c.Args().First() == "help" {
		cli.ShowCommandHelpAndExit(c, "renew", 1) // last argument is exit code
	}

	// Renew keys from this dir.
	certsDir := c.String("dir")

	// Get email and domain.
	email := c.Args().Get(0)

	// Renew a certificate.
	newCertificates, err := renewCerts(certsDir, email)
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
	log.Printf("Renewed certificates for user %s under %s\n", email, certsDir)
}
