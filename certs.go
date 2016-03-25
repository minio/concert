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
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/xenolf/lego/acme"
)

// renewDaysLimit - renewal is not initated if cert is valid with in this limit.
const renewDaysLimit = 45 // Number of days.

func genCerts(email, domain string) (acme.CertificateResource, error) {
	// Create a user. New accounts need an email and private key to start with.
	const rsaKeySize = 2048
	privateKey, err := rsa.GenerateKey(rand.Reader, rsaKeySize)
	if err != nil {
		return acme.CertificateResource{}, err
	}

	// Initialize user.
	user := conUser{
		Email: email,
		key:   privateKey,
	}

	// A client facilitates communication with the CA server.
	client, err := acme.NewClient(acmeServer, &user, acme.RSA2048)
	if err != nil {
		return acme.CertificateResource{}, err
	}

	client.ExcludeChallenges([]acme.Challenge{acme.DNS01})

	// New users will need to register; be sure to save it
	reg, err := client.Register()
	if err != nil {
		return acme.CertificateResource{}, err
	}
	user.Registration = reg

	// The client has a URL to the current Let's Encrypt Subscriber
	// Agreement. The user will need to agree to it.
	err = client.AgreeToTOS()
	if err != nil {
		return acme.CertificateResource{}, err
	}

	// The acme library takes care of completing the challenges to
	// obtain the certificate(s). Of course, the hostnames must
	// resolve to this machine or it will fail.
	isBundle := false
	newCertificates, failures := client.ObtainCertificate([]string{domain}, isBundle, nil)
	if len(failures) > 0 {
		return acme.CertificateResource{}, failures[domain]
	}
	return newCertificates, nil
}

// Renew certs.
func renewCerts(certsDir, email string) (acme.CertificateResource, error) {
	certBytes, err := loadCert(certsDir)
	if err != nil {
		return acme.CertificateResource{}, err
	}

	expTime, err := acme.GetPEMCertExpiration(certBytes)
	if int(expTime.Sub(time.Now()).Hours()/24.0) > renewDaysLimit {
		return acme.CertificateResource{}, errors.New("Keys have not expired yet, will not renew.")
	}

	// Create a user. New accounts need an email and private key to start with.
	const rsaKeySize = 2048
	privateKey, err := rsa.GenerateKey(rand.Reader, rsaKeySize)
	if err != nil {
		return acme.CertificateResource{}, err
	}

	// Initialize user.
	user := conUser{
		Email: email,
		key:   privateKey,
	}

	// A client facilitates communication with the CA server. This CA
	// URL is configured for a local dev instance of Boulder running
	// in Docker in a VM.
	client, err := acme.NewClient(acmeServer, &user, acme.RSA2048)
	if err != nil {
		return acme.CertificateResource{}, err
	}

	client.ExcludeChallenges([]acme.Challenge{acme.DNS01})

	certMeta, err := loadCertMeta(certsDir)
	if err != nil {
		return acme.CertificateResource{}, err
	}

	certMeta.Certificate = certBytes

	isBundle := false
	newCertificates, err := client.RenewCertificate(certMeta, isBundle)
	if err != nil {
		return acme.CertificateResource{}, err
	}
	return newCertificates, nil
}

// load certificate meta resource.
func loadCertMeta(certsDir string) (acme.CertificateResource, error) {
	metaBytes, err := ioutil.ReadFile(filepath.Join(certsDir, "certs.json"))
	if err != nil {
		return acme.CertificateResource{}, err
	}
	var certRes acme.CertificateResource
	err = json.Unmarshal(metaBytes, &certRes)
	if err != nil {
		return acme.CertificateResource{}, err
	}
	return certRes, nil
}

// load certs.
func loadCert(certsDir string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(certsDir, "public.crt"))
}

// saveCerts saves the certificates to disk. This includes the
// certificate file itself, the private key, and the json metadata file.
func saveCerts(certsDir string, cert acme.CertificateResource) error {
	// Save cert file.
	err := ioutil.WriteFile(filepath.Join(certsDir, "public.crt"), cert.Certificate, 0600)
	if err != nil {
		return err
	}

	// Save private key.
	err = ioutil.WriteFile(filepath.Join(certsDir, "private.key"), cert.PrivateKey, 0600)
	if err != nil {
		return err
	}

	// Save cert metadata.
	jsonBytes, err := json.MarshalIndent(&cert, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(certsDir, "certs.json"), jsonBytes, 0600)
	if err != nil {
		return err
	}

	// Return success.
	return nil
}

// Verify if certs are available in a certs dir.
func isCertAvailable(certsDir string) bool {
	_, err := os.Stat(filepath.Join(certsDir, "public.crt"))
	if err != nil {
		return false
	}
	_, err = os.Stat(filepath.Join(certsDir, "private.key"))
	if err != nil {
		return false
	}
	return true
}
