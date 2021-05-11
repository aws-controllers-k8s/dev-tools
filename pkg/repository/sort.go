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

import "sort"

// Gently stolen from github.com/aws-controllers-k8s/code-generator/pkg/model/printer_column.go

// By can sort two Repositories
type By func(a, b *Repository) bool

// Sort does an in-place sort of the supplied repositories
func (by By) Sort(subject []*Repository) {
	pcs := repositorySorter{
		cols: subject,
		by:   by,
	}
	sort.Sort(pcs)
}

// repositorySorter sorts repositories
type repositorySorter struct {
	cols []*Repository
	by   By
}

// Len implements sort.Interface.Len
func (pcs repositorySorter) Len() int {
	return len(pcs.cols)
}

// Swap implements sort.Interface.Swap
func (pcs repositorySorter) Swap(i, j int) {
	pcs.cols[i], pcs.cols[j] = pcs.cols[j], pcs.cols[i]
}

// Less implements sort.Interface.Less
func (pcs repositorySorter) Less(i, j int) bool {
	return pcs.by(pcs.cols[i], pcs.cols[j])
}

// Sort two repositories by name
func ByName(a, b *Repository) bool {
	return a.Name < b.Name
}

// Sort two repositories by branch
func ByBranch(a, b *Repository) bool {
	return a.GitHead < b.GitHead
}

// Sort two repositories by type
func ByType(a, b *Repository) bool {
	return a.Type < b.Type
}

// SortBy takes a field path and returns the equivalent Sorter function
func SortBy(fieldPath string) By {
	switch fieldPath {
	case "name":
		return ByName
	case "branch":
		return ByBranch
	case "type":
		return ByType
	default:
		panic("unknown sort field")
	}
}
