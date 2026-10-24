package notify

import (
	"strings"
	"testing"
)

func TestLocalizedBatchTitlesAndLinks(t *testing.T) {
	prev := Lang
	t.Cleanup(func() { Lang = prev })
	for _, lang := range []string{"en", "zh", "es", "ja"} {
		Lang = func() string { return lang }
		m := Message{Batch: true, HomeURL: "https://example.invalid", WindowMinutes: 10}
		title := markdownTitle(m)
		body := htmlBody(m, 0)
		if strings.Contains(title, "취약점") || strings.Contains(body, "플랫폼에서") || strings.Contains(body, "최근") {
			t.Fatalf("%s: untranslated batch template", lang)
		}
	}
}
