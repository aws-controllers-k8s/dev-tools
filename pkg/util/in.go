package util

// Gently copied from github.com/aws-controllers-k8s/runtime/pkg/util

// InStrings returns true if the subject string is contained in the supplied
// slice of strings
func InStrings(subject string, collection []string) bool {
	for _, item := range collection {
		if subject == item {
			return true
		}
	}
	return false
}

// InStringPs returns true if the subject string is contained in the supplied
// slice of string pointers
func InStringPs(subject string, collection []*string) bool {
	for _, item := range collection {
		if subject == *item {
			return true
		}
	}
	return false
}

// InsertString inserts a string in a given slice index
func InsertString(a []string, index int, value string) []string {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}

// InsertInterface inserts a interface{} implemetation in a given slice index
func InsertInterface(a []interface{}, index int, value interface{}) []interface{} {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}
