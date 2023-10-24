package signature_header

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

func GenerateKey(bits int) (privateKey *rsa.PrivateKey, privateKeyBytes []byte, publicKeyBytes []byte) {
	privateKey, _ = rsa.GenerateKey(rand.Reader, bits)
	privateKeyBytes = pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)

	publicKeyASN1, _ := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	publicKeyBytes = pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyASN1,
	})
	return
}

func PrivateKeyFromBytes(privateKeyBytes []byte) (*rsa.PrivateKey, error) {
	privateKeyPem, _ := pem.Decode(privateKeyBytes)
	return x509.ParsePKCS1PrivateKey(privateKeyPem.Bytes)
}
