package chatgpt

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type identityClaims struct {
	jwt.RegisteredClaims
	Nonce           string `json:"nonce"`
	Email           string `json:"email"`
	AuthorizedParty string `json:"azp"`
}

func (s *Store) publicJSON(ctx context.Context, endpoint string, out any) error {
	u, err := url.Parse(endpoint)
	if err != nil || u.Scheme != "https" || u.Host != "auth.openai.com" || u.User != nil || u.Fragment != "" {
		return ErrExchange
	}
	req, _ := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	resp, err := s.client.Do(req)
	if err != nil {
		return ErrExchange
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 || json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(out) != nil {
		return ErrExchange
	}
	return nil
}

func (s *Store) identity(ctx context.Context, raw, client, nonce string) (identityClaims, error) {
	var claims identityClaims
	var jwks struct {
		Keys []struct{ Kty, Kid, Use, Alg, N, E string } `json:"keys"`
	}
	if raw == "" || s.publicJSON(ctx, Issuer+"/.well-known/jwks.json", &jwks) != nil {
		return claims, ErrExchange
	}
	parsed, err := jwt.ParseWithClaims(raw, &claims, func(t *jwt.Token) (any, error) {
		kid, _ := t.Header["kid"].(string)
		for _, k := range jwks.Keys {
			if kid == "" || k.Kid != kid || k.Kty != "RSA" || (k.Use != "" && k.Use != "sig") || (k.Alg != "" && k.Alg != "RS256") {
				continue
			}
			n, en := base64.RawURLEncoding.DecodeString(k.N)
			e, ee := base64.RawURLEncoding.DecodeString(k.E)
			if en != nil || ee != nil || len(n) < 256 || len(e) == 0 || len(e) > 4 {
				return nil, ErrExchange
			}
			exponent := int64(0)
			for _, b := range e {
				exponent = exponent*256 + int64(b)
			}
			if exponent < 3 || exponent > 2147483647 || exponent%2 == 0 {
				return nil, ErrExchange
			}
			return &rsa.PublicKey{N: new(big.Int).SetBytes(n), E: int(exponent)}, nil
		}
		return nil, ErrExchange
	}, jwt.WithValidMethods([]string{"RS256"}), jwt.WithIssuer(Issuer), jwt.WithAudience(client), jwt.WithExpirationRequired(), jwt.WithIssuedAt(), jwt.WithLeeway(5*time.Second))
	if err != nil || !parsed.Valid || claims.Subject == "" || claims.IssuedAt == nil || (nonce != "" && claims.Nonce != nonce) || (claims.AuthorizedParty != "" && claims.AuthorizedParty != client) || (len(claims.Audience) > 1 && claims.AuthorizedParty != client) {
		return identityClaims{}, ErrExchange
	}
	return claims, nil
}

func (s *Store) revocationEndpoint(ctx context.Context) (string, error) {
	var discovery struct {
		Issuer     string `json:"issuer"`
		Revocation string `json:"revocation_endpoint"`
	}
	if s.publicJSON(ctx, Issuer+"/.well-known/openid-configuration", &discovery) != nil || discovery.Issuer != Issuer {
		return "", ErrExchange
	}
	u, err := url.Parse(discovery.Revocation)
	if err != nil || u.Scheme != "https" || u.Host != "auth.openai.com" || u.User != nil || u.Fragment != "" {
		return "", ErrExchange
	}
	return u.String(), nil
}
