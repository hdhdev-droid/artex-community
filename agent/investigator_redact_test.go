package agent

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Autumn-27/norma/llm"
)

func TestInvestigationRedactsSecretsWithoutChangingStoredEvidence(t *testing.T) {
	original := llm.CompletionRequest{Messages: []llm.Message{{
		Role: llm.RoleUser,
		Content: []llm.ContentBlock{llm.TextBlock(`{"CLOUD_KEY":"fixture-one","nested":{"private_key":"fixture-two"},"raw":"SERVICE_KEY=fixture-three Authorization: Bearer fixture-four","path":"/audit","id":9223372036854775806}`),
			llm.TextBlock("PRIVATE_KEY=\"fixture-five\nfixture-six\"\nnon_secret=retained")},
	}}}
	before, _ := json.Marshal(original)
	request := original
	if err := redactInvestigationRequest(&request); err != nil {
		t.Fatal(err)
	}
	result, _ := json.Marshal(request)
	for _, value := range []string{"fixture-one", "fixture-two", "fixture-three", "fixture-four", "fixture-five", "fixture-six"} {
		if strings.Contains(string(result), value) {
			t.Fatal("secret remained in outgoing investigation request")
		}
	}
	if !strings.Contains(string(result), "/audit") || !strings.Contains(string(result), "9223372036854775806") || !strings.Contains(string(result), "non_secret=retained") {
		t.Fatal("non-secret evidence was removed")
	}
	after, _ := json.Marshal(original)
	if string(before) != string(after) {
		t.Fatal("original evidence was modified")
	}
}
