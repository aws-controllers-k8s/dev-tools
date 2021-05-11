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
	"errors"
	"fmt"
	"strings"
)

var (
	ErrMalformatedFilterExpression error = errors.New("malformated filter expression")
	ErrUnknownFilterKey            error = errors.New("unknown filter key")
)

// BuildFilters takes an expression string and returns a list
// of Filter functions. Example: "branch=main type=controller"
func BuildFilters(expression string) ([]Filter, error) {
	expression = strings.TrimSpace(expression)
	if len(expression) == 0 {
		return []Filter{NoFilter}, nil
	}

	var filters []Filter
	expressionElements := strings.Split(expression, " ")
	for _, filterStr := range expressionElements {
		// TODO(hilalymh) probably we can use a regular expression instead..
		filterArgs := strings.Split(filterStr, "=")
		if len(filterArgs) != 2 {
			return nil, ErrMalformatedFilterExpression
		}
		key := strings.ToLower(filterArgs[0])
		value := filterArgs[1]
		switch key {
		case "type":
			filters = append(filters, TypeFilter(value))
		case "name":
			filters = append(filters, NameFilter(value))
		case "branch":
			filters = append(filters, BranchFilter(value))
		default:
			return nil, fmt.Errorf("unknown filter key: %s", key)
		}
	}
	return filters, nil
}

// A Filter is a prototype for a function that can be used to filter the
// results from a call to the List() method on the Manager.
type Filter func(r *Repository) bool

// NoFilter will not filter out any repository.
func NoFilter(r *Repository) bool { return true }

// NameFilter filters all repositories whose names matches the specified name
func NameFilter(name string) Filter {
	return func(r *Repository) bool {
		return r.Name == name
	}
}

// NamePrefixFilter filters all repositories whose name prefix matches the
// the given namePrefix
func NamePrefixFilter(namePrefix string) Filter {
	return func(r *Repository) bool {
		return strings.HasPrefix(r.Name, namePrefix)
	}
}

// TypeFilter filters a repository by a name prefix
// The only two possible types are 'controller' and 'core'
func TypeFilter(t string) Filter {
	return func(r *Repository) bool {
		return r.Type == repositoryTypeFromString(t)
	}
}

// BranchFilter filters all repositories whose current branch matches the
// given branch name.
func BranchFilter(branch string) Filter {
	return func(r *Repository) bool {
		return r.GitHead == branch
	}
}
