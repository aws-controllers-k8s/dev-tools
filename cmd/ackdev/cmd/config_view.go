package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/aws-controllers-k8s/dev-tools/pkg/config"
)

var (
	optViewConfigOutput string
)

func init() {
	viewConfigCmd.PersistentFlags().StringVarP(&optViewConfigOutput, "output", "o", "yaml", "output format (json|yaml)")
}

var viewConfigCmd = &cobra.Command{
	Use:  "view",
	Args: cobra.NoArgs,
	RunE: printConfig,
}

func printConfig(*cobra.Command, []string) error {
	cfg, err := config.Load(ackConfigPath)
	if err != nil {
		return err
	}

	var b []byte
	switch optViewConfigOutput {
	case "json":
		b, err = json.Marshal(cfg)
		if err != nil {
			return err
		}
	case "yaml":
		b, err = yaml.Marshal(cfg)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported output type: %s", optViewConfigOutput)
	}

	fmt.Println(string(b))
	return nil
}
