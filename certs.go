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
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/xenolf/lego/acme"
)

// saveCerts saves the certificates to disk. This includes the
// certificate file itself, the private key, and the json metadata
// file.
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
