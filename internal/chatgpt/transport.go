package chatgpt

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type transport struct {
	base  http.RoundTripper
	store *Store
	id    string
}

// RequiresAction prevents an authentication/plan error from silently replaying
// a request or falling through to a separately billed provider.
func RequiresAction(err error) bool {
	if err == nil {
		return false
	}
	for _, code := range []string{"chatgpt_reconnect_required", "chatgpt_scope_missing", "chatgpt_storage_failed", "chatgpt_exchange_failed", "subscription_sharing_usage_limit_exceeded", "subscription_sharing_usage_unavailable"} {
		if strings.Contains(err.Error(), code) {
			return true
		}
	}
	return false
}

// Client isolates subscription requests from custom endpoints, proxies and redirects.
func Client(base *http.Client, ref string) *http.Client {
	c := *base
	c.Transport = &transport{base.Transport, Default(), AccountID(ref)}
	c.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	return &c
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.store == nil || t.id == "" {
		return nil, ErrReconnect
	}
	if req.Method != "POST" || req.URL.String() != BaseURL+"/responses" {
		return nil, errors.New("chatgpt_endpoint_invalid")
	}
	// Fail closed if a caller bypasses the streaming provider wrapper.
	body, err := io.ReadAll(io.LimitReader(req.Body, 16<<20+1))
	req.Body.Close()
	var flags struct {
		Store  bool `json:"store"`
		Stream bool `json:"stream"`
	}
	if err != nil || len(body) > 16<<20 || json.Unmarshal(body, &flags) != nil || flags.Store || !flags.Stream {
		return nil, errors.New("chatgpt_request_invalid")
	}
	token, err := t.store.AccessToken(req.Context(), t.id)
	if err != nil {
		return nil, err
	}
	r := req.Clone(req.Context())
	r.Body = io.NopCloser(bytes.NewReader(body))
	r.Header = req.Header.Clone()
	r.Header.Set("Authorization", "Bearer "+token)
	resp, err := t.base.RoundTrip(r)
	if err != nil {
		return nil, errors.New("chatgpt_request_failed")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		var failure struct {
			Error struct {
				Code string `json:"code"`
			} `json:"error"`
		}
		_ = json.NewDecoder(io.LimitReader(resp.Body, 65536)).Decode(&failure)
		if resp.StatusCode == 401 || resp.StatusCode == 403 {
			return nil, ErrReconnect
		}
		return nil, fmt.Errorf("%w (status %d)", streamError(failure.Error.Code), resp.StatusCode)
	}
	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 4096), 4<<20)
	resp.Body = &completedStream{body: resp.Body, scanner: scanner}
	return resp, nil
}

func streamError(code string) error {
	switch code {
	case "subscription_sharing_usage_limit_exceeded", "subscription_sharing_usage_unavailable", "context_length_exceeded":
		return errors.New(code)
	default:
		return errors.New("chatgpt_response_failed")
	}
}

// The SDK accepts EOF and response.incomplete as a finished response. This reader
// forwards complete SSE frames but requires response.completed before clean EOF.
type completedStream struct {
	body      io.ReadCloser
	scanner   *bufio.Scanner
	pending   []byte
	completed bool
	terminal  error
}

func (s *completedStream) Close() error { return s.body.Close() }
func (s *completedStream) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	for len(s.pending) == 0 {
		if s.terminal != nil {
			return 0, s.terminal
		}
		if s.completed {
			return 0, io.EOF
		}
		var frame, data strings.Builder
		for s.scanner.Scan() {
			line := s.scanner.Text()
			frame.WriteString(line)
			frame.WriteByte('\n')
			if strings.HasPrefix(line, "data:") {
				if data.Len() > 0 {
					data.WriteByte('\n')
				}
				data.WriteString(strings.TrimPrefix(strings.TrimPrefix(line, "data:"), " "))
			}
			if frame.Len() > 4<<20 {
				s.terminal = errors.New("chatgpt_response_failed")
				return 0, s.terminal
			}
			if line == "" {
				break
			}
		}
		if frame.Len() == 0 || s.scanner.Err() != nil {
			s.terminal = errors.New("chatgpt_response_incomplete")
			return 0, s.terminal
		}
		if data.Len() > 0 {
			var event struct {
				Type  string `json:"type"`
				Error struct {
					Code string `json:"code"`
				} `json:"error"`
				Response struct {
					Error struct {
						Code string `json:"code"`
					} `json:"error"`
				} `json:"response"`
			}
			if json.Unmarshal([]byte(data.String()), &event) != nil {
				s.terminal = errors.New("chatgpt_response_incomplete")
				return 0, s.terminal
			}
			switch event.Type {
			case "response.completed":
				s.completed = true
			case "response.incomplete":
				s.terminal = errors.New("chatgpt_response_incomplete")
			case "response.failed", "error":
				code := event.Response.Error.Code
				if code == "" {
					code = event.Error.Code
				}
				s.terminal = streamError(code)
			}
			if s.terminal != nil {
				return 0, s.terminal
			}
		}
		frame.WriteByte('\n')
		s.pending = []byte(frame.String())
	}
	n := copy(p, s.pending)
	s.pending = s.pending[n:]
	return n, nil
}
