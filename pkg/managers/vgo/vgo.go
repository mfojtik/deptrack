package vgo

import (
	"github.com/mfojtik/deptrack/pkg/managers/version"
	"github.com/mfojtik/deptrack/pkg/managers/vgo/modfile"
)

func ParseManifest(manifest map[string][]byte) ([]version.Dependency, error) {
	f, err := modfile.Parse("", manifest["go.mod"], nil)
	if err != nil {
		return nil, err
	}
	list := []version.Dependency{}
	for _, r := range f.Require {
		list = append(list, version.Dependency{
			Name:    r.Mod.Path,
			Version: r.Mod.Version,
		})
	}

	return list, nil
}
