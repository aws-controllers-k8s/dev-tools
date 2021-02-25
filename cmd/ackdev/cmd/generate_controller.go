package cmd

import (
	"fmt"
	"os/exec"
	"strconv"
	"sync"

	"github.com/spf13/cobra"

	"github.com/aws-controllers-k8s/dev-tools/pkg/aexec"
	"github.com/aws-controllers-k8s/dev-tools/pkg/config"
	"github.com/aws-controllers-k8s/dev-tools/pkg/repository"
	"github.com/aws-controllers-k8s/dev-tools/pkg/table"
)

var (
	optGenerateControllerParallel     bool
	optGenerateControllerSummaryTable bool
	optGenerateControllerStreamLogs   bool

	generateControllerResultsTableColumns = []interface{}{"SERVICE", "SUCCESS", "LOGFILE"}
)

func init() {
	generateControllerCmd.PersistentFlags().BoolVar(&optGenerateControllerParallel, "parallel", false, "generate controllers parallely")
	generateControllerCmd.PersistentFlags().BoolVar(&optGenerateControllerSummaryTable, "summary", true, "print summary table at the end")
	generateControllerCmd.PersistentFlags().BoolVar(&optGenerateControllerStreamLogs, "stream-logs", false, "stream code-generator logs")
}

var generateControllerCmd = &cobra.Command{
	Use:     "controller",
	Aliases: []string{"ctrl"},
	RunE:    generateController,
	Args:    cobra.MinimumNArgs(1),
}

func generateController(cmd *cobra.Command, args []string) error {
	if optGenerateControllerParallel && optGenerateControllerStreamLogs {
		return fmt.Errorf("flag conflict: --parallel, --stream-logs")
	}

	cfg, err := config.Load(ackConfigPath)
	if err != nil {
		return err
	}

	repoManager, err := repository.NewManager(cfg)
	if err != nil {
		return err
	}

	services := args
	if len(args) == 0 {
		services = cfg.Repositories.Services
	}

	codeGeneratorRepo, err := repoManager.LoadRepository("code-generator", repository.RepositoryTypeCore)
	if err != nil {
		return err
	}

	// Dragons start here

	//					\||/
	//					|  @___oo
	//		/\  /\   / (__,,,,|
	//		) /^\) ^\/ _)
	//		)   /^\/   _)
	//		)   _ /  / _)
	//	/\  )/\/ ||  | )_)
	//	<  >      |(,,) )__)
	//	||      /    \)___)\
	//	| \____(      )___) )___
	//	\______(_______;;; __;;;

	if optGenerateControllerParallel {
		genResults, err := generateControllersParallel(repoManager, codeGeneratorRepo.FullPath, services)
		if err != nil {
			return err
		}

		if optGenerateControllerSummaryTable {
			tablePrintGenerateResults(genResults)
		}
	} else {
		if optGenerateControllerStreamLogs {
			err := generateControllersVerbose(repoManager, codeGeneratorRepo.FullPath, services)
			if err != nil {
				return err
			}
		} else {
			genResults, err := generateControllersSequencial(repoManager, codeGeneratorRepo.FullPath, services)
			if err != nil {
				return err
			}

			if optGenerateControllerSummaryTable {
				tablePrintGenerateResults(genResults)
			}
		}
	}

	return nil
}

func tablePrintGenerateResults(gr []aexec.CmdResult) {
	tableHeaderColumns := generateControllerResultsTableColumns

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
	for _, result := range gr {
		rawArgs := []interface{}{result.Service, strconv.FormatBool(result.Success), "-"}
		if err := tw.AddRaw(rawArgs...); err != nil {
			panic(err)
		}
	}
}

func generateControllersVerbose(rl repository.Loader, contextDir string, services []string) error {
	for _, service := range services {
		err := runACKGenerateWithLogStream(contextDir, service)
		if err != nil {
			return err
		}
	}
	return nil
}

func generateControllersSequencial(rl repository.Loader, contextDir string, services []string) ([]aexec.CmdResult, error) {
	results := []aexec.CmdResult{}
	for _, service := range services {
		result := runACKGenerate(contextDir, service)
		results = append(results, result)
	}
	return results, nil
}

func generateControllersParallel(rl repository.Loader, contextDir string, services []string) ([]aexec.CmdResult, error) {
	results := []aexec.CmdResult{}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, service := range services {
		wg.Add(1)
		go func(service string) {
			defer wg.Done()

			result := runACKGenerate(contextDir, service)
			mu.Lock()
			defer mu.Unlock()
			results = append(results, result)
		}(service)
	}
	wg.Wait()
	return results, nil
}

func runACKGenerate(codeGeneratorPath, service string) aexec.CmdResult {
	cmd := exec.Command("make", "build-controller", "SERVICE="+service)
	cmd.Dir = codeGeneratorPath
	acmd := aexec.New(cmd, 50)
	err := acmd.Run(false)
	if err != nil {
		return aexec.CmdResult{
			Error:   err,
			Service: service,
			Success: false,
		}
	}
	err = acmd.Wait()
	return aexec.CmdResult{
		Error:   err,
		Service: service,
		Success: err == nil,
		Stdout:  acmd.Stdout,
		Stderr:  acmd.Stderr,
	}
}

func runACKGenerateWithLogStream(codeGeneratorPath, service string) error {
	cmd := exec.Command("make", "build-controller", "SERVICE="+service)
	cmd.Dir = codeGeneratorPath
	acmd := aexec.New(cmd, 50)
	err := acmd.Run(true)
	if err != nil {
		return err
	}
	done := make(chan struct{})

	go func() {
		for b := range acmd.StdoutStream() {
			fmt.Println(string(b))
		}
		done <- struct{}{}

	}()
	go func() {
		for b := range acmd.StderrStream() {
			fmt.Println(string(b))
		}
		done <- struct{}{}

	}()

	// wait for printers to finish
	defer func() { _, _ = <-done, <-done }()

	err = acmd.Wait()
	if err != nil {
		return err
	}

	return nil
}
