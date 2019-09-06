package release

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func ReadConfig(configFile string) (*Release, error) {
	configBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	config := Release{}
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
