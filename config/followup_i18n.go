package config

import (
	_ "embed"
	"encoding/json"
)

// Only compile-time system templates belong here. Format parameters, evidence
// and user-authored values must be passed after translation.
//
//go:embed messages_followup.json
var followupSource []byte

var followupMessages = func() map[string]map[string]string {
	var rows [][]string
	if err := json.Unmarshal(followupSource, &rows); err != nil {
		panic(err)
	}
	out := make(map[string]map[string]string, len(rows))
	for _, row := range rows {
		if len(row) != 6 {
			panic("invalid system translation row")
		}
		out[row[0]] = map[string]string{
			"ko": row[1], "en": row[2], "zh": row[3], "es": row[4], "ja": row[5],
		}
	}
	return out
}()

func followupText(source, lang string) string {
	if text := followupMessages[source][lang]; text != "" {
		return text
	}
	return source
}
