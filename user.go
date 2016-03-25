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

	"github.com/xenolf/lego/acme"
)

// conUser - You'll need a user or account type that implements acme.User
type conUser struct {
	Email        string
	Registration *acme.RegistrationResource
	key          crypto.PrivateKey
}

// GetEmail -
func (u conUser) GetEmail() string {
	return u.Email
}

// GetRegistration -
func (u conUser) GetRegistration() *acme.RegistrationResource {
	return u.Registration
}

// GetPrivateKey -
func (u conUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}
