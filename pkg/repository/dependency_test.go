package repository

import (
	"log"
	"testing"

	"github.com/mfojtik/deptrack/pkg/release"
)

func TestDependencyStatusFor(t *testing.T) {
	if _, err := ForComponent(release.Component{
		Name:   "github.com/openshift/api",
		Branch: "master",
	}); err != nil {
		t.Fatal(err)
	}
	repo, err := ForComponent(release.Component{
		Name:   "github.com/openshift/client-go",
		Branch: "master",
		Vendor: []release.Dependency{"github.com/openshift/api"},
	})
	if err != nil {
		t.Fatal(err)
	}

	status, err := repo.DependencyStatusFor("github.com/openshift/api")
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("status: %#v", status.MissingCommits)
}
