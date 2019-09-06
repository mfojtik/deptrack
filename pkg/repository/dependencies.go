package repository

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/mfojtik/deptrack/pkg/managers"
	"github.com/mfojtik/deptrack/pkg/release"
)

type DependencyStatus struct {
	Dependency release.Dependency

	Branch string

	// Level that is currently used in repo
	Current string

	// Level that is the HEAD in dependency repo
	Desired string

	// Number of missing commits
	MissingCommits []string
}

func (r *Repository) DependencyStatusFor(targetDependency release.Dependency) (*DependencyStatus, error) {
	manager, err := managers.GetRepositoryDependencies(release.ComponentShortName(r.Component), r.Component.Branch)
	if err != nil {
		return nil, err
	}
	var status *DependencyStatus

	for _, d := range manager.Dependencies {
		if !strings.HasPrefix(d.Name, string(targetDependency)) {
			continue
		}
		if len(d.Digest) == 0 {
			continue
		}
		status = &DependencyStatus{
			Dependency: targetDependency,
			Current:    d.Digest,
			Branch:     d.Version,
		}
	}
	if status == nil {
		return nil, fmt.Errorf("unable to find %q in repository package manifest", string(targetDependency))
	}

	dependency := GetRepository(status.Branch, string(status.Dependency))
	if dependency == nil {
		return nil, fmt.Errorf("repository %s@%s is not cached", status.Dependency, status.Branch)
	}

	head, err := dependency.repo.Head()
	if err != nil {
		return nil, err
	}
	status.Desired = head.Hash().String()[0:8]
	status.MissingCommits, err = gitGetMissingCommits(dependency.RootDir, status.Desired, status.Current)

	return status, nil
}

func gitGetMissingCommits(root string, desired, current string) ([]string, error) {
	command := exec.Command("git", "rev-list", "--no-merges", desired, "^"+current)
	command.Dir = root
	out, err := command.CombinedOutput()
	if err != nil {
		return nil, err
	}
	commits := []string{}
	for _, c := range strings.Split(string(out), "\n") {
		commits = append(commits, strings.TrimSpace(c))
	}
	return commits, nil
}
