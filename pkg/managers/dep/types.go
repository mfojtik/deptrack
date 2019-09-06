package dep

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type rawLock struct {
	Projects []rawLockedProject `toml:"projects"`
}

type rawLockedProject struct {
	Name      string   `toml:"name"`
	Branch    string   `toml:"branch,omitempty"`
	Revision  string   `toml:"revision"`
	Version   string   `toml:"version,omitempty"`
	Source    string   `toml:"source,omitempty"`
	Packages  []string `toml:"packages"`
	PruneOpts string   `toml:"pruneopts"`
	Digest    string   `toml:"digest"`
}

func readLock(lockBytes []byte) (*rawLock, error) {
	raw := rawLock{}
	err := toml.Unmarshal(lockBytes, &raw)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse the lock as TOML")
	}

	return &raw, nil
}
