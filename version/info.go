// Package version records the community release target and its upstream provenance.
package version

import (
	_ "embed"
	"encoding/json"
	"strings"
)

//go:embed VERSION
var declaredVersion string

//go:embed upstream.json
var upstreamJSON []byte

type Upstream struct {
	Repository         string `json:"repository"`
	Version            string `json:"version"`
	Commit             string `json:"commit"`
	Describe           string `json:"describe"`
	OriginalRepository string `json:"original_repository"`
	OriginalVersion    string `json:"original_version"`
	OriginalEvidence   string `json:"original_evidence"`
}

type Info struct {
	Version  string   `json:"version"`
	Upstream Upstream `json:"upstream"`
}

func Version() string { return strings.TrimSpace(declaredVersion) }

func Metadata() Info {
	var upstream Upstream
	if err := json.Unmarshal(upstreamJSON, &upstream); err != nil {
		panic("invalid embedded upstream metadata: " + err.Error())
	}
	return Info{Version: Version(), Upstream: upstream}
}
