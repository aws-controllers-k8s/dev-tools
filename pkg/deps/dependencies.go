package deps

import (
	"errors"
	"os/exec"
	"regexp"
)

var (
	ErrorVersionNotFound = errors.New("version not found")
)

const (
	ScriptsRepositoryName = "community"
)

var (
	DevelopmentTools = []Tool{

		{
			Name:           "go",
			GetVersionArgs: []string{"version"},
		},
		{
			Name:           "kind",
			GetVersionArgs: []string{"--version"},
		},
		{
			Name:           "helm",
			GetVersionArgs: []string{"version", "--short"},
		},
		{Name: "mockery"},
		{
			Name:           "kubectl",
			GetVersionArgs: []string{"version", "--client", "--short"},
		},
		{
			Name:           "kustomize",
			GetVersionArgs: []string{"version", "--short"},
		},
		{
			Name:           "controller-gen",
			GetVersionArgs: []string{"--version"},
		},
	}
)

type Tool struct {
	Name string
	// This is currently not used. Thinking about embedding scripts
	// into the binary ...
	InstallScriptName string
	GetVersionArgs    []string
}

func (t *Tool) BinPath() (string, error) {
	path, err := exec.LookPath(t.Name)
	if err != nil {
		return "", err
	}
	return path, nil
}

func (t *Tool) Install(scriptPath string) error {
	return nil
}

func (t *Tool) Version() (string, error) {
	cmd := exec.Command(t.Name, t.GetVersionArgs...)
	b, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	matches := getVersionFromString(string(b))
	if err != nil {
		return "", err
	}
	if len(matches) == 0 {
		return "", ErrorVersionNotFound
	}

	return matches[0], nil
}

func getVersionFromString(s string) []string {
	regexExpr := `v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?` +
		`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?` +
		`(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?`
	re := regexp.MustCompile(regexExpr) // panics if Compile errors
	match := re.FindStringSubmatch(s)
	return match
}
