package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/aws-controllers-k8s/dev-tools/pkg/config"
	"github.com/aws-controllers-k8s/dev-tools/pkg/repository"
	"github.com/aws-controllers-k8s/dev-tools/pkg/table"
)

var (
	listTableHeaderColumns = []interface{}{"NAME", "TYPE"}

	optListFilterExpression string
	optListShowCloneURL     bool
	optListShowBranch       bool
)

func init() {
	listRepositoriesCmd.PersistentFlags().StringVarP(&optListFilterExpression, "filter", "f", "", "filter expression")
	listRepositoriesCmd.PersistentFlags().BoolVar(&optListShowCloneURL, "show-url", false, "display project remote URL or not")
	listRepositoriesCmd.PersistentFlags().BoolVar(&optListShowBranch, "show-branch", true, "display project current branch or not")
}

var listRepositoriesCmd = &cobra.Command{
	Use:     "repository",
	Aliases: []string{"repo", "repos", "repositories"},
	RunE:    printRepositories,
	Args:    cobra.NoArgs,
}

func printRepositories(cmd *cobra.Command, args []string) error {
	filters, err := repository.NewFiltersFromExpression(optListFilterExpression)
	if err != nil {
		return err
	}

	repos, err := listRepositories(filters...)
	if err != nil {
		return err
	}

	switch optListOutput {

	case "yaml":
		b, err := yaml.Marshal(repos)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	case "json":
		b, err := json.Marshal(repos)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	default:
		tablePrintRepositories(repos)
	}

	return nil
}

func listRepositories(filters ...repository.Filter) ([]*repository.Repository, error) {
	cfg, err := config.Load(ackConfigPath)
	if err != nil {
		return nil, err
	}
	repoManager, err := repository.NewManager(cfg)
	if err != nil {
		return nil, err
	}

	// Try to load all repositories
	err = repoManager.LoadAll()
	if err != nil {
		return nil, err
	}

	// List repositories
	repos := repoManager.List(filters...)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

func tablePrintRepositories(repos []*repository.Repository) {

	tableHeaderColumns := listTableHeaderColumns
	if optListShowBranch {
		tableHeaderColumns = append(tableHeaderColumns, "BRANCH")
	}
	if optListShowCloneURL {
		tableHeaderColumns = append(tableHeaderColumns, "URL")
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

	// print raws
	for _, repo := range repos {
		rawArgs := []interface{}{repo.Name, repo.Type.String()}
		if optListShowBranch {
			rawArgs = append(rawArgs, repo.GitHead)
		}
		if optListShowCloneURL {
			rawArgs = append(rawArgs, repo.RemoteURL)
		}
		if err := tw.AddRaw(rawArgs...); err != nil {
			panic(err)
		}
	}
}
