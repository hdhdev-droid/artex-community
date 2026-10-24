package notify

import "github.com/Autumn-27/artex/internal/l10n"

func text(literal string) string {
	return l10n.Template(labelLang(), literal)
}
