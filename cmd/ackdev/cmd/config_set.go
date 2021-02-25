package cmd

import (
	"encoding/json"

	"github.com/spf13/cobra"
	"github.com/tidwall/sjson"

	"github.com/aws-controllers-k8s/dev-tools/pkg/config"
)

func init() {}

var setConfigCmd = &cobra.Command{
	Use:  "set",
	Args: cobra.ExactArgs(2),
	RunE: setConfigJSONPath,
}

func setConfigJSONPath(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(ackConfigPath)
	if err != nil {
		return err
	}

	jsonBytes, err := json.Marshal(cfg)
	if err != nil {
		return err // maybe panic ?
	}

	newJSONbytes, err := sjson.SetBytes(jsonBytes, args[0], args[1])
	if err != nil {
		return err
	}

	newCfg, err := config.LoadFromJSONBytes(newJSONbytes)
	if err != nil {
		return err
	}

	return config.Save(newCfg, ackConfigPath)
}
