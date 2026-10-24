package chatgpt

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestEncryptedStoreRestart(t *testing.T) {
	s := testStore(t)
	if _, err := os.Stat(s.path); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("constructor must not write files")
	}
	host, client, err := s.Registration("")
	if err != nil || host == "" || client != "dynamic_agent_client" {
		t.Fatal("initial registration failed")
	}
	hostAgain, _, err := s.Registration("")
	if err != nil || hostAgain != host {
		t.Fatal("host registration must survive reload")
	}
	testCredential(t, s, false)
	raw, err := os.ReadFile(s.path)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(raw, []byte("test-access")) || bytes.Contains(raw, []byte("test-refresh")) {
		t.Fatal("credentials stored in plaintext")
	}
	accounts, err := s.Accounts()
	if err != nil || len(accounts) != 1 || !accounts[0].Connected {
		t.Fatal("account metadata not preserved")
	}
	raw[len(raw)-1] ^= 1
	if err := os.WriteFile(s.path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Accounts(); !errors.Is(err, ErrStorage) {
		t.Fatal("tampered vault accepted")
	}
}

func TestRefreshRotationSerialized(t *testing.T) {
	s := testStore(t)
	testCredential(t, s, true)
	calls := 0
	s.client.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.URL.String() != TokenURL || r.ParseForm() != nil || r.Form.Get("grant_type") != "refresh_token" || r.Form.Get("refresh_token") != "test-refresh" || r.Form.Has("scope") {
			t.Error("invalid refresh request")
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewBufferString(`{"access_token":"rotated-access","refresh_token":"rotated-refresh","token_type":"Bearer","expires_in":3600}`))}, nil
	})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			token, err := s.AccessToken(context.Background(), "account")
			if err != nil || token != "rotated-access" {
				t.Error("refresh failed")
			}
		}()
	}
	wg.Wait()
	if calls != 1 {
		t.Fatalf("refresh rotation raced: %d calls", calls)
	}
	v, err := s.read()
	if err != nil || v.Accounts[0].RefreshToken != "rotated-refresh" || v.Accounts[0].Scope != Scope {
		t.Fatal("rotation/scope not persisted")
	}
}

func TestInvalidGrantClearsTokens(t *testing.T) {
	s := testStore(t)
	testCredential(t, s, true)
	s.client.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 400, Body: io.NopCloser(bytes.NewBufferString(`{"error":"invalid_grant"}`))}, nil
	})
	if _, err := s.AccessToken(context.Background(), "account"); !errors.Is(err, ErrReconnect) {
		t.Fatal("must require sign-in")
	}
	v, err := s.read()
	if err != nil || v.Accounts[0].RefreshToken != "" || v.Accounts[0].AccessToken != "" || v.Accounts[0].ClientID != "oaiapp_test" {
		t.Fatal("invalid tokens retained or registration removed")
	}
}

func TestCredentialsSurviveNewStore(t *testing.T) {
	s := testStore(t)
	testCredential(t, s, false)
	dir := filepath.Dir(filepath.Dir(s.path))
	reopened := New(dir, []byte("test-only-installation-secret"))
	reopened.client.Transport = s.client.Transport
	token, err := reopened.AccessToken(context.Background(), "account")
	if err != nil || token != "test-access" {
		t.Fatal("credentials lost across restart")
	}
	wrongKey := New(dir, []byte("different-installation-secret"))
	if _, err := wrongKey.Accounts(); !errors.Is(err, ErrStorage) {
		t.Fatal("wrong installation key accepted")
	}
	if _, err := reopened.Accounts(); err != nil {
		t.Fatal("failed read damaged existing vault")
	}
}

func TestDisconnectOfflineKeepsRegistrationOnly(t *testing.T) {
	s := testStore(t)
	testCredential(t, s, false)
	s.client.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("offline fixture")
	})
	confirmed, err := s.Disconnect(context.Background(), "account")
	if err != nil || confirmed {
		t.Fatal("offline revocation falsely confirmed")
	}
	v, err := s.read()
	if err != nil || v.Accounts[0].ClientID != "oaiapp_test" || v.Accounts[0].RefreshToken != "" || v.Accounts[0].AccessToken != "" {
		t.Fatal("local disconnect did not preserve registration and remove credentials")
	}
	if _, err := s.AccessToken(context.Background(), "account"); !errors.Is(err, ErrReconnect) {
		t.Fatal("disconnected account remained usable")
	}
}
