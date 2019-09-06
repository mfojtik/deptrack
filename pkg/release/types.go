package release

import (
	"strings"
)

// Release describe a single OpenShift release (eg. release-4.2)
type Release struct {
	Name string `yaml:"name",json:"name"`

	// Components is a list of components we want to track for the release
	Components []Component `yaml:"components",json:"components"`
}

type Dependency string

// Component describe a single component included in the release (eg. "openshift/library-go")
type Component struct {
	Name string `yaml:"name",json:"name"`

	// Branch is the name of the Git branch that matches the release
	Branch string `yaml:"branch",json:"branch"`

	// Vendor is a list of dependencies the component vendor that we are interested to track changes in.
	Vendor []Dependency `yaml:"vendor",json:"vendor"`
}

func ComponentShortName(component Component) string {
	return strings.TrimPrefix(component.Name, "github.com/")
}

type Config struct {
	Releases []Release `yaml:"releases"`
}
