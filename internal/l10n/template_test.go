package l10n

import (
	"encoding/json"
	"reflect"
	"regexp"
	"testing"
)

func TestTemplatesPreserveFormatArguments(t *testing.T) {
	var rows [][]string
	if err := json.Unmarshal(source, &rows); err != nil {
		t.Fatal(err)
	}
	verbs := regexp.MustCompile(`%[+#0 -]*[0-9.]*[a-zA-Z%]`)
	seen := map[string]bool{}
	for _, row := range rows {
		if seen[row[0]] {
			t.Fatalf("duplicate source: %q", row[0])
		}
		seen[row[0]] = true
		for _, value := range row[1:] {
			if !reflect.DeepEqual(verbs.FindAllString(row[0], -1), verbs.FindAllString(value, -1)) {
				t.Fatalf("format arguments differ: %q", row[0])
			}
		}
	}
}
