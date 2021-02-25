package crd

import (
	"context"
	"io/ioutil"
	"path/filepath"

	apiextensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensiontypesv1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/aws-controllers-k8s/dev-tools/pkg/repository"
)

const (
	defaultCRDsPath = "config/crd/bases"
)

func New(rl repository.Loader, kubeconfig string) (*Manager, error) {
	aev1CRDClient, err := NewClient(kubeconfig)
	if err != nil {
		return nil, err
	}
	return &Manager{
		rl:        rl,
		crdClient: aev1CRDClient,
	}, nil
}

type Manager struct {
	rl        repository.Loader
	crdClient apiextensiontypesv1.CustomResourceDefinitionInterface
}

func (m *Manager) getCRD(ctx context.Context, name string) (bool, *apiextensionv1.CustomResourceDefinition, error) {
	crd, err := m.crdClient.Get(ctx, name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return false, nil, nil
	}
	if err != nil {
		return false, nil, err
	}
	return true, crd, nil
}

func (m *Manager) createOrUpdateCRD(ctx context.Context, crd *apiextensionv1.CustomResourceDefinition) error {
	exist, gcrd, err := m.getCRD(ctx, crd.Name)
	if err != nil {
		return err
	}

	if exist {
		crd.ObjectMeta.ResourceVersion = gcrd.ObjectMeta.ResourceVersion
		_, err := m.crdClient.Update(ctx, crd, metav1.UpdateOptions{})
		return err
	}

	_, err = m.crdClient.Create(ctx, crd, metav1.CreateOptions{})
	return err
}

func (m *Manager) DeployCRD(ctx context.Context, crd *apiextensionv1.CustomResourceDefinition) error {
	return m.createOrUpdateCRD(ctx, crd)
}

func (m *Manager) DeployRepositoryCRDs(ctx context.Context, repositoryName string) error {
	repo, err := m.rl.LoadRepository(repositoryName, repository.RepositoryTypeController)
	if err != nil {
		return err
	}

	crdDirectory := filepath.Join(repo.FullPath, defaultCRDsPath)
	crdFiles, err := ioutil.ReadDir(crdDirectory)
	if err != nil {
		return err
	}

	for _, crdFile := range crdFiles {
		filePath := filepath.Join(crdDirectory, crdFile.Name())
		crd, err := GetCRDFromFile(filePath)
		if err != nil {
			return err
		}
		return m.DeployCRD(ctx, crd)
	}
	return nil
}
