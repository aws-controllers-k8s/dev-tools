package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/aws-controllers-k8s/dev-tools/pkg/aexec"
	"github.com/aws-controllers-k8s/dev-tools/pkg/config"
	"github.com/aws-controllers-k8s/dev-tools/pkg/repository"
)

var (
	optRunUseLocalModules string
)

func init() {}

var runCmd = &cobra.Command{
	Use:           "run",
	SilenceErrors: true,
	RunE:          runController,
	Args:          cobra.ExactArgs(1),
	Short:         "Runs a controller binary locally",
}

func runController(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(ackConfigPath)
	if err != nil {
		return err
	}

	repoManager, err := repository.NewManager(cfg)
	if err != nil {
		return err
	}

	repo, err := repoManager.LoadRepository(args[0], repository.RepositoryTypeController)
	if err != nil {
		return err
	}

	controllerRunFlags := []string{}
	for flag, value := range cfg.RunConfig.Flags {
		controllerRunFlags = append(
			controllerRunFlags,
			fmt.Sprintf("--%s=%s", flag, value),
		)
	}

	err = cmdRunController(repo.FullPath, controllerRunFlags...)
	if err != nil {
		return err
	}

	return nil
}

func cmdRunController(path string, flags ...string) error {
	cmdArgs := append([]string{
		"run",
		"-tags=codegen",
		"./cmd/controller/main.go",
	}, flags...)
	cmd := exec.Command("go", cmdArgs...)
	cmd.Dir = path
	acmd := aexec.New(cmd, 50)
	err := acmd.Run(true)
	if err != nil {
		return err
	}
	go func() {
		for b := range acmd.StdoutStream() {
			fmt.Println(string(b))
		}
	}()
	go func() {
		for b := range acmd.StderrStream() {
			fmt.Println(string(b))
		}
	}()

	acmd.Wait()
	return nil
}
