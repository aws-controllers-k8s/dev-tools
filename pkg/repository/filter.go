package repository

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrorMalformatedFilterExpression error = errors.New("malformated filter expression")
)

// NewFiltersFromExpression takes an expression string and returns a list
// of Filter functions. Example: "branch=main type=controller"
func NewFiltersFromExpression(expression string) ([]Filter, error) {
	if expression == "" {
		return nil, nil
	}
	expressionElements := strings.Split(expression, " ")
	filters := []Filter{}

	for _, filterStr := range expressionElements {
		filterArgs := strings.Split(filterStr, "=")
		if len(filterArgs) != 2 {
			return nil, fmt.Errorf("malformatted filter expression")
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
			return nil, fmt.Errorf("unknown filter key %s", key)
		}
	}
	return filters, nil
}

// Filter is a function that filtrate a single repository
type Filter func(r *Repository) bool

// NoFilter always return true
func NoFilter() Filter {
	return func(r *Repository) bool { return true }
}

// NameFilter filters a repository by a exact name
func NameFilter(name string) Filter {
	return func(r *Repository) bool {
		return r.Name == name
	}
}

// NamePrefixFilter filters a repository by a exact name
func NamePrefixFilter(namePrefix string) Filter {
	return func(r *Repository) bool {
		return strings.HasPrefix(r.Name, namePrefix)
	}
}

// TypeFilter filters a repository by a name prefix
// The only two possible types are 'controller' and 'core'
func TypeFilter(t string) Filter {
	return func(r *Repository) bool {
		return r.Type.String() == t
	}
}

// BranchFilter filters a repository by the current branch name
func BranchFilter(branch string) Filter {
	return func(r *Repository) bool {
		return r.GitHead == branch
	}
}
