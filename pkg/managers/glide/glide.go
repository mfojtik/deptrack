package glide

import (
	"github.com/Masterminds/glide/cfg"
	"github.com/mfojtik/deptrack/pkg/managers/version"
)

func ParseManifest(manifest map[string][]byte) ([]version.Dependency, error) {
	lockFile, err := cfg.LockfileFromYaml(manifest["glide.lock"])
	if err != nil {
		return nil, err
	}
	configFile, err := cfg.ConfigFromYaml(manifest["glide.yaml"])
	if err != nil {
		return nil, err
	}

	list := []version.Dependency{}
	for _, i := range lockFile.Imports {
		v := ""
		for _, c := range configFile.Imports {
			if c.Name == i.Name {
				v = c.Reference
				break
			}
		}
		list = append(list, version.Dependency{
			Name:       i.Name,
			Digest:     i.Version[0:8],
			Version:    v,
			Repository: i.Repository,
		})
	}
	return list, nil
}
