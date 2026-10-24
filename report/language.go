package report

import "github.com/Autumn-27/artex/internal/l10n"

// Lang is wired by the server. Standalone exports retain their historical language.
var Lang func() string

func text(literal string) string {
	lang := "ko"
	if Lang != nil {
		lang = Lang()
	}
	return l10n.Template(lang, literal)
}
