// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package util

import (
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

// newSigner returns a ssh.Signer from a PEM encoded private key path.
// If the PEM file is encrypted it will try to read the passphrase from
// a terminal without local echo.
func NewSigner(sshKeyPath string) (ssh.Signer, error) {
	pemBytes, err := ioutil.ReadFile(sshKeyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("invalid ssh certificate")
	}

	if encryptedBlock(block) {
		passphrase, err := promptPassphrase()
		if err != nil {
			return nil, err
		}

		return ssh.ParsePrivateKeyWithPassphrase(pemBytes, passphrase)
	}

	return ssh.ParsePrivateKey(pemBytes)
}

// encryptedBlock tells whether a private key is
// encrypted by examining its Proc-Type header
// for a mention of ENCRYPTED
// according to RFC 1421 Section 4.6.1.1.
func encryptedBlock(block *pem.Block) bool {
	return strings.Contains(block.Headers["Proc-Type"], "ENCRYPTED")
}

func promptPassphrase() ([]byte, error) {
	fmt.Printf("type your ssh key passphrase:")
	passphrase, err := terminal.ReadPassword(0)
	if err != nil {
		return nil, err
	}
	// TODO(hilalymh): retry + validation mechanisms.
	return passphrase, nil
}
