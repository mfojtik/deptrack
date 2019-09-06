package report

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/mfojtik/deptrack/pkg/release"
	"github.com/mfojtik/deptrack/pkg/repository"
)

type ReleaseStatus struct {
	Name string `yaml:"name",json:"name"`

	Components       []release.Component `yaml:"components",json:"components"`
	ComponentsStatus []ComponentStatus   `yaml:"status",json:"status"`
}

type ComponentStatus struct {
	Component release.Component             `json:"component"`
	Status    []repository.DependencyStatus `json:"status"`
}

func WriteReport(targetDir string, status ReleaseStatus) error {
	jsonBytes, err := json.Marshal(status)
	if err != nil {
		return err
	}
	yamlBytes, err := yaml.Marshal(status)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(targetDir, fmt.Sprintf("%s_report.yaml", status.Name)), yamlBytes, os.ModePerm); err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(targetDir, fmt.Sprintf("%s_report.json", status.Name)), jsonBytes, os.ModePerm)
}

func NewReport(gitDir string, r release.Release) (*ReleaseStatus, error) {
	releaseStatus := ReleaseStatus{
		Name:             r.Name,
		Components:       r.Components,
		ComponentsStatus: []ComponentStatus{},
	}
	for _, c := range r.Components {
		repo, err := repository.ForComponent(gitDir, c, r.Components)
		if err != nil {
			return nil, err
		}
		status := []repository.DependencyStatus{}
		for _, d := range c.Vendor {
			dependencyStatus, err := repo.DependencyStatusFor(d)
			if err != nil {
				return nil, err
			}
			status = append(status, *dependencyStatus)
		}
		releaseStatus.ComponentsStatus = append(releaseStatus.ComponentsStatus, ComponentStatus{
			Component: c,
			Status:    status,
		})
	}
	return &releaseStatus, nil
}
