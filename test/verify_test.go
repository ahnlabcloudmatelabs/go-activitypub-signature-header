package signature_header_test

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"

	signature_header "github.com/cloudmatelabs/go-activitypub-signature-header"
)

var wait sync.WaitGroup

func TestMain(m *testing.M) {
	wait.Add(1)
	setup()
	wait.Wait()

	code := m.Run()
	shutdown()
	os.Exit(code)
}

var headers = createSignature()
var block, _ = pem.Decode([]byte(publicKeyStr))
var publicKey, _ = x509.ParsePKIXPublicKey(block.Bytes)

func Test_VerifyWithPublicKey(t *testing.T) {
	verify := signature_header.Verifier{
		Headers: headers,
		Method:  "POST",
		URL:     "https://snippet.social/@juunini/inbox",
	}

	if err := verify.VerifyWithPublicKey(publicKey); err != nil {
		t.Error(err)
	}
}

func Test_VerifyWithPublicKeyStr(t *testing.T) {
	verify := signature_header.Verifier{
		Headers: headers,
		Method:  "POST",
		URL:     "https://snippet.social/@juunini/inbox",
	}

	if err := verify.VerifyWithPublicKeyStr(publicKeyStr); err != nil {
		t.Error(err)
	}
}

func Test_VerifyWithActor(t *testing.T) {
	verify := signature_header.Verifier{
		Headers: headers,
		Method:  "POST",
		URL:     "https://snippet.social/@juunini/inbox",
	}

	if err := verify.VerifyWithActor("http://localhost:8000/@juunini"); err != nil {
		t.Error(err)
	}
}

func Test_VerifyWithBody(t *testing.T) {
	verify := signature_header.Verifier{
		Headers: headers,
		Method:  "POST",
		URL:     "https://snippet.social/@juunini/inbox",
	}

	if err := verify.VerifyWithBody([]byte(requestMessage)); err != nil {
		t.Error(err)
	}
}

func createSignature() map[string]string {
	message := []byte(requestMessage)
	const host = "snippet.social"
	const path = "/@juunini/inbox"
	const keyID = "http://localhost:8000/@juunini#main-key"
	const algorithm = crypto.SHA256

	headers, err := signature_header.Generate(signature_header.GenerateInput{
		PrivateKeyBytes: []byte(privateKeyStr),
		Algorithm:       algorithm,
		Host:            host,
		Path:            path,
		Body:            message,
		KeyID:           keyID,
	})
	if err != nil {
		panic(err)
	}

	return map[string]string{
		"Signature":    headers.Signature,
		"Date":         headers.Date,
		"Host":         headers.Host,
		"Digest":       headers.Digest,
		"Content-Type": "application/activity+json",
	}
}

var srv *http.Server

func setup() {
	srv = &http.Server{Addr: "localhost:8000"}

	srv.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"")
		io.WriteString(w, userJSON)
	})

	go func() {
		srv.ListenAndServe()
	}()

	for {
		resp, err := http.Get("http://localhost:8000/@juunini")
		if err == nil {
			resp.Body.Close()
			wait.Done()
			break
		}
	}
}

func shutdown() {
	srv.Shutdown(context.TODO())
}

var publicKeyStr = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDQ/IVYdqmPXub2a3DEnYoJJw3c
QY/C13+2Xi9PAZkszfVvsW8woB+MLi0krM0d0cvn5VZEyOxuQBLYWMMY4i/GrJSs
GqU4eBnpoT1+LDGqbeemu0FYNQvkI2ogVVToZnjcXrlIYM0CCXiv/TEIkI+Cfyk1
gIiLoNn0jdL1n7cZKwIDAQAB
-----END PUBLIC KEY-----
`
var privateKeyStr = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDQ/IVYdqmPXub2a3DEnYoJJw3cQY/C13+2Xi9PAZkszfVvsW8w
oB+MLi0krM0d0cvn5VZEyOxuQBLYWMMY4i/GrJSsGqU4eBnpoT1+LDGqbeemu0FY
NQvkI2ogVVToZnjcXrlIYM0CCXiv/TEIkI+Cfyk1gIiLoNn0jdL1n7cZKwIDAQAB
AoGBANB2pBjrPRYrh13VLIaj0xZwI45Kw7mKtvVWHADGSmH8DFBfANfTUcTGJvfH
e0+8f5aRGK3Ccr1DIsk2OV5v+VB542bd+3ZZiHmFBFp05pz7SynO2XQy6vlpzRUj
m+mkodxcnRtBRq2CC19KgVhrHx1hG/5pyRh9mrpVg5aSkCb5AkEA3ZRRvM1WGR4x
wG7+KA5tEXogyWYNPebbi2PON27JIOoyjNxYaopCiobI+qNFK4SnlbQZfinlB2Lo
B3zG5/LqfwJBAPFzaN6xmKkODoB4SaZEZpfGKkh7L7ltkHG/dA3Q1HZhPuhQw9fJ
aXRwVSzSCfF+wDQ6WBHU/KuRW1ZGqL9QQ1UCQDrp4aWybs704UOJ/1eFJmi8MRV7
Zc/snrj8C2tfsGho/JHJUFTbd/+/AJbrbEu61JgQL6sE1plVKd47xeMMCl8CQC8n
aeCr8HN7oktmsoN9MkgL1HApVq2w/xen20NjeErSPRXjyAuZczXhRlElh/mY1nKc
vlxlKx9amOrli8kpJK0CQQDMGOYc4Z91c23sTzSDxTCAujV26+nHc1WSU+x+FqUY
3gqSqQp2pmV2AO5U+X2x2152thp3YsEb14JMmQTT4/Ly
-----END RSA PRIVATE KEY-----
`
var userJSON = fmt.Sprintf(`{
	"@context": [
			"activitystreams.json",
			"security.json"
	],
	"id": "http://localhost:8000/@juunini",
	"inbox": "http://localhost:8000/@juunini/inbox",
	"publicKey": {
			"id": "http://localhost:8000/@juunini#main-key",
			"owner": "http://localhost:8000/@juunini",
			"publicKeyPem": "%s",
			"type": "Key"
	},
	"type": "Person",
	"url": "http://localhost:8000/@juunini"
}`, strings.ReplaceAll(publicKeyStr, "\n", "\\n"))
var requestMessage = `{
	"@context": "activitystreams.json",
	"id": "http://localhost:8000/@juunini",
	"type": "Follow",
	"actor": "http://localhost:8000/@juunini",
	"object": "https://snippet.social/@juunini"
}`
