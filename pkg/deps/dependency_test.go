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
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_getVersionFromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"empty input",
			args{""},
			"",
		},
		{
			"expression not containing any version",
			args{"this is an input that doesn't contain any version"},
			"",
		},
		{
			"expression containing a version with a 'v' prefix",
			args{"someoutput v1.0.0-rc1 someotheroutput"},
			"v1.0.0-rc1",
		},
		{
			"expression containing a version without a 'v' prefix",
			args{"someoutput 2.0.0-rc3 someotheroutput"},
			"2.0.0-rc3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := getVersionFromString(tt.args.s)
			require.EqualValues(t, matches, tt.want)
		})
	}
}
