package report

import (
	"strings"
	"testing"
	"time"

	"github.com/Autumn-27/artex/db"
)

func TestLocalizedReportPreservesEvidence(t *testing.T) {
	previous := Lang
	t.Cleanup(func() { Lang = previous })
	for _, lang := range []string{"en", "zh", "es", "ja"} {
		Lang = func() string { return lang }
		got := FindingsMarkdown([]*db.DBFinding{{Name: "사용자 제목", Evidence: "원본 증거", Report: "사용자 보고서"}}, time.Now())
		for _, original := range []string{"사용자 제목", "원본 증거", "사용자 보고서"} {
			if !strings.Contains(got, original) {
				t.Fatalf("%s: lost original content", lang)
			}
		}
		if strings.Contains(got, "취약점 요약 보고서") || strings.Contains(got, "생성 시각") {
			t.Fatalf("%s: untranslated report template", lang)
		}
	}
}
