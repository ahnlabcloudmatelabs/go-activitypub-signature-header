package signature_header

import (
	"crypto"
	"encoding/base64"
	"fmt"
)

// Example: signature_header.Digest(crypto.SHA256, []byte("hello world"))
//
// Output: SHA-256=7Uq...Q==
func Digest(hash crypto.Hash, message []byte) string {
	digest := base64.StdEncoding.EncodeToString(hashing(hash, message))

	return fmt.Sprintf("%s=%s", hash.String(), digest)
}
