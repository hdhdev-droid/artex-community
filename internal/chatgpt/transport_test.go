package chatgpt

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func testStore(t *testing.T) *Store {
	t.Helper()
	previous := Default()
	s := New(t.TempDir(), []byte("test-only-installation-secret"))
	t.Cleanup(func() { current.Store(previous) })
	s.client.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		t.Fatal("unexpected network request")
		return nil, errors.New("network disabled")
	})
	return s
}

func testCredential(t *testing.T, s *Store, expired bool) {
	t.Helper()
	expires := time.Now().Add(time.Hour)
	if expired {
		expires = time.Now().Add(-time.Hour)
	}
	if err := s.write(vault{HostID: "test-host", Accounts: []credential{{
		ID: "account", ClientID: "oaiapp_test", Subject: "test-subject", Label: "Test account",
		AccessToken: "test-access", RefreshToken: "test-refresh", Scope: Scope, Expires: expires,
	}}}); err != nil {
		t.Fatal(err)
	}
}

func TestCompletedStream(t *testing.T) {
	for _, tc := range []struct{ name, events, want string }{
		{"completed", "data: {\"type\":\"response.completed\"}\n\n", ""},
		{"early eof", "data: {\"type\":\"response.output_text.delta\",\"delta\":\"partial\"}\n\n", "chatgpt_response_incomplete"},
		{"incomplete", "data: {\"type\":\"response.incomplete\"}\n\n", "chatgpt_response_incomplete"},
		{"quota", "data: {\"type\":\"response.failed\",\"response\":{\"error\":{\"code\":\"subscription_sharing_usage_limit_exceeded\"}}}\n\n", "subscription_sharing_usage_limit_exceeded"},
		{"redacted error", "data: {\"type\":\"error\",\"error\":{\"code\":\"private-upstream-value\"}}\n\n", "chatgpt_response_failed"},
		{"malformed", "data: not-json\n\n", "chatgpt_response_incomplete"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := io.NopCloser(strings.NewReader(tc.events))
			stream := &completedStream{body: body, scanner: bufio.NewScanner(body)}
			_, err := io.ReadAll(stream)
			if tc.want == "" {
				if err != nil {
					t.Fatal(err)
				}
				return
			}
			if err == nil || err.Error() != tc.want {
				t.Fatalf("expected %s, got %v", tc.want, err)
			}
		})
	}
}

func TestTransportBoundaries(t *testing.T) {
	s := testStore(t)
	testCredential(t, s, false)
	calls := 0
	base := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.Header.Get("Authorization") != "Bearer test-access" {
			t.Fatal("credential was not replaced")
		}
		return &http.Response{StatusCode: 200, Header: http.Header{}, Body: io.NopCloser(strings.NewReader("data: {\"type\":\"response.completed\"}\n\n"))}, nil
	})
	client := &http.Client{Transport: &transport{base: base, store: s, id: "account"}}
	for _, tc := range []struct {
		endpoint, body string
		valid          bool
	}{
		{BaseURL + "/responses", `{"stream":true,"store":false}`, true},
		{"https://example.invalid/responses", `{"stream":true,"store":false}`, false},
		{BaseURL + "/responses", `{"stream":false,"store":false}`, false},
		{BaseURL + "/responses", `{"stream":true,"store":true}`, false},
	} {
		req, _ := http.NewRequestWithContext(context.Background(), "POST", tc.endpoint, strings.NewReader(tc.body))
		resp, err := client.Do(req)
		if tc.valid {
			if err != nil {
				t.Fatal(err)
			}
			_, err = io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
		} else if err == nil {
			resp.Body.Close()
			t.Fatal("unsafe request accepted")
		}
	}
	if calls != 1 {
		t.Fatalf("unexpected outgoing request count: %d", calls)
	}
}

func TestRequiresAction(t *testing.T) {
	for _, code := range []string{"chatgpt_reconnect_required", "chatgpt_scope_missing", "subscription_sharing_usage_limit_exceeded", "subscription_sharing_usage_unavailable"} {
		if !RequiresAction(errors.New("wrapped: " + code)) {
			t.Fatalf("must stop retries for %s", code)
		}
	}
	if RequiresAction(errors.New("connection reset")) || RequiresAction(nil) {
		t.Fatal("ordinary errors must retain existing retry behavior")
	}
}
