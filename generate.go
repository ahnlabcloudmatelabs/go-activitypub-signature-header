package signature_header

import "crypto"

/*
Example:

	signature_header.GenerateInput{
		PrivateKeyBytes: []byte("-----BEGIN RSA PRIVATE KEY-----..."),
		Algorithm: crypto.SHA256,
		Host: "example.com",
		Path: "/inbox",
		KeyID: "https://snippet.social/@juunini#main-key",
	}
*/
type GenerateInput struct {
	PrivateKeyBytes []byte
	Algorithm       crypto.Hash
	Host            string
	Path            string
	Body            []byte
	KeyID           string
}

type GenerateOutput struct {
	Date      string
	Host      string
	Digest    string
	Signature string
}

/*
Example:

	headers, err := signature_header.Generate(signature_header.GenerateInput{
		PrivateKeyBytes: []byte("-----BEGIN RSA PRIVATE KEY-----..."),
		Algorithm: crypto.SHA256,
		Host: "example.com",
		Path: "/inbox",
		KeyID: "https://snippet.social/@juunini#main-key",
	})
*/
func Generate(input GenerateInput) (*GenerateOutput, error) {
	date := Date()
	digest := Digest(input.Algorithm, input.Body)

	privateKey, err := PrivateKeyFromBytes(input.PrivateKeyBytes)
	if err != nil {
		return nil, err
	}

	signature, err := Signature{
		PrivateKey: privateKey,
		Algorithm:  input.Algorithm,
		Date:       date,
		Digest:     digest,
		Host:       input.Host,
		Path:       input.Path,
		KeyID:      input.KeyID,
	}.String()
	if err != nil {
		return nil, err
	}

	return &GenerateOutput{
		Date:      date,
		Host:      input.Host,
		Digest:    digest,
		Signature: signature,
	}, nil
}
