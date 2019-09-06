package version

type Dependency struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Digest     string `json:"digest"`
	Repository string `json:"repository,omitempty"`
}
