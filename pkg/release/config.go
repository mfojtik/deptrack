package release

import (
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

func ReadConfig(configFile string) (*Config, error) {
	configBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	config := Config{}
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// IsLeaf is true if no other component depend on passed component
func IsLeaf(component Component, components []Component) bool {
	for _, c := range components {
		for _, d := range c.Vendor {
			if strings.HasPrefix(string(d), component.Name) {
				return false
			}
		}
	}
	return true
}
