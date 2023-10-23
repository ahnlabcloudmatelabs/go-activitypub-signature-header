package signature_header

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"strings"
)

// Example:
//
//	privateKey, err := signature_header.PrivateKeyFromBytes(privateKeyBytes)
//	date := signature_header.Date()
//	digest := signature_header.Digest(crypto.SHA256, message)
//	signature, err := signature_header.Signature{
//	  PrivateKey: privateKey,
//	  Algorithm:  crypto.SHA256,
//	  Date:       date,
//	  Digest:     digest,
//	  Host:       "yodangang.express",
//	  Path:       "/users/9iffvxhojp/inbox",
//	  KeyID:      "https://api.dev.snippet.cloudmt.co.kr/@juunini#main-key",
//	}.String()
type Signature struct {
	PrivateKey *rsa.PrivateKey
	Algorithm  crypto.Hash
	Date       string
	Digest     string
	Host       string
	Path       string
	KeyID      string
}

func (s Signature) String() (string, error) {
	message := fmt.Sprintf(`(request-target): post %s
date: %s
host: %s
digest: %s`,
		s.Path,
		s.Date,
		s.Host,
		s.Digest,
	)

	signed, err := rsa.SignPKCS1v15(
		rand.Reader,
		s.PrivateKey,
		s.Algorithm,
		hashing(s.Algorithm, []byte(message)),
	)
	if err != nil {
		return "", err
	}
	signature := base64.StdEncoding.EncodeToString(signed)

	return strings.Join([]string{
		fmt.Sprintf(`keyId="%s"`, s.KeyID),
		fmt.Sprintf(`algorithm="rsa-%s"`, strings.Replace(strings.ToLower(s.Algorithm.String()), "-", "", 1)),
		`headers="(request-target) date host digest"`,
		fmt.Sprintf(`signature="%s"`, signature),
	}, ","), nil
}
