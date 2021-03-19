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

package deps

import (
	"errors"
	"os/exec"
	"regexp"
)

var (
	ErrorVersionNotFound = errors.New("version not found in output")
)

var (
	// DevelopmentTools is the list of ACK development tools
	DevelopmentTools = []Dependency{
		{
			BinaryName:     "go",
			GetVersionArgs: []string{"version"},
		},
		{
			BinaryName:     "kind",
			GetVersionArgs: []string{"--version"},
		},
		{
			BinaryName:     "helm",
			GetVersionArgs: []string{"version", "--short"},
		},
		{
			BinaryName:     "mockery",
			GetVersionArgs: []string{"--version", "--quiet"},
		},
		{
			BinaryName:     "kubectl",
			GetVersionArgs: []string{"version", "--client", "--short"},
		},
		{
			BinaryName:     "kustomize",
			GetVersionArgs: []string{"version", "--short"},
		},
		{
			BinaryName:     "controller-gen",
			GetVersionArgs: []string{"--version"},
		},
	}
)

// Dependency represent an ACK development dependency. Generally
// dependencies are binaries.
type Dependency struct {
	// Expected binary name
	BinaryName string
	// Arguments passed to the binary in order to get it version
	GetVersionArgs []string
}

// BinPath returns the path of a binary if it exists
func (t *Dependency) BinPath() (string, error) {
	path, err := exec.LookPath(t.BinaryName)
	if err != nil {
		return "", err
	}
	return path, nil
}

// Version returns the version of the binary.
func (t *Dependency) Version() (string, error) {
	cmd := exec.Command(t.BinaryName, t.GetVersionArgs...)
	b, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	version := getVersionFromString(string(b))
	if len(version) == 0 {
		return "", ErrorVersionNotFound
	}

	return version, nil
}

// getVersionFromString parses a string expression and returns the first
// observed semantic version.
// TODO(a-hilaly) find a better regex expression - currently this one returns
// v1 if it gets 'thisisnotav1ersion'
func getVersionFromString(s string) string {
	regexExpr := `v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?` +
		`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?` +
		`(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?`
	re := regexp.MustCompile(regexExpr) // panics if it doesn't compile
	matches := re.FindStringSubmatch(s)
	if len(matches) == 0 {
		return ""
	}
	return matches[0]
}
