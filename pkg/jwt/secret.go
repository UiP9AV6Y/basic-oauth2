package jwt

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"gopkg.in/square/go-jose.v2"
)

var ErrDecoding = errors.New("unable to decode private key")

type Secret struct {
	use       string
	key       *rsa.PrivateKey
	algorithm jose.SignatureAlgorithm
}

func NewSecret(key *rsa.PrivateKey) *Secret {
	secret := &Secret{
		algorithm: jose.RS256,
		use:       "sig",
		key:       key,
	}

	return secret
}

func ParseSecretFile(path string) (*Secret, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return ParseSecretBytes(bytes)
}

func ParseSecretBytes(b []byte) (*Secret, error) {
	block, _ := pem.Decode(b)
	if block == nil {
		return nil, ErrDecoding
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse key: %w", err)
	}

	return NewSecret(key), nil
}

func GenerateSecret(random io.Reader, bits int) (*Secret, error) {
	key, err := rsa.GenerateKey(random, bits)
	if err != nil {
		return nil, err
	}

	return NewSecret(key), nil
}

func (s *Secret) NewSigner() (jose.Signer, error) {
	key := jose.SigningKey{
		Algorithm: s.algorithm,
		Key:       s.key,
	}

	return jose.NewSigner(key, nil)
}

func (s *Secret) CalculateThumbprint() (string, error) {
	enc, err := x509.MarshalPKIXPublicKey(&s.key.PublicKey)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(enc)
	return hex.EncodeToString(hash[:16]), nil
}

func (s *Secret) NewPrivateKey() (*jose.JSONWebKey, error) {
	thumb, err := s.CalculateThumbprint()
	if err != nil {
		return nil, err
	}

	key := &jose.JSONWebKey{
		Key:       s.key,
		Use:       s.use,
		Algorithm: string(s.algorithm),
		KeyID:     thumb,
	}
	return key, nil
}

func (s *Secret) NewPublicKey() (*jose.JSONWebKey, error) {
	thumb, err := s.CalculateThumbprint()
	if err != nil {
		return nil, err
	}

	key := &jose.JSONWebKey{
		Key:       &s.key.PublicKey,
		Use:       s.use,
		Algorithm: string(s.algorithm),
		KeyID:     thumb,
	}
	return key, nil
}

func (s *Secret) NewPublicKeySet() (*jose.JSONWebKeySet, error) {
	key, err := s.NewPublicKey()
	if err != nil {
		return nil, err
	}

	keys := &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			*key,
		},
	}
	return keys, nil
}
