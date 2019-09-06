package repository

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/mfojtik/deptrack/pkg/release"
)

func TestDependencyStatusFor(t *testing.T) {
	components := []release.Component{
		{
			Name:   "github.com/openshift/api",
			Branch: "master",
		},
		{
			Name:   "github.com/openshift/client-go",
			Branch: "master",
			Vendor: []release.Dependency{"github.com/openshift/api"},
		},
	}
	tmpDir, _ := ioutil.TempDir("", "test")
	t.Logf("Using %q as root directory", tmpDir)
	if _, err := ForComponent(tmpDir, components[0], components); err != nil {
		t.Fatal(err)
	}
	repo, err := ForComponent(tmpDir, components[1], components)
	if err != nil {
		t.Fatal(err)
	}

	status, err := repo.DependencyStatusFor("github.com/openshift/api")
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("status: %#v", status.MissingCommits)
}
