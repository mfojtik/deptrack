package repository

import (
	"fmt"
	"io/ioutil"
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
	repo    *git.Repository
}

func ForComponent(component release.Component) (*Repository, error) {
	for _, r := range repositories {
		if reflect.DeepEqual(r.Component, component) {
			return r, nil
		}
	}
	repository := &Repository{Component: component}
	if err := repository.Synchronize(); err != nil {
		return nil, err
	}
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

func (r *Repository) Synchronize() error {
	if r.repo != nil {
		if err := r.repo.Fetch(&git.FetchOptions{}); err != nil {
			return err
		}
		return nil
	}
	var err error
	r.RootDir, err = ioutil.TempDir("", componentNameToFileName(r.Component.Name)+"_"+r.Component.Branch)
	if err != nil {
		return err
	}
	r.repo, err = git.PlainClone(r.RootDir, true, &git.CloneOptions{
		URL:           "https://" + r.Component.Name,
		ReferenceName: plumbing.NewBranchReferenceName(r.Component.Branch),
		SingleBranch:  true,
	})
	if err != nil {
		return fmt.Errorf("unable to clone repository %q: %v", "https://"+r.Component.Name, err)
	}
	return nil
}
