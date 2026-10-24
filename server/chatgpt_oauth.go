package server

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Autumn-27/artex/db"
	"github.com/Autumn-27/artex/internal/chatgpt"
)

type chatGPTAttempt struct {
	owner                                              [32]byte
	origin, accountID, clientID, verifier, nonce, code string
	expires                                            time.Time
	returned, failed                                   bool
}
type chatGPTCatalog struct {
	models  []chatgpt.Model
	expires time.Time
}
type chatGPTOAuth struct {
	mu       sync.Mutex
	store    *chatgpt.Store
	attempts map[string]*chatGPTAttempt
	catalogs map[string]chatGPTCatalog
	listener *http.Server
}

// SDKs may wrap safe error codes with HTTP/request context. Show only the known
// translated reason for subscription requests, never an upstream body.
func chatGPTPublicError(err error) string {
	for code := range errCatalog {
		if (strings.HasPrefix(code, "chatgpt_") || strings.HasPrefix(code, "subscription_sharing_")) && strings.Contains(err.Error(), code) {
			return trMsg(code)
		}
	}
	return trMsg("chatgpt_request_failed")
}

func (s *Server) startChatGPTOAuth(w http.ResponseWriter, r *http.Request) {
	oauthHeaders(w)
	var req struct {
		Origin    string `json:"origin"`
		AccountID string `json:"account_id"`
	}
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		writeErr(w, 400, "oauth_invalid_request")
		return
	}
	origin, ok := oauthOrigin(r, req.Origin)
	u, _ := url.Parse(origin)
	if !ok || u.Hostname() != "localhost" {
		writeErr(w, 400, "chatgpt_localhost_required")
		return
	}
	state, e1 := randomToken(32)
	verifier, e2 := randomToken(48)
	nonce, e3 := randomToken(32)
	if e1 != nil || e2 != nil || e3 != nil {
		writeErr(w, 500, "oauth_unavailable")
		return
	}
	o := &s.chatGPT
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.attempts == nil {
		o.attempts = make(map[string]*chatGPTAttempt)
	}
	for id, a := range o.attempts {
		if time.Now().After(a.expires) {
			delete(o.attempts, id)
		}
	}
	if len(o.attempts) >= 32 {
		writeErr(w, 429, "oauth_try_again")
		return
	}
	if o.listener == nil {
		// Docker publishes this listener only to the host loopback interface.
		bind := "127.0.0.1:8789"
		if os.Getenv("ARTEX_CHATGPT_CALLBACK_LISTEN") == "0.0.0.0:8789" {
			bind = "0.0.0.0:8789"
		}
		listener, err := net.Listen("tcp", bind)
		if err != nil {
			writeErr(w, 503, "chatgpt_callback_unavailable")
			return
		}
		callback := &http.Server{Handler: http.HandlerFunc(s.chatGPTCallback), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, IdleTimeout: 10 * time.Second, MaxHeaderBytes: 16384}
		o.listener = callback
		go func() { _ = callback.Serve(listener) }()
		go func() {
			ticker := time.NewTicker(time.Minute)
			defer ticker.Stop()
			for {
				select {
				case <-s.ctx.Done():
				case <-ticker.C:
				}
				o.mu.Lock()
				for id, a := range o.attempts {
					if time.Now().After(a.expires) {
						delete(o.attempts, id)
					}
				}
				stop := len(o.attempts) == 0 || s.ctx.Err() != nil
				if stop && o.listener == callback {
					o.listener = nil
					_ = callback.Close()
				}
				o.mu.Unlock()
				if stop {
					return
				}
			}
		}()
	}
	host, client, err := o.store.Registration(req.AccountID)
	if err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	o.attempts[state] = &chatGPTAttempt{owner: oauthSession(r), origin: origin, accountID: req.AccountID, clientID: client, verifier: verifier, nonce: nonce, expires: time.Now().Add(5 * time.Minute)}
	challenge := sha256.Sum256([]byte(verifier))
	query := url.Values{"client_id": {client}, "response_type": {"code"}, "redirect_uri": {chatgpt.CallbackURL}, "scope": {chatgpt.Scope}, "resource": {chatgpt.BaseURL}, "state": {state}, "nonce": {nonce}, "code_challenge": {base64.RawURLEncoding.EncodeToString(challenge[:])}, "code_challenge_method": {"S256"}, "ext_agent_host_id": {host}}
	if client == "dynamic_agent_client" {
		query.Set("agent_name_hint", "ARTEX Community")
	}
	writeJSON(w, 200, map[string]string{"url": chatgpt.Issuer + "/api/accounts/authorize?" + query.Encode()})
}

func (s *Server) chatGPTCallback(w http.ResponseWriter, r *http.Request) {
	oauthHeaders(w)
	if r.Method != "GET" || r.URL.Path != "/auth/callback" || r.Host != "127.0.0.1:8789" || len(r.URL.RawQuery) > 16384 {
		http.NotFound(w, r)
		return
	}
	q := r.URL.Query()
	o := &s.chatGPT
	o.mu.Lock()
	a, ok := o.attempts[q.Get("state")]
	if !ok || a.returned || time.Now().After(a.expires) {
		o.mu.Unlock()
		http.Error(w, trMsg("oauth_expired"), 400)
		return
	}
	a.returned = true
	a.failed = q.Get("error") != "" || q.Get("code") == "" || len(q.Get("code")) > 8192
	issued := q.Get("client_id")
	if a.clientID == "dynamic_agent_client" {
		if !strings.HasPrefix(issued, "oaiapp_") || len(issued) > 256 {
			a.failed = true
		} else {
			a.clientID = issued
		}
	} else if issued != "" && issued != a.clientID {
		a.failed = true
	}
	if !a.failed {
		a.code = q.Get("code")
	}
	destination := a.origin + "/system/llm/chatgpt/callback/?attempt=" + url.QueryEscape(q.Get("state"))
	o.mu.Unlock()
	// Return to the original localhost origin. Authorization codes stay server-side.
	http.Redirect(w, r, destination, http.StatusSeeOther)
}

func (s *Server) completeChatGPTOAuth(w http.ResponseWriter, r *http.Request) {
	oauthHeaders(w)
	var req struct {
		Attempt string `json:"attempt"`
	}
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		writeErr(w, 400, "oauth_invalid_request")
		return
	}
	o := &s.chatGPT
	o.mu.Lock()
	a, ok := o.attempts[req.Attempt]
	if ok && a.owner == oauthSession(r) && a.returned {
		delete(o.attempts, req.Attempt)
	} else {
		ok = false
	}
	defer o.mu.Unlock() // Serialize credential replacement against disconnect.
	if !ok || a.failed || time.Now().After(a.expires) {
		writeErr(w, 400, "oauth_expired")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
	defer cancel()
	account, err := o.store.Complete(ctx, a.accountID, a.clientID, a.code, a.verifier, a.nonce)
	if err != nil {
		writeErr(w, 502, err.Error())
		return
	}
	delete(o.catalogs, account.ID)
	writeJSON(w, 200, account)
}

func (s *Server) chatGPTAccounts(w http.ResponseWriter, r *http.Request) {
	oauthHeaders(w)
	accounts, err := s.chatGPT.store.Accounts()
	if err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, map[string]any{"accounts": accounts})
}

func (s *Server) chatGPTAccountAction(w http.ResponseWriter, r *http.Request) {
	oauthHeaders(w)
	var req struct {
		AccountID string `json:"account_id"`
		Name      string `json:"name"`
		Model     string `json:"model"`
	}
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	if json.NewDecoder(r.Body).Decode(&req) != nil || req.AccountID == "" {
		writeErr(w, 400, "oauth_invalid_request")
		return
	}
	o := &s.chatGPT
	ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
	defer cancel()
	switch r.PathValue("action") {
	case "models":
		models, err := o.store.Models(ctx, req.AccountID)
		if err != nil {
			writeErr(w, 502, err.Error())
			return
		}
		o.mu.Lock()
		if o.catalogs == nil {
			o.catalogs = make(map[string]chatGPTCatalog)
		}
		o.catalogs[req.AccountID] = chatGPTCatalog{models, time.Now().Add(10 * time.Minute)}
		o.mu.Unlock()
		writeJSON(w, 200, map[string]any{"models": models})
	case "disconnect":
		o.mu.Lock()
		confirmed, err := o.store.Disconnect(ctx, req.AccountID)
		delete(o.catalogs, req.AccountID)
		for state, a := range o.attempts {
			if a.accountID == req.AccountID {
				delete(o.attempts, state)
			}
		}
		o.mu.Unlock()
		if err != nil {
			writeErr(w, 500, err.Error())
			return
		}
		writeJSON(w, 200, map[string]bool{"revoked": confirmed})
	case "profile":
		o.mu.Lock()
		catalog := o.catalogs[req.AccountID]
		allowed := false
		for _, m := range catalog.models {
			if m.Slug == req.Model && time.Now().Before(catalog.expires) {
				allowed = true
			}
		}
		o.mu.Unlock()
		if !allowed || strings.TrimSpace(req.Name) == "" || len(req.Name) > 120 {
			writeErr(w, 400, "chatgpt_select_model")
			return
		}
		accounts, err := o.store.Accounts()
		connected := false
		for _, a := range accounts {
			if a.ID == req.AccountID && a.Connected {
				connected = true
			}
		}
		if err != nil || !connected {
			writeErr(w, 409, "chatgpt_reconnect_required")
			return
		}
		pg := s.pg(w)
		if pg == nil {
			return
		}
		p := db.LLMProfile{Name: strings.TrimSpace(req.Name), Format: "openai-responses", BaseURL: chatgpt.BaseURL, Model: req.Model, APIKey: chatgpt.Prefix + req.AccountID, Streaming: true, PoolExclude: true}
		id, err := pg.SaveProfile(&p)
		if err != nil {
			writeErr(w, 500, "oauth_save_failed")
			return
		}
		s.invalidateProfileAgents()
		s.reapplyActiveProfile()
		writeJSON(w, 200, map[string]any{"id": id})
	default:
		http.NotFound(w, r)
	}
}
