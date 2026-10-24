// Package chatgpt implements the documented open-source Sign in with ChatGPT flow.
// Credentials never enter the profile database: profiles contain an opaque reference.
package chatgpt

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

const (
	BaseURL     = "https://api.openai.com/v1"
	Issuer      = "https://auth.openai.com"
	TokenURL    = Issuer + "/api/accounts/oauth/token"
	CallbackURL = "http://127.0.0.1:8789/auth/callback"
	Prefix      = "artex-chatgpt:"
	Scope       = "openid profile email offline_access resource.invoke chatgpt.tokens.use.direct"
)

var (
	ErrStorage   = errors.New("chatgpt_storage_failed")
	ErrReconnect = errors.New("chatgpt_reconnect_required")
	ErrExchange  = errors.New("chatgpt_exchange_failed")
	ErrScope     = errors.New("chatgpt_scope_missing")
	ErrModels    = errors.New("chatgpt_models_failed")
	current      atomic.Pointer[Store]
)

type Account struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	Connected bool   `json:"connected"`
}

type credential struct {
	ID, ClientID, Subject, Label              string
	AccessToken, RefreshToken, IDToken, Scope string
	Expires                                   time.Time
}

type vault struct {
	HostID   string
	Accounts []credential
}

type Model struct {
	Slug        string `json:"slug"`
	DisplayName string `json:"display_name"`
	Visibility  string `json:"visibility"`
}

type Store struct {
	mu     sync.Mutex // Also serializes refresh rotation against sign-out and reconnect.
	path   string
	aead   cipher.AEAD
	client *http.Client
}

// New performs no disk or network operations. Startup must not register an account.
func New(dir string, secret []byte) *Store {
	key := sha256.Sum256(append([]byte("artex-chatgpt/v1\x00"), secret...))
	block, _ := aes.NewCipher(key[:])
	aead, _ := cipher.NewGCM(block)
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = nil
	s := &Store{path: filepath.Join(dir, "chatgpt-oauth", "credentials.enc"), aead: aead,
		client: &http.Client{Transport: transport, Timeout: 20 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}}
	current.Store(s)
	return s
}

func Default() *Store { return current.Load() }
func AccountID(ref string) string {
	if strings.HasPrefix(ref, Prefix) {
		return strings.TrimPrefix(ref, Prefix)
	}
	return ""
}

func newID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", ErrStorage
	}
	return id.String(), nil
}

func (s *Store) read() (vault, error) {
	v := vault{}
	b, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return v, nil
	}
	if err != nil || len(b) < s.aead.NonceSize() {
		return v, ErrStorage
	}
	n := s.aead.NonceSize()
	raw, err := s.aead.Open(nil, b[:n], b[n:], []byte("artex-chatgpt/v1"))
	if err != nil || json.Unmarshal(raw, &v) != nil {
		return vault{}, ErrStorage
	}
	return v, nil
}

func (s *Store) write(v vault) error {
	raw, err := json.Marshal(v)
	if err != nil {
		return ErrStorage
	}
	nonce := make([]byte, s.aead.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return ErrStorage
	}
	sealed := s.aead.Seal(nonce, nonce, raw, []byte("artex-chatgpt/v1"))
	if os.MkdirAll(filepath.Dir(s.path), 0700) != nil {
		return ErrStorage
	}
	f, err := os.CreateTemp(filepath.Dir(s.path), ".credentials-*")
	if err != nil {
		return ErrStorage
	}
	name := f.Name()
	defer os.Remove(name)
	if err = f.Chmod(0600); err == nil {
		_, err = f.Write(sealed)
	}
	if err == nil {
		err = f.Sync()
	}
	closeErr := f.Close()
	if err != nil || closeErr != nil || os.Rename(name, s.path) != nil {
		return ErrStorage
	}
	return nil
}

func direct(scope string) bool {
	fields := " " + strings.Join(strings.Fields(scope), " ") + " "
	return strings.Contains(fields, " chatgpt.tokens.use.direct ") && strings.Contains(fields, " resource.invoke ")
}
func find(v *vault, id string) *credential {
	for i := range v.Accounts {
		if v.Accounts[i].ID == id {
			return &v.Accounts[i]
		}
	}
	return nil
}
func metadata(c credential) Account {
	return Account{c.ID, c.Label, c.Subject != "" && c.RefreshToken != "" && direct(c.Scope)}
}

func (s *Store) Accounts() ([]Account, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, err := s.read()
	result := make([]Account, 0, len(v.Accounts))
	for _, c := range v.Accounts {
		result = append(result, metadata(c))
	}
	return result, err
}

func (s *Store) Registration(id string) (host, client string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, err := s.read()
	if err != nil {
		return "", "", err
	}
	client = "dynamic_agent_client"
	if id != "" {
		c := find(&v, id)
		if c == nil {
			return "", "", ErrReconnect
		}
		client = c.ClientID
	}
	if v.HostID == "" {
		uid, e := newID()
		if e != nil {
			return "", "", e
		}
		v.HostID = "urn:uuid:" + uid
		if e = s.write(v); e != nil {
			return "", "", e
		}
	}
	return v.HostID, client, nil
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int64  `json:"expires_in"`
	Error        string `json:"error"`
}

func (s *Store) token(ctx context.Context, form url.Values) (tokenResponse, error) {
	var out tokenResponse
	req, _ := http.NewRequestWithContext(ctx, "POST", TokenURL, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := s.client.Do(req)
	if err != nil {
		return out, ErrExchange
	}
	defer resp.Body.Close()
	if json.NewDecoder(io.LimitReader(resp.Body, 128<<10)).Decode(&out) != nil {
		return tokenResponse{}, ErrExchange
	}
	if out.Error == "invalid_grant" {
		return tokenResponse{}, ErrReconnect
	}
	if resp.StatusCode != 200 || out.AccessToken == "" || !strings.EqualFold(out.TokenType, "Bearer") || out.ExpiresIn <= 0 || out.ExpiresIn > 31*86400 {
		return tokenResponse{}, ErrExchange
	}
	return out, nil
}

func (s *Store) Complete(ctx context.Context, id, client, code, verifier, nonce string) (Account, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, err := s.read()
	if err != nil {
		return Account{}, err
	}
	if !strings.HasPrefix(client, "oaiapp_") || len(client) > 256 {
		return Account{}, ErrExchange
	}
	old := find(&v, id)
	if id != "" && (old == nil || old.ClientID != client) {
		return Account{}, ErrReconnect
	}
	// Keep an issued registration even if exchange fails with invalid_grant.
	// It has no usable credentials or identity until signature validation succeeds.
	if old == nil {
		for i := range v.Accounts {
			if v.Accounts[i].ClientID == client {
				old = &v.Accounts[i]
				break
			}
		}
	}
	if old == nil {
		uid, e := newID()
		if e != nil {
			return Account{}, e
		}
		v.Accounts = append(v.Accounts, credential{ID: uid, ClientID: client, Label: "ChatGPT · " + uid[:8]})
		old = &v.Accounts[len(v.Accounts)-1]
		if e = s.write(v); e != nil {
			return Account{}, e
		}
	}
	t, err := s.token(ctx, url.Values{"grant_type": {"authorization_code"}, "client_id": {client}, "code": {code}, "code_verifier": {verifier}, "redirect_uri": {CallbackURL}, "resource": {BaseURL}})
	if err != nil {
		return Account{}, err
	}
	claims, err := s.identity(ctx, t.IDToken, client, nonce)
	if err != nil {
		return Account{}, err
	}
	permitted := direct(t.Scope) && t.RefreshToken != ""
	if old.Subject != "" && old.Subject != claims.Subject {
		return Account{}, ErrReconnect
	}
	old.Subject = claims.Subject
	old.Label = claims.Email
	if old.Label == "" {
		old.Label = "ChatGPT · " + old.ID[:8]
	}
	old.AccessToken, old.RefreshToken, old.IDToken, old.Scope = t.AccessToken, t.RefreshToken, t.IDToken, t.Scope
	old.Expires = time.Now().Add(time.Duration(t.ExpiresIn) * time.Second)
	if !permitted {
		old.AccessToken, old.RefreshToken = "", ""
	}
	if err = s.write(v); err != nil {
		return Account{}, err
	}
	if !permitted {
		return metadata(*old), ErrScope
	}
	return metadata(*old), nil
}

func (s *Store) AccessToken(ctx context.Context, id string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, err := s.read()
	if err != nil {
		return "", err
	}
	c := find(&v, id)
	if c == nil || c.RefreshToken == "" || !direct(c.Scope) {
		return "", ErrReconnect
	}
	if c.AccessToken != "" && time.Until(c.Expires) > time.Minute {
		return c.AccessToken, nil
	}
	t, err := s.token(ctx, url.Values{"grant_type": {"refresh_token"}, "client_id": {c.ClientID}, "refresh_token": {c.RefreshToken}, "resource": {BaseURL}})
	if err != nil {
		if errors.Is(err, ErrReconnect) {
			c.AccessToken, c.RefreshToken, c.IDToken = "", "", ""
			if s.write(v) != nil {
				return "", ErrStorage
			}
		}
		return "", err
	}
	// An omitted scope on refresh preserves the previously granted scope (OAuth 2.0).
	if t.Scope != "" {
		c.Scope = t.Scope
	}
	if t.RefreshToken != "" {
		c.RefreshToken = t.RefreshToken
	}
	if t.IDToken != "" {
		claims, e := s.identity(ctx, t.IDToken, c.ClientID, "")
		if e != nil || claims.Subject != c.Subject {
			// A rotated refresh token cannot safely be reused after identity
			// validation fails. Require a fresh sign-in instead of looping.
			c.AccessToken, c.RefreshToken, c.IDToken = "", "", ""
			if s.write(v) != nil {
				return "", ErrStorage
			}
			return "", ErrReconnect
		}
		c.IDToken = t.IDToken
	}
	c.AccessToken = t.AccessToken
	c.Expires = time.Now().Add(time.Duration(t.ExpiresIn) * time.Second)
	if err = s.write(v); err != nil {
		return "", err
	}
	if !direct(c.Scope) {
		return "", ErrScope
	}
	return c.AccessToken, nil
}

func (s *Store) Models(ctx context.Context, id string) ([]Model, error) {
	token, err := s.AccessToken(ctx, id)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequestWithContext(ctx, "GET", BaseURL+"/models", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, ErrModels
	}
	defer resp.Body.Close()
	if resp.StatusCode == 401 {
		return nil, ErrReconnect
	}
	var body struct {
		Models []Model `json:"models"`
	}
	if resp.StatusCode != 200 || json.NewDecoder(io.LimitReader(resp.Body, 2<<20)).Decode(&body) != nil {
		return nil, ErrModels
	}
	models := []Model{}
	for _, m := range body.Models {
		if m.Visibility == "list" && m.Slug != "" {
			models = append(models, m)
		}
	}
	return models, nil
}

// Disconnect blocks refresh while revocation is attempted, then erases local
// tokens even when the remote server cannot confirm revocation. Registration stays.
func (s *Store) Disconnect(ctx context.Context, id string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, err := s.read()
	if err != nil {
		return false, err
	}
	c := find(&v, id)
	if c == nil {
		return false, ErrReconnect
	}
	confirmed := c.RefreshToken == ""
	if !confirmed {
		endpoint, e := s.revocationEndpoint(ctx)
		if e == nil {
			for attempt := 0; attempt < 2; attempt++ {
				form := url.Values{"token": {c.RefreshToken}, "token_type_hint": {"refresh_token"}, "client_id": {c.ClientID}}
				req, _ := http.NewRequestWithContext(ctx, "POST", endpoint, strings.NewReader(form.Encode()))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				resp, requestErr := s.client.Do(req)
				if requestErr == nil {
					resp.Body.Close()
					confirmed = resp.StatusCode == 200
					if resp.StatusCode < 500 {
						break
					}
				}
				if ctx.Err() != nil {
					break
				}
			}
		}
	}
	c.AccessToken, c.RefreshToken, c.IDToken, c.Scope = "", "", "", ""
	c.Expires = time.Time{}
	return confirmed, s.write(v)
}
