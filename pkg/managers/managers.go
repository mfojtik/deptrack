package managers

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mfojtik/deptrack/pkg/managers/dep"
	"github.com/mfojtik/deptrack/pkg/managers/glide"
	"github.com/mfojtik/deptrack/pkg/managers/version"
	"github.com/mfojtik/deptrack/pkg/managers/vgo"
)

type Repository struct {
	Name       string
	BranchName string
}

type PackageManagerType string

var (
	PackageManagerTypeGlide  PackageManagerType = "glide"
	PackageManagerTypeGodeps PackageManagerType = "godeps"
	PackageManagerTypeDep    PackageManagerType = "dep"
	PackageManagerTypeGoMod  PackageManagerType = "gomod"
)

type RepositoryPackageManager struct {
	Manifests    map[string][]byte
	ManifestType PackageManagerType

	Dependencies []version.Dependency

	*Repository
}

func GetRepositoryDependencies(name, branchName string) (*RepositoryPackageManager, error) {
	repo := &Repository{
		Name:       name,
		BranchName: branchName,
	}

	var manager *RepositoryPackageManager

	if bytes, err := httpGet(buildManifestURLs(PackageManagerTypeGoMod, repo)); err == nil {
		manager = &RepositoryPackageManager{Repository: repo, Manifests: bytes, ManifestType: PackageManagerTypeGoMod}
	}

	if manager == nil {
		if bytes, err := httpGet(buildManifestURLs(PackageManagerTypeGlide, repo)); err == nil {
			manager = &RepositoryPackageManager{Repository: repo, Manifests: bytes, ManifestType: PackageManagerTypeGlide}
		}
	}
	if manager == nil {
		if bytes, err := httpGet(buildManifestURLs(PackageManagerTypeDep, repo)); err == nil {
			manager = &RepositoryPackageManager{Repository: repo, Manifests: bytes, ManifestType: PackageManagerTypeDep}
		}
	}
	if manager == nil {
		if bytes, err := httpGet(buildManifestURLs(PackageManagerTypeGodeps, repo)); err == nil {
			manager = &RepositoryPackageManager{Repository: repo, Manifests: bytes, ManifestType: PackageManagerTypeGodeps}
		}
	}
	if manager != nil {
		var err error
		manager.Dependencies, err = manager.fetchPackageManagerManifest()
		if err != nil {
			return nil, err
		}
	}

	return manager, nil
}

func (r *RepositoryPackageManager) fetchPackageManagerManifest() ([]version.Dependency, error) {
	var (
		err  error
		deps []version.Dependency
	)
	switch r.ManifestType {
	case PackageManagerTypeGoMod:
		deps, err = vgo.ParseManifest(r.Manifests)
	case PackageManagerTypeGlide:
		deps, err = glide.ParseManifest(r.Manifests)
	case PackageManagerTypeDep:
		deps, err = dep.ParseManifest(r.Manifests)
	default:
		deps = []version.Dependency{}
	}
	return deps, err
}

func httpGet(urls []string) (map[string][]byte, error) {
	result := map[string][]byte{}
	for _, u := range urls {
		parts := strings.Split(u, "/")
		name := parts[len(parts)-1]
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		response, err := client.Get(u)
		if err != nil {
			return nil, err
		}

		if response.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to get %q, HTTP code: %d", u, response.StatusCode)
		}
		result[name], err = ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		if err := response.Body.Close(); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func getGithubRawURL(repoName string, repoBranch string) string {
	return "https://raw.githubusercontent.com/" + repoName + "/" + repoBranch + "/"
}

func buildManifestURLs(manifestType PackageManagerType, repo *Repository) []string {
	switch manifestType {
	case PackageManagerTypeGlide:
		return []string{
			getGithubRawURL(repo.Name, repo.BranchName) + "glide.yaml",
			getGithubRawURL(repo.Name, repo.BranchName) + "glide.lock",
		}
	case PackageManagerTypeDep:
		return []string{
			getGithubRawURL(repo.Name, repo.BranchName) + "Gopkg.lock",
		}
	case PackageManagerTypeGoMod:
		return []string{getGithubRawURL(repo.Name, repo.BranchName) + "go.mod"}
	case PackageManagerTypeGodeps:
		return []string{getGithubRawURL(repo.Name, repo.BranchName) + "Godeps/Godeps.json"}
	default:
		panic("unknown type")
	}
}
