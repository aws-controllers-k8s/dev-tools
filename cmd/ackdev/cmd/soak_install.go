// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/release"

	"github.com/aws-controllers-k8s/dev-tools/pkg/config"
	"github.com/aws-controllers-k8s/dev-tools/pkg/repository"
	"github.com/aws-controllers-k8s/dev-tools/pkg/soak"
)

const (
	defaultControllerRegion      = "us-west-2"
	defaultControllerNamespace   = "ack-system"
	defaultControllerReleaseName = "soak-test"
	defaultPrometheusRepoName    = "prometheus-community"
	defaultPrometheusNamespace   = "prometheus"
	defaultPrometheusReleaseName = "kube-prom"
)

func init() {
	// soakInstallCmd.PersistentFlags().StringVar(&optBootstrapSoakClusterConfig, "cluster-config", defaultSoakClusterConfig, "eksctl cluster config file")
}

var soakInstallCmd = &cobra.Command{
	Use:     "install",
	RunE:    soakInstall,
	Args:    cobra.NoArgs,
	Short:   "Installs the soak test framework onto the cluster",
	Example: "ackdev soak install --cluster-config=./cluster-config.yml --service=s3",
}

func soakInstall(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(ackConfigPath)
	if err != nil {
		return err
	}

	repoManager, err := repository.NewManager(cfg)
	if err != nil {
		return err
	}

	// controllerRepo, err := repoManager.LoadRepository(optSoakService, repository.RepositoryTypeController)
	// if err != nil {
	// 	return err
	// }

	testInfraRepo, err := repoManager.LoadRepository("test-infra", repository.RepositoryTypeCore)
	if err != nil {
		return err
	}

	fmt.Printf("Adding Prometheus Helm chart repository ... ")
	// soak.AddHelmRepo(defaultPrometheusRepoName, "https://prometheus-community.github.io/helm-charts")
	fmt.Println("👍")

	fmt.Printf("Installing ACK %s controller Helm chart ... ", optSoakService)
	// controllerChart := filepath.Join(controllerRepo.FullPath, "helm")
	// controllerRelease, err := installController(optSoakService, controllerChart)
	// if err != nil {
	// 	return err
	// }
	fmt.Println("👍")

	fmt.Printf("Installing Prometheus Helm chart ... ")
	// _, err := installPrometheus(optSoakService)
	// if err != nil {
	// 	return err
	// }
	fmt.Println("👍")

	fmt.Printf("Applying ACK Grafana dashboard ... ")
	grafanaKustomization := filepath.Join(testInfraRepo.FullPath, "soak/monitoring/grafana")
	applyGrafanaDashboard(grafanaKustomization)
	fmt.Println("👍")

	return nil
}

func installController(service string, controllerChartPath string) (*release.Release, error) {
	chartValues := map[string]interface{}{
		"metrics": map[string]interface{}{
			"service": map[string]interface{}{
				"create": true,
				"type":   "ClusterIP",
			},
		},
		"aws": map[string]interface{}{
			"region": defaultControllerRegion,
		},
		"serviceAccount": map[string]interface{}{
			"create": false,
			"name":   soak.GetDefaultServiceAccountName(service),
		},
	}

	controllerRelease, err := soak.InstallLocalChart(controllerChartPath, defaultControllerNamespace, defaultControllerReleaseName, chartValues)
	if err != nil {
		return nil, err
	}

	return controllerRelease, nil
}

func installPrometheus(service string) (*release.Release, error) {
	truncService := service
	if len(truncService) > 44 {
		truncService = truncService[:44]
	}

	staticTargetURL := fmt.Sprintf("%s-controller-metrics.ack-system:8080", truncService)

	chartValues := map[string]interface{}{
		"prometheus": map[string]interface{}{
			"prometheusSpec": map[string]interface{}{
				"additionalScrapeConfigs": []interface{}{
					map[string]interface{}{
						"job_name": "ack-controller",
						"static_configs": []interface{}{
							map[string]interface{}{
								"targets": []string{staticTargetURL},
							}},
					},
				},
			},
		},
	}

	promRelease, err := soak.InstallRepoChart(
		defaultPrometheusRepoName,
		"kube-prometheus-stack",
		defaultPrometheusNamespace,
		defaultPrometheusReleaseName,
		chartValues,
	)
	if err != nil {
		return nil, err
	}

	return promRelease, nil
}

func applyGrafanaDashboard(dashboardKustomizationBasePath string) error {
	return soak.ApplyKustomization(dashboardKustomizationBasePath, defaultPrometheusNamespace)
}
