package repository

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/mfojtik/deptrack/pkg/managers"
	"github.com/mfojtik/deptrack/pkg/release"
)

type DependencyStatus struct {
	Dependency release.Dependency `json:"dependency"`

	Branch string `json:"branch"`

	// Level that is currently used in repo
	Current string `json:"current"`

	// Level that is the HEAD in dependency repo
	Desired string `json:"desired"`

	// Number of missing commits
	MissingCommits []string `json:"commits"`
}

// DependencyStatusFor return a status object describing the state of the given dependency.
func (r *Repository) DependencyStatusFor(targetDependency release.Dependency) (*DependencyStatus, error) {
	manager, err := managers.GetRepositoryDependencies(release.ComponentShortName(r.Component), r.Component.Branch)
	if err != nil {
		return nil, err
	}

	var status *DependencyStatus

	for _, d := range manager.Dependencies {
		if !strings.HasPrefix(string(targetDependency), d.Name+"@") {
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
		return nil, fmt.Errorf("unable to find %q in %q repository package manifest", string(targetDependency), r.Component.Name)
	}

	repoName, branchName := splitDependencyName(status.Dependency)

	dependency := GetRepository(branchName, repoName)
	if dependency == nil {
		return nil, fmt.Errorf("repository %s@%s is not cached", status.Dependency, status.Branch)
	}

	head, err := dependency.repo.Head()
	if err != nil {
		return nil, fmt.Errorf("git head: %v", err)
	}
	status.Desired = head.Hash().String()[0:8]
	status.MissingCommits, err = gitGetMissingCommits(dependency.RepoDir, status.Desired, status.Current)
	fmt.Printf("  -> Dependency %s is currently %s and desired %s (%d commits behind) ...\n", targetDependency, status.Current, status.Desired, len(status.MissingCommits))

	return status, nil
}

func splitDependencyName(d release.Dependency) (string, string) {
	parts := strings.Split(string(d), "@")
	return parts[0], parts[1]
}

func gitGetMissingCommits(root string, desired, current string) ([]string, error) {
	command := exec.Command("git", "rev-list", "--no-merges", desired, "^"+current)
	command.Dir = root
	out, err := command.CombinedOutput()
	if err != nil {
		return nil, err
	}
	var commits []string
	for _, c := range strings.Split(string(out), "\n") {
		if len(c) > 0 {
			commits = append(commits, strings.TrimSpace(c))
		}
	}
	return commits, nil
}
