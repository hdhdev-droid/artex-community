package chatgpt

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestIdentityValidation(t *testing.T) {
	s := testStore(t)
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	public, err := json.Marshal(map[string]any{"keys": []map[string]string{{
		"kty": "RSA", "kid": "test-signature", "alg": "RS256", "use": "sig",
		"n": base64.RawURLEncoding.EncodeToString(key.N.Bytes()),
		"e": base64.RawURLEncoding.EncodeToString(big.NewInt(int64(key.E)).Bytes()),
	}}})
	if err != nil {
		t.Fatal(err)
	}
	s.client.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() != Issuer+"/.well-known/jwks.json" {
			t.Fatal("unexpected endpoint")
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(public))}, nil
	})
	for _, tc := range []struct {
		name   string
		change func(*identityClaims)
		valid  bool
	}{
		{"valid", func(*identityClaims) {}, true},
		{"nonce", func(c *identityClaims) { c.Nonce = "other" }, false},
		{"audience", func(c *identityClaims) { c.Audience = jwt.ClaimStrings{"other-client"} }, false},
		{"issuer", func(c *identityClaims) { c.Issuer = "https://example.invalid" }, false},
		{"subject", func(c *identityClaims) { c.Subject = "" }, false},
		{"expired", func(c *identityClaims) { c.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-time.Minute)) }, false},
		{"issued at", func(c *identityClaims) { c.IssuedAt = nil }, false},
		{"authorized party", func(c *identityClaims) { c.AuthorizedParty = "other-client" }, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			claims := identityClaims{RegisteredClaims: jwt.RegisteredClaims{Issuer: Issuer, Subject: "test-subject", Audience: jwt.ClaimStrings{"oaiapp_test"}, ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)), IssuedAt: jwt.NewNumericDate(time.Now())}, Nonce: "test-nonce"}
			tc.change(&claims)
			token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
			token.Header["kid"] = "test-signature"
			raw, err := token.SignedString(key)
			if err != nil {
				t.Fatal(err)
			}
			_, err = s.identity(context.Background(), raw, "oaiapp_test", "test-nonce")
			if (err == nil) != tc.valid {
				t.Fatalf("identity validation mismatch: %v", err)
			}
		})
	}
}
