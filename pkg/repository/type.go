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

type RepositoryType int

const (
	RepositoryTypeCore RepositoryType = iota
	RepositoryTypeTooling
	RepositoryTypeController
	RepositoryTypeUnknown
)

// String stringifies a Repository type
func (rt RepositoryType) String() string {
	switch rt {
	case RepositoryTypeCore:
		return "core"
	case RepositoryTypeController:
		return "controller"
	case RepositoryTypeUnknown:
		return "UNKNOWN"
	default:
		panic("unsupported repository type")
	}
}

// repositoryTypeFromString casts a string to a RepositoryType
func repositoryTypeFromString(s string) RepositoryType {
	switch s {
	case "core":
		return RepositoryTypeCore
	case "controller":
		return RepositoryTypeController
	default:
		panic("unsupported repository type")
	}
}
