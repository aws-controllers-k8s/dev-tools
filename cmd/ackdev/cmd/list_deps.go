package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/aws-controllers-k8s/dev-tools/pkg/deps"
	"github.com/aws-controllers-k8s/dev-tools/pkg/table"
)

var (
	listDepsTableHeaderColumns = []interface{}{"NAME", "STATUS"}
	optDepsListShowPath        bool
	optDepsListShowVersion     bool
)

func init() {
	listDependenciesCmd.PersistentFlags().BoolVar(&optDepsListShowPath, "show-path", true, "display binary path")
	listDependenciesCmd.PersistentFlags().BoolVar(&optDepsListShowVersion, "show-version", true, "display binary version")
}

var listDependenciesCmd = &cobra.Command{
	Use:     "dependency",
	Aliases: []string{"dep", "deps", "dependencies"},
	RunE:    printDependencies,
}

type toolInfo struct {
	Name    string
	Version string
	Path    string
	Status  string
}

func printDependencies(cmd *cobra.Command, args []string) error {
	toolsInfo, err := listToolsInfo()
	if err != nil {
		return err
	}

	switch optListOutput {
	case "yaml":
		b, err := yaml.Marshal(toolsInfo)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	case "json":
		b, err := json.Marshal(toolsInfo)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	default:
		tablePrintDependencies(toolsInfo)
	}

	return nil
}

func tablePrintDependencies([]*toolInfo) error {
	// table headers
	tableHeaderColumns := listDepsTableHeaderColumns
	if optDepsListShowVersion {
		tableHeaderColumns = append(tableHeaderColumns, "VERSION")
	}
	if optDepsListShowPath {
		tableHeaderColumns = append(tableHeaderColumns, "PATH")
	}

	tw := table.NewPrinter(len(tableHeaderColumns))
	defer func() {
		if err := tw.Print(); err != nil {
			panic(err)
		}
	}()

	// print header
	if err := tw.AddRaw(tableHeaderColumns...); err != nil {
		panic(err)
	}

	dependencies, err := listToolsInfo()
	if err != nil {
		return err
	}

	for _, tool := range dependencies {
		rawArgs := []interface{}{tool.Name, tool.Status}
		if optDepsListShowVersion {
			rawArgs = append(rawArgs, tool.Version)
		}
		if optDepsListShowPath {
			rawArgs = append(rawArgs, tool.Path)
		}
		if err := tw.AddRaw(rawArgs...); err != nil {
			panic(err)
		}
	}

	return nil
}

func listToolsInfo() ([]*toolInfo, error) {
	list := make([]*toolInfo, 0, len(deps.DevelopmentTools))
	for _, tool := range deps.DevelopmentTools {
		status := ""
		path, err := tool.BinPath()
		if err != nil {
			status = "NOT FOUND"
		} else {
			status = "OK"
		}

		version := "-"
		if status == "OK" {
			v, err := tool.Version()
			if err != nil && err != deps.ErrorVersionNotFound {
				return nil, err
			}
			version = v
		}
		list = append(list, &toolInfo{
			Name:    tool.Name,
			Version: version,
			Path:    path,
			Status:  status,
		})
	}
	return list, nil
}
