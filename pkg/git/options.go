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

package git

import "golang.org/x/crypto/ssh"

type Option func(*Git)

// WithRemote sets the remote attached to cloning URL.
func WithRemote(remote string) Option {
	return func(g *Git) {
		g.remote = remote
	}
}

// WithGithubCredentials sets the Github username and password used to clone
// repositories with HTTPS protocol.
func WithGithubCredentials(username, token string) Option {
	return func(g *Git) {
		g.githubUsername = username
		g.githubToken = token
	}
}

// WithSSHSigner sets the ssh.Signer used to clone repositories with
// ssh protocol.
func WithSSHSigner(signer ssh.Signer) Option {
	return func(g *Git) {
		g.signer = signer
	}
}
