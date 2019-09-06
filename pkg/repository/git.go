package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"

	"github.com/mfojtik/deptrack/pkg/release"
)

var repositories = []*Repository{}

type Repository struct {
	Component release.Component

	RootDir string
	RepoDir string

	repo *git.Repository
}

func ForComponent(rootDir string, component release.Component, components []release.Component) (*Repository, error) {
	for _, r := range repositories {
		if reflect.DeepEqual(r.Component, component) {
			return r, nil
		}
	}
	repository := &Repository{Component: component, RootDir: rootDir}
	if err := repository.Synchronize(components); err != nil {
		return nil, err
	}
	fmt.Printf("* Processing %s ...\n", component.Name)
	(&sync.Once{}).Do(func() {
		repositories = append(repositories, repository)
	})
	return repository, nil
}

func GetRepository(branchName, name string) *Repository {
	var result *Repository
	(&sync.Once{}).Do(func() {
		for _, r := range repositories {
			if r.Component.Name == name && r.Component.Branch == branchName {
				result = r
				break
			}
		}
	})
	return result
}

func componentNameToFileName(name string) string {
	return strings.Replace(name, "/", "_", -1)
}

func (r *Repository) Synchronize(components []release.Component) error {
	if release.IsLeaf(r.Component, components) {
		return nil
	}
	r.RepoDir = filepath.Join(r.RootDir, componentNameToFileName(r.Component.Name)+"_"+r.Component.Branch)

	if _, err := os.Stat(r.RepoDir); !os.IsNotExist(err) {
		r.repo, err = git.PlainOpen(r.RepoDir)
		if err != nil {
			return fmt.Errorf("unable to open repository %q: %v", "https://"+r.Component.Name, err)
		}
	}
	if r.repo != nil {
		if err := r.repo.Fetch(&git.FetchOptions{}); err != nil && err != git.NoErrAlreadyUpToDate {
			return fmt.Errorf("fetch error: %v", err)
		}
		return nil
	}
	var err error
	r.repo, err = git.PlainClone(r.RepoDir, true, &git.CloneOptions{
		URL:           "https://" + r.Component.Name,
		ReferenceName: plumbing.NewBranchReferenceName(r.Component.Branch),
		SingleBranch:  true,
	})
	if err != nil {
		return fmt.Errorf("unable to clone repository %q: %v", "https://"+r.Component.Name, err)
	}
	return nil
}
