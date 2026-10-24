package server

// OpenRouter's official PKCE flow exchanges a one-use code for a user-owned
// API credential. Only the backend handles that credential; existing profile
// persistence and activation are reused without changing the database schema.
import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/Autumn-27/artex/db"
)

type llmOAuthAttempt struct {
	owner                 [32]byte
	verifier, name, model string
	expires               time.Time
}

type llmOAuthState struct {
	mu       sync.Mutex
	attempts map[string]llmOAuthAttempt
}

func oauthHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
}

func oauthSession(r *http.Request) [32]byte {
	return sha256.Sum256([]byte(r.Header.Get("Authorization")))
}

// Refuse arbitrary callback destinations, including userinfo and relative URLs.
// The authenticated same-origin browser supplies Origin on this POST.
func oauthOrigin(r *http.Request, raw string) (string, bool) {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || u.Path != "" {
		return "", false
	}
	if u.Scheme != "https" && !(u.Scheme == "http" && u.Hostname() == "localhost") {
		return "", false
	}
	if r.Header.Get("Origin") != raw || !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
		return "", false
	}
	return raw, true
}

func (s *Server) startOpenRouterOAuth(w http.ResponseWriter, r *http.Request) {
	oauthHeaders(w)
	var req struct{ Origin, Name, Model string }
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		writeErr(w, 400, "oauth_invalid_request")
		return
	}
	origin, ok := oauthOrigin(r, req.Origin)
	if !ok || strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Model) == "" || len(req.Name) > 120 || len(req.Model) > 200 {
		writeErr(w, 400, "oauth_invalid_request")
		return
	}
	state, err := randomToken(32)
	if err != nil {
		writeErr(w, 500, "oauth_unavailable")
		return
	}
	verifier, err := randomToken(48)
	if err != nil {
		writeErr(w, 500, "oauth_unavailable")
		return
	}
	s.llmOAuth.mu.Lock()
	if s.llmOAuth.attempts == nil {
		s.llmOAuth.attempts = make(map[string]llmOAuthAttempt)
	}
	for key, attempt := range s.llmOAuth.attempts {
		if time.Now().After(attempt.expires) {
			delete(s.llmOAuth.attempts, key)
		}
	}
	if len(s.llmOAuth.attempts) >= 32 {
		s.llmOAuth.mu.Unlock()
		writeErr(w, 429, "oauth_try_again")
		return
	}
	s.llmOAuth.attempts[state] = llmOAuthAttempt{owner: oauthSession(r), verifier: verifier, name: strings.TrimSpace(req.Name), model: strings.TrimSpace(req.Model), expires: time.Now().Add(5 * time.Minute)}
	s.llmOAuth.mu.Unlock()
	challenge := sha256.Sum256([]byte(verifier))
	query := url.Values{
		"callback_url":          {origin + "/system/llm/oauth/callback/"},
		"code_challenge":        {base64.RawURLEncoding.EncodeToString(challenge[:])},
		"code_challenge_method": {"S256"}, "state": {state}, "key_label": {"ARTEX Community"},
	}
	writeJSON(w, 200, map[string]string{"url": "https://openrouter.ai/auth?" + query.Encode()})
}

func (s *Server) completeOpenRouterOAuth(w http.ResponseWriter, r *http.Request) {
	oauthHeaders(w)
	var req struct{ Code, State string }
	r.Body = http.MaxBytesReader(w, r.Body, 8192)
	if json.NewDecoder(r.Body).Decode(&req) != nil || req.Code == "" || len(req.Code) > 4096 {
		writeErr(w, 400, "oauth_invalid_request")
		return
	}
	s.llmOAuth.mu.Lock()
	attempt, ok := s.llmOAuth.attempts[req.State]
	if ok && attempt.owner == oauthSession(r) {
		delete(s.llmOAuth.attempts, req.State)
	} else {
		ok = false
	}
	s.llmOAuth.mu.Unlock()
	if !ok || time.Now().After(attempt.expires) {
		writeErr(w, 400, "oauth_expired")
		return
	}
	pg := s.pg(w)
	if pg == nil {
		return
	}
	body, _ := json.Marshal(map[string]string{"code": req.Code, "code_verifier": attempt.verifier, "code_challenge_method": "S256"})
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	exchange, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://openrouter.ai/api/v1/auth/keys", bytes.NewReader(body))
	if err != nil {
		writeErr(w, 500, "oauth_unavailable")
		return
	}
	exchange.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 20 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	response, err := client.Do(exchange)
	if err != nil {
		writeErr(w, 502, "oauth_exchange_failed")
		return
	}
	defer response.Body.Close()
	var result struct {
		Key string `json:"key"`
	}
	if response.StatusCode != http.StatusOK || json.NewDecoder(io.LimitReader(response.Body, 65536)).Decode(&result) != nil || result.Key == "" {
		// Never return the upstream body: it may contain credential material.
		writeErr(w, 502, "oauth_exchange_failed")
		return
	}
	profile := db.LLMProfile{Name: attempt.name, Format: "openai", BaseURL: "https://openrouter.ai/api/v1", Model: attempt.model, APIKey: result.Key, Streaming: true, PoolExclude: true}
	id, err := pg.SaveProfile(&profile)
	if err != nil {
		writeErr(w, 500, "oauth_save_failed")
		return
	}
	s.invalidateProfileAgents()
	s.reapplyActiveProfile()
	writeJSON(w, 200, map[string]any{"id": id})
}
