package soak

import (
	"fmt"
	"strings"

	"github.com/aws-controllers-k8s/dev-tools/pkg/util"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecrpublic"
)

var (
	ecrPublicRegion    = "us-east-1"
	defaultRepoNameFmt = "ack-%s-soak"
)

// EnsureECRRepository ensures that the default soak test image ECR public
// repository has been created, or creates it if it doesn't exist.
func EnsureECRRepository(service string) (string, error) {
	repoName := fmt.Sprintf(defaultRepoNameFmt, service)

	if existing, err := getRepo(repoName); err != nil {
		return "", err
	} else if existing != nil {
		return *existing.RepositoryUri, nil
	}

	repoUri, err := createRepo(repoName)
	if err != nil {
		return "", err
	}

	return *repoUri, nil
}

// GetRepoURL returns the repository URL of the default soak test image ECR
// public repository if it existed, otherwise it will return an empty string.
func GetRepoURL(service string) (string, error) {
	repoName := fmt.Sprintf(defaultRepoNameFmt, service)

	if existing, err := getRepo(repoName); err != nil {
		return "", err
	} else if existing != nil {
		return *existing.RepositoryUri, nil
	}

	return "", nil
}

func createRepo(repoName string) (*string, error) {
	input := &ecrpublic.CreateRepositoryInput{
		RepositoryName: &repoName,
	}

	sess := util.NewSessionWithRegion(ecrPublicRegion)
	client := ecrpublic.New(sess)

	out, err := client.CreateRepository(input)
	if err != nil {
		return nil, err
	}

	return out.Repository.RepositoryUri, nil
}

func getRepo(repoName string) (*ecrpublic.Repository, error) {
	input := &ecrpublic.DescribeRepositoriesInput{
		RepositoryNames: []*string{&repoName},
	}

	sess := util.NewSessionWithRegion(ecrPublicRegion)
	client := ecrpublic.New(sess)

	out, err := client.DescribeRepositories(input)
	if err != nil {
		awsErr, _ := err.(awserr.Error)

		if strings.HasPrefix(awsErr.Code(), "RepositoryNotFoundException") {
			return nil, nil
		}

		return nil, err
	}

	if len(out.Repositories) == 0 {
		return nil, nil
	}

	return out.Repositories[0], nil
}
