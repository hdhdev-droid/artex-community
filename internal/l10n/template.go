// Package l10n translates compile-time document templates. Never pass user text,
// evidence, model output or an already formatted message to Template.
package l10n

import (
	_ "embed"
	"encoding/json"
	"sort"
	"strings"
)

//go:embed templates.json
var source []byte

var replacers = build()

// Language is set by the host, keeping validation packages independent of the DB.
var Language func() string

func T(literal string) string {
	if Language != nil {
		return Template(Language(), literal)
	}
	return literal
}

func build() map[string]*strings.Replacer {
	var rows [][]string
	if err := json.Unmarshal(source, &rows); err != nil {
		panic(err)
	}
	sort.SliceStable(rows, func(i, j int) bool { return len(rows[i][0]) > len(rows[j][0]) })
	out := make(map[string]*strings.Replacer)
	for index, lang := range []string{"en", "zh", "es", "ja"} {
		var pairs []string
		for _, row := range rows {
			if len(row) != 5 {
				panic("invalid document translation row")
			}
			pairs = append(pairs, row[0], row[index+1])
		}
		out[lang] = strings.NewReplacer(pairs...)
	}
	return out
}

func Template(lang, literal string) string {
	if r := replacers[lang]; r != nil {
		return r.Replace(literal)
	}
	return literal
}
