package agent

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/Autumn-27/norma/llm"
)

var investigationSecretAssignment = regexp.MustCompile(`(?im)(["']?[A-Za-z_][A-Za-z0-9_]*_KEY["']?\s*[:=]\s*)(?:"(?:\\.|[^"\\])*"|'[^']*'|[^\s,;}\r\n]+)`)
var investigationBearer = regexp.MustCompile(`(?i)\bBearer\s+[A-Za-z0-9._~+/=-]+`)

// Sanitize a copy at the provider boundary, including embedded JSON tool results.
// Local original evidence and hashes remain unchanged. Credentials for the log
// source never enter the request in the first place.
func redactInvestigationRequest(req *llm.CompletionRequest) error {
	raw, err := json.Marshal(req)
	if err != nil {
		return err
	}
	var tree any
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	if err := decoder.Decode(&tree); err != nil {
		return err
	}
	var clean func(any) any
	clean = func(value any) any {
		switch v := value.(type) {
		case map[string]any:
			for key, child := range v {
				lower := strings.ToLower(key)
				if strings.HasSuffix(lower, "_key") || lower == "password" || lower == "authorization" || lower == "cookie" {
					v[key] = "[REDACTED]"
				} else {
					v[key] = clean(child)
				}
			}
		case []any:
			for i, child := range v {
				v[i] = clean(child)
			}
		case string:
			var embedded any
			trimmed := strings.TrimSpace(v)
			decoder := json.NewDecoder(strings.NewReader(v))
			decoder.UseNumber()
			if (strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[")) && json.Valid([]byte(v)) && decoder.Decode(&embedded) == nil {
				data, _ := json.Marshal(clean(embedded))
				return string(data)
			}
			v = investigationSecretAssignment.ReplaceAllString(v, "${1}[REDACTED]")
			return investigationBearer.ReplaceAllString(v, "Bearer [REDACTED]")
		}
		return value
	}
	raw, err = json.Marshal(clean(tree))
	if err != nil {
		return err
	}
	var sanitized llm.CompletionRequest
	if err := json.Unmarshal(raw, &sanitized); err != nil {
		return err
	}
	*req = sanitized
	return nil
}
