package jwt

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gopkg.in/square/go-jose.v2"
)

// The ID Token represents a JWT passed to the client as part of the token response.
//
// https://openid.net/specs/openid-connect-core-1_0.html#IDToken
type IDToken struct {
	Issuer     string `json:"iss"`
	UserID     string `json:"sub"`
	ClientID   string `json:"aud"`
	Expiration int64  `json:"exp"`
	IssuedAt   int64  `json:"iat"`

	Nonce string `json:"nonce,omitempty"`

	Email         string `json:"email,omitempty"`
	EmailVerified *bool  `json:"email_verified,omitempty"`
}

func NewIDToken(userID string, ttl time.Duration) *IDToken {
	now := time.Now()
	idToken := &IDToken{
		Expiration: now.Add(ttl).Unix(),
		IssuedAt:   now.Unix(),
		UserID:     userID,
	}

	return idToken
}

func (t *IDToken) UpdateEmail(domain string) {
	if strings.Contains(t.UserID, "@") {
		v := true
		t.EmailVerified = &v
		t.Email = t.UserID
	} else {
		v := false
		t.EmailVerified = &v
		t.Email = fmt.Sprintf("%s@%s", t.UserID, domain)
	}
}

// Encode serializes and signs the instance.
func (t *IDToken) Encode(signer jose.Signer) (output string, fail error) {
	payload, err := json.Marshal(t)
	if err != nil {
		fail = fmt.Errorf("failed to marshal token: %w", err)
		return
	}

	jws, err := signer.Sign(payload)
	if err != nil {
		fail = fmt.Errorf("failed to sign token: %w", err)
		return
	}

	raw, err := jws.CompactSerialize()
	if err != nil {
		fail = fmt.Errorf("failed to serialize token: %w", err)
		return
	}

	output = raw
	return
}
