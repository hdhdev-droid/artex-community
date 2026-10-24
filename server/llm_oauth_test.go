package server

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestOAuthOriginBoundary(t *testing.T) {
	for _, tc := range []struct {
		origin, header string
		accepted       bool
	}{
		{"http://localhost:5173", "http://localhost:5173", true},
		{"https://example.com", "https://example.com", true},
		{"http://127.0.0.1:5173", "http://127.0.0.1:5173", false},
		{"http://example.com", "http://example.com", false},
		{"http://localhost:5173/path", "http://localhost:5173/path", false},
		{"http://localhost:5173", "https://example.com", false},
		{"https://user@example.com", "https://user@example.com", false},
	} {
		r := httptest.NewRequest("POST", "/", nil)
		r.Header.Set("Origin", tc.header)
		r.Header.Set("Authorization", "Bearer test-session")
		if _, ok := oauthOrigin(r, tc.origin); ok != tc.accepted {
			t.Errorf("origin acceptance mismatch for %s", tc.origin)
		}
	}
}

func TestChatGPTCallbackKeepsCodeServerSide(t *testing.T) {
	a := &chatGPTAttempt{origin: "http://localhost:5173", clientID: "dynamic_agent_client", expires: time.Now().Add(time.Minute)}
	s := &Server{chatGPT: chatGPTOAuth{attempts: map[string]*chatGPTAttempt{"test-state": a}}}
	r := httptest.NewRequest("GET", "http://127.0.0.1:8789/auth/callback?state=test-state&code=test-code&client_id=oaiapp_test", nil)
	w := httptest.NewRecorder()
	s.chatGPTCallback(w, r)
	if w.Code != 303 || w.Header().Get("Location") != "http://localhost:5173/system/llm/chatgpt/callback/?attempt=test-state" {
		t.Fatal("incorrect local callback redirect")
	}
	if strings.Contains(w.Body.String(), "test-code") || a.code != "test-code" || a.clientID != "oaiapp_test" || !a.returned || a.failed {
		t.Fatal("callback code leaked or registration not retained")
	}
	w = httptest.NewRecorder()
	s.chatGPTCallback(w, r)
	if w.Code != 400 {
		t.Fatal("replayed callback accepted")
	}
}

func TestChatGPTCallbackRejectsWrongHostAndExpiredState(t *testing.T) {
	a := &chatGPTAttempt{origin: "http://localhost:5173", expires: time.Now().Add(-time.Minute)}
	s := &Server{chatGPT: chatGPTOAuth{attempts: map[string]*chatGPTAttempt{"expired": a}}}
	for _, tc := range []struct {
		url    string
		status int
	}{
		{"http://localhost:8789/auth/callback?state=expired", 404},
		{"http://127.0.0.1:8789/auth/callback?state=expired", 400},
		{"http://127.0.0.1:8789/auth/callback?state=unknown", 400},
	} {
		w := httptest.NewRecorder()
		s.chatGPTCallback(w, httptest.NewRequest("GET", tc.url, nil))
		if w.Code != tc.status {
			t.Errorf("callback status: got %d, want %d", w.Code, tc.status)
		}
	}
	if a.returned {
		t.Fatal("expired state consumed as a valid callback")
	}
}

func TestOpenRouterCompletionRejectsAnotherSession(t *testing.T) {
	owner := httptest.NewRequest("POST", "/", nil)
	owner.Header.Set("Authorization", "Bearer owner")
	s := &Server{llmOAuth: llmOAuthState{attempts: map[string]llmOAuthAttempt{"test-state": {owner: oauthSession(owner), expires: time.Now().Add(time.Minute)}}}}
	r := httptest.NewRequest("POST", "/", strings.NewReader(`{"code":"test-code","state":"test-state"}`))
	r.Header.Set("Authorization", "Bearer other-session")
	w := httptest.NewRecorder()
	s.completeOpenRouterOAuth(w, r)
	if w.Code != 400 || len(s.llmOAuth.attempts) != 1 {
		t.Fatal("foreign session consumed an OAuth attempt")
	}
}
