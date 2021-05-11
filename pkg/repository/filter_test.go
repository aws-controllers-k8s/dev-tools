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

package repository

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildFilters(t *testing.T) {
	assert := assert.New(t)

	type args struct {
		expression string
	}
	tests := []struct {
		name          string
		args          args
		wantLenFilter int
		wantErr       bool
	}{
		{
			name:          "empty expression",
			args:          args{expression: "    "},
			wantLenFilter: 1,
			wantErr:       false,
		},
		{
			name:          "malformated expression - empty key",
			args:          args{expression: "key=value =value"},
			wantLenFilter: 0,
			wantErr:       true,
		},
		{
			name:          "malformated expression - empty value",
			args:          args{expression: "key=value key="},
			wantLenFilter: 0,
			wantErr:       true,
		},
		{
			name:          "unkown filter key",
			args:          args{expression: "name=runtime somekey=value"},
			wantLenFilter: 0,
			wantErr:       true,
		},
		{
			name:          "correct expression",
			args:          args{expression: "type=core"},
			wantLenFilter: 1,
			wantErr:       false,
		},
		{
			name:          "correct expression - all filters",
			args:          args{expression: "type=core branch=main name=runtime"},
			wantLenFilter: 3,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filters, err := BuildFilters(tt.args.expression)
			if (err != nil) != tt.wantErr {
				assert.Fail(fmt.Sprintf("BuildFilters() error = %v, wantErr %v", err, tt.wantErr))
			}
			assert.Len(filters, tt.wantLenFilter)
		})
	}
}

func TestNoFilter(t *testing.T) {
	repo := &Repository{}
	assert.True(t, NoFilter(repo))
	assert.True(t, NoFilter(nil))
}

func TestNameFilter(t *testing.T) {
	nameFilter := NameFilter("runtime")
	runtimeRepo := &Repository{
		Name: "runtime",
	}
	sqsRepo := &Repository{
		Name: "sqs",
	}
	assert.True(t, nameFilter(runtimeRepo))
	assert.False(t, nameFilter(sqsRepo))
}

func TestNamePrefixFilter(t *testing.T) {
	namePrefixfilter := NamePrefixFilter("com")
	runtimeRepo := &Repository{
		Name: "runtime",
	}
	communityRepo := &Repository{
		Name: "community",
	}
	assert.False(t, namePrefixfilter(runtimeRepo))
	assert.True(t, namePrefixfilter(communityRepo))
}

func TestTypeFilter(t *testing.T) {
	repoTypeFilter := TypeFilter(RepositoryTypeCore.String())
	runtimeRepo := &Repository{
		Name: "runtime",
		Type: RepositoryTypeCore,
	}
	sqsRepo := &Repository{
		Name: "sqs",
		Type: RepositoryTypeController,
	}
	assert.True(t, repoTypeFilter(runtimeRepo))
	assert.False(t, repoTypeFilter(sqsRepo))
}

func TestBranchFilter(t *testing.T) {
	branchFilter := BranchFilter("main")
	runtimeRepo := &Repository{
		Name:    "runtime",
		GitHead: "main",
		Type:    RepositoryTypeCore,
	}
	sqsRepo := &Repository{
		Name:    "sqs",
		GitHead: "feature-xyz",
		Type:    RepositoryTypeController,
	}
	assert.True(t, branchFilter(runtimeRepo))
	assert.False(t, branchFilter(sqsRepo))
}
