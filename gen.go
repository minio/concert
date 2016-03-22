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
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"log"
	"os"

	"github.com/minio/cli"
	"github.com/xenolf/lego/acme"
)

const acmeServer = "https://acme-staging.api.letsencrypt.org/directory"

// User - You'll need a user or account type that implements acme.User
type User struct {
	Email        string
	Registration *acme.RegistrationResource
	key          crypto.PrivateKey
}

// GetEmail -
func (u User) GetEmail() string {
	return u.Email
}

// GetRegistration -
func (u User) GetRegistration() *acme.RegistrationResource {
	return u.Registration
}

// GetPrivateKey -
func (u User) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func genMain(c *cli.Context) {
	if !c.Args().Present() || c.Args().First() == "help" {
		cli.ShowCommandHelpAndExit(c, "gen", 1) // last argument is exit code
	}

	// Create certs folder.
	certsDir := c.String("folder")
	if err := os.MkdirAll(certsDir, 0600); err != nil {
		log.Fatalln(err)
	}

	// Get email and domain.
	email := c.Args().Get(0)
	domain := c.Args().Get(1)

	// Create a user. New accounts need an email and private key to
	// start with.
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

	// We specify an http port of 5002 and an tls port of 5001 on all
	// interfaces because we aren't running as root and can't bind a
	// listener to port 80 and 443 (used later when we attempt to pass
	// challenges). Keep in mind that we still need to proxy challenge
	// traffic to port 5002 and 5001.
	err = client.SetHTTPAddress(":5002")
	if err != nil {
		log.Fatalln(err)
	}
	err = client.SetTLSAddress(":5001")
	if err != nil {
		log.Fatalln(err)
	}

	// New users will need to register; be sure to save it
	reg, err := client.Register()
	if err != nil {
		log.Fatalln(err)
	}
	user.Registration = reg

	// The client has a URL to the current Let's Encrypt Subscriber
	// Agreement. The user will need to agree to it.
	err = client.AgreeToTOS()
	if err != nil {
		log.Fatalln(err)
	}

	// The acme library takes care of completing the challenges to
	// obtain the certificate(s). Of course, the hostnames must
	// resolve to this machine or it will fail.
	isBundle := false
	certificates, failures := client.ObtainCertificate([]string{domain}, isBundle, nil)
	if len(failures) > 0 {
		log.Fatalln(failures)
	}

	// Each certificate comes back with the cert bytes, the bytes of
	// the client's private key, and a certificate URL. This is where
	// you should save them to files!
	err = saveCerts(certsDir, certificates)
	if err != nil {
		log.Fatalln(err)
	}
}
