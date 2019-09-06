package dep

import (
	"strings"

	"github.com/mfojtik/deptrack/pkg/managers/version"
)

func ParseManifest(manifest map[string][]byte) ([]version.Dependency, error) {
	lock, err := readLock(manifest["Gopkg.lock"])
	if err != nil {
		return nil, err
	}
	list := []version.Dependency{}
	for _, p := range lock.Projects {
		digest := strings.TrimPrefix(p.Digest, "1:")
		if len(digest) > 8 {
			digest = digest[0:8]
		}
		list = append(list, version.Dependency{
			Name:       p.Name,
			Version:    p.Version,
			Digest:     digest,
			Repository: p.Source,
		})
	}
	return list, nil
}
