package soak

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/aws-controllers-k8s/dev-tools/pkg/util"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/eks"
)

var (
	defaultClusterName                = "ack-soak-test"
	defaultControllerInstallNamespace = "ack-system"
	defaultServiceAccountNameFmt      = "ack-%s-controller"
)

func EnsureCluster(clusterConfigPath string, service string) error {
	existing, err := getCluster()
	if err != nil {
		return err
	}

	if existing == nil {
		if err = CreateCluster(clusterConfigPath); err != nil {
			return err
		}
		existing, err := getCluster()
		if err != nil {
			return err
		} else if existing == nil {
			return fmt.Errorf("cluster not found after creating")
		}
	}

	if err = SetupIRSA(service); err != nil {
		return err
	}

	return nil
}

func CreateCluster(clusterConfigPath string) error {
	cmd := exec.Command("eksctl", "create", "cluster", "-f", clusterConfigPath)

	if err := cmd.Run(); err != nil {
		return err
	}

	// TODO(RedbackThomson): Add logic to verify that the cluster was created
	// correctly and can be accessed locally.

	return nil
}

func SetupIRSA(service string) error {
	cmd := exec.Command("eksctl",
		"utils",
		"associate-iam-oidc-provider",
		"--cluster",
		defaultClusterName,
		"--approve",
	)
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("eksctl",
		"create",
		"iamserviceaccount",
		"--cluster",
		defaultClusterName,
		"--namespace",
		defaultControllerInstallNamespace,
		"--name",
		GetDefaultServiceAccountName(service),
		"--attach-policy-arn",
		// TODO(RedbackThomson): Load recommended policy from the controller
		// repository
		"arn:aws:iam::aws:policy/PowerUserAccess",
		"--approve",
	)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func GetDefaultServiceAccountName(service string) string {
	return fmt.Sprintf(defaultServiceAccountNameFmt, service)
}

func getCluster() (*eks.Cluster, error) {
	input := &eks.DescribeClusterInput{
		Name: &defaultClusterName,
	}

	sess := util.NewSession()
	client := eks.New(sess)

	out, err := client.DescribeCluster(input)
	if err != nil {
		awsErr, _ := err.(awserr.Error)

		if strings.HasPrefix(awsErr.Code(), "ResourceNotFoundException") {
			return nil, nil
		}

		return nil, err
	}

	return out.Cluster, err
}
