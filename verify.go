package signature_header

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"strings"

	jsonld_helper "github.com/cloudmatelabs/go-jsonld-helper"
	"github.com/go-fed/httpsig"
)

type Verifier struct {
	Method    string
	URL       string
	Headers   map[string]string
	request   *http.Request
	algorithm httpsig.Algorithm
	verifier  httpsig.Verifier
}

func (v *Verifier) VerifyWithPublicKey(publicKey crypto.PublicKey) error {
	if err := v.init(); err != nil {
		return err
	}

	return v.verifier.Verify(publicKey, v.algorithm)
}

func (v *Verifier) VerifyWithPublicKeyStr(publicKeyStr string) error {
	publicKey, err := publicKeyFromStr(publicKeyStr)
	if err != nil {
		return err
	}

	return v.VerifyWithPublicKey(publicKey)
}

func (v *Verifier) VerifyWithActor(actor string) error {
	publicKey, err := publicKeyFromActor(actor)
	if err != nil {
		return err
	}

	return v.VerifyWithPublicKey(publicKey)
}

func (v *Verifier) VerifyWithBody(body []byte) error {
	actor, err := getActor(body)
	if err != nil {
		return err
	}

	publicKey, err := publicKeyFromActor(actor)
	if err != nil {
		return err
	}

	return v.VerifyWithPublicKey(publicKey)
}

func (v *Verifier) init() (err error) {
	if v.Headers == nil {
		return fmt.Errorf("headers is required")
	}

	v.request, _ = http.NewRequest(v.Method, v.URL, nil)

	for key, value := range v.Headers {
		v.request.Header.Set(key, value)
	}

	v.algorithm, err = algorithmFromHeaders(v.request)
	if err != nil {
		return err
	}

	v.verifier, _ = httpsig.NewVerifier(v.request)
	return
}

func getActor(body []byte) (string, error) {
	jsonld, err := jsonld_helper.ParseJsonLD(body, nil)
	if err != nil {
		return "", err
	}

	return jsonld.ReadKey("actor").StringOrThrow(nil)
}

func publicKeyFromActor(actor string) (crypto.PublicKey, error) {
	publicKeyStr, err := fetchPublicKey(actor)
	if err != nil {
		return nil, err
	}

	return publicKeyFromStr(publicKeyStr)
}

func publicKeyFromStr(publicKeyStr string) (crypto.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyStr))
	return x509.ParsePKIXPublicKey(block.Bytes)
}
func fetchPublicKey(actor string) (string, error) {
	request, err := http.NewRequest("GET", actor, nil)
	if err != nil {
		return "", err
	}

	request.Header.Add("Accept", "application/ld+json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	jsonld, err := jsonld_helper.ParseJsonLD(body, nil)
	if err != nil {
		return "", err
	}

	publicKey, err := jsonld.ReadKey("publicKey").ReadKey("publicKeyPem").StringOrThrow(nil)
	return publicKey, err
}

func algorithmFromHeaders(r *http.Request) (httpsig.Algorithm, error) {
	signature := r.Header.Get("Signature")
	if signature != "" {
		return algorithmFromSignature(string(signature))
	}

	authorization := r.Header.Get("Authorization")
	if authorization != "" {
		return algorithmFromSignature(strings.Replace(string(authorization), "Signature ", "", 1))
	}

	return httpsig.RSA_SHA256, fmt.Errorf("algorithm not found")
}

func algorithmFromSignature(signature string) (httpsig.Algorithm, error) {
	for _, field := range strings.Split(signature, ",") {

		if strings.HasPrefix(field, "algorithm=") {
			return httpsig.Algorithm(strings.ReplaceAll(
				strings.Replace(field, "algorithm=", "", -1),
				`"`,
				"",
			)), nil
		}
	}

	return httpsig.RSA_SHA256, fmt.Errorf("algorithm not found")
}
