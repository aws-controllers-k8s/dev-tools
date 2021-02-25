package cmd

import (
	"context"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/aws-controllers-k8s/dev-tools/pkg/config"
	"github.com/aws-controllers-k8s/dev-tools/pkg/crd"
	"github.com/aws-controllers-k8s/dev-tools/pkg/repository"
	"github.com/aws-controllers-k8s/dev-tools/pkg/table"
	"github.com/aws-controllers-k8s/dev-tools/pkg/util"
)

func init() {}

var (
	optDeployCRDsShowGroup      bool
	optDeployCRDsKubeConfigPath string
	deployCRDsTableHeader       = []interface{}{"SERVICE", "NAME", "VERSIONS", "STATUS"}
)

func init() {
	deployCRDsCmd.PersistentFlags().StringVarP(&optDeployCRDsKubeConfigPath, "kubeconfig", "k", filepath.Join(homeDirectory, ".kube", "config"), "kube config path")
	deployCRDsCmd.PersistentFlags().BoolVar(&optDeployCRDsShowGroup, "show-group", false, "display CRD group")
}

var deployCRDsCmd = &cobra.Command{
	Use:     "crd",
	Aliases: []string{"crds"},
	RunE:    deployCRDs,
	Args:    cobra.MinimumNArgs(1),
}

func deployCRDs(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(ackConfigPath)
	if err != nil {
		return err
	}

	repoManager, err := repository.NewManager(cfg)
	if err != nil {
		return err
	}

	services := args

	// CRD Client
	client, err := crd.New(repoManager, optDeployCRDsKubeConfigPath)
	if err != nil {
		return err
	}

	// table headers
	tableHeaderColumns := deployCRDsTableHeader
	if optDeployCRDsShowGroup {
		tableHeaderColumns = util.InsertInterface(tableHeaderColumns, 2, "GROUP")
	}

	// print headers
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

	addErrorRaw := func(service, msg string) {
		fmtArgs := []interface{}{service, "", "", msg}
		if optDeployCRDsShowGroup {
			fmtArgs = util.InsertInterface(fmtArgs, 2, "")
		}
		if err := tw.AddRaw(fmtArgs...); err != nil {
			panic(err)
		}
	}

	for _, service := range services {
		serviceRepo, err := repoManager.LoadRepository(service, repository.RepositoryTypeController)
		if err != nil || serviceRepo.GitHead == "" {
			addErrorRaw(service, "REPO ERROR")
			continue
		}

		crdDirectory := filepath.Join(serviceRepo.FullPath, "config/crd/bases")
		crdFiles, err := ioutil.ReadDir(crdDirectory)
		if err != nil {
			addErrorRaw(service, "CRD DIRECTORY ERROR")
			continue
		}

		for _, crdFile := range crdFiles {
			filePath := filepath.Join(crdDirectory, crdFile.Name())
			crdObject, err := crd.GetCRDFromFile(filePath)
			if err != nil {
				addErrorRaw(service, "CRD FORMAT ERROR")
				continue
			}

			err = client.DeployCRD(context.TODO(), crdObject)
			if err != nil {
				addErrorRaw(service, "K8S CLIENT ERROR")
				continue
			}

			// write raw
			versions := []string{}
			for _, crdVersion := range crdObject.Spec.Versions {
				versions = append(versions, crdVersion.Name)
			}
			rawArgs := []interface{}{service, crdObject.Name, versions, "OK"}
			if optDeployCRDsShowGroup {
				rawArgs = util.InsertInterface(rawArgs, 2, crdObject.Spec.Group)
			}
			if err := tw.AddRaw(rawArgs...); err != nil {
				panic(err)
			}
		}
	}

	return nil
}
