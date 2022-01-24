package soak

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gofrs/flock"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
)

// InstallLocalChart installs a Helm chart from a local path.
func InstallLocalChart(
	chartPath string,
	namespace string,
	releaseName string,
	values map[string]interface{},
) (*release.Release, error) {
	actionConfig, err := getHelmConfig(namespace)
	if err != nil {
		return nil, err
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		return nil, err
	}

	client := action.NewInstall(actionConfig)
	client.Namespace = namespace
	client.CreateNamespace = true
	if releaseName != "" {
		client.ReleaseName = releaseName
	} else {
		client.GenerateName = true
	}

	release, err := client.Run(chart, values)
	if err != nil {
		return nil, err
	}

	return release, nil
}

// InstallRepoChart installs a Helm chart from a remote chart repository.
func InstallRepoChart(
	repo string,
	chart string,
	namespace string,
	releaseName string,
	values map[string]interface{},
) (*release.Release, error) {
	settings := cli.New()
	actionConfig, err := getHelmConfig(namespace)
	if err != nil {
		return nil, err
	}

	client := action.NewInstall(actionConfig)

	cp, err := client.ChartPathOptions.LocateChart(fmt.Sprintf("%s/%s", repo, chart), settings)
	if err != nil {
		log.Fatal(err)
	}

	return InstallLocalChart(cp, namespace, releaseName, values)
}

// AddHelmRepo adds a remote chart repository and updates the indexes.
func AddHelmRepo(name string, url string) error {
	settings := cli.New()
	repoFile := settings.RepositoryConfig

	//Ensure the file directory exists as it is required for file locking
	err := os.MkdirAll(filepath.Dir(repoFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	// Acquire a file lock for process synchronization
	fileLock := flock.New(strings.Replace(repoFile, filepath.Ext(repoFile), ".lock", 1))
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer fileLock.Unlock()
	}
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(repoFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var f repo.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		return err
	}

	if f.Has(name) {
		return nil
	}

	c := repo.Entry{
		Name: name,
		URL:  url,
	}

	r, err := repo.NewChartRepository(&c, getter.All(settings))
	if err != nil {
		return err
	}

	if _, err := r.DownloadIndexFile(); err != nil {
		err := errors.Wrapf(err, "looks like %q is not a valid chart repository or cannot be reached", url)
		return err
	}

	f.Update(&c)

	if err := f.WriteFile(repoFile, 0644); err != nil {
		return err
	}

	// Ensure repository indexes are updated
	if err = repoUpdate(); err != nil {
		return err
	}

	return nil
}

func repoUpdate() error {
	settings := cli.New()
	repoFile := settings.RepositoryConfig

	f, err := repo.LoadFile(repoFile)
	if os.IsNotExist(errors.Cause(err)) || len(f.Repositories) == 0 {
		return errors.New("no repositories found. You must add one before updating")
	}
	var repos []*repo.ChartRepository
	for _, cfg := range f.Repositories {
		r, err := repo.NewChartRepository(cfg, getter.All(settings))
		if err != nil {
			return err
		}
		repos = append(repos, r)
	}

	var wg sync.WaitGroup
	for _, re := range repos {
		wg.Add(1)
		go func(re *repo.ChartRepository) {
			defer wg.Done()
			re.DownloadIndexFile()
		}(re)
	}
	wg.Wait()

	return nil
}

func getHelmConfig(namespace string) (*action.Configuration, error) {
	settings := cli.New()

	emptyLog := func(format string, v ...interface{}) {}

	actionConfig := &action.Configuration{}
	if err := actionConfig.Init(
		settings.RESTClientGetter(),
		namespace,
		os.Getenv("HELM_DRIVER"),
		emptyLog,
	); err != nil {
		return nil, err
	}

	return actionConfig, nil
}
