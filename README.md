<div align="center">

![cloudmate logo](https://avatars.githubusercontent.com/u/69299682?s=200&v=4)

# ActivityPub Signature header

<small style="opacity: 0.7;">by Cloudmate</small>

---

![Golang](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)

[![Go Test](https://github.com/cloudmatelabs/go-activitypub-signature-header/actions/workflows/gotest.yml/badge.svg)](https://github.com/cloudmatelabs/go-activitypub-signature-header/actions/workflows/gotest.yml)

</div>

## Install

```bash
go get -u github.com/cloudmatelabs/go-activitypub-signature-header
```

## Introduce

This library is generate `Signature` header for the connect with ActivityPub federations.  
And verify the `Signature` header.

## Usage

### Sign `Signature` header

```go
import (
  "crypto"
  "net/url"

  "github.com/go-resty/resty/v2"
  signature_header "github.com/cloudmatelabs/go-activitypub-signature-header"
)

requestURL, _ := url.Parse("https://yodangang.express/users/9iffvxhojp/inbox")
message := []byte(`{
  "@context": "https://www.w3.org/ns/activitystreams",
  "id": "https://snippet.social/@juunini",
  "type": "Follow",
  "actor": "https://snippet.social/@juunini",
  "object": "https://yodangang.express/users/9iffvxhojp"
}`)

headers, err := signature_header.Generate(signature_header.GenerateInput{
  PrivateKeyBytes: []byte("-----BEGIN RSA PRIVATE KEY-----..."),
  // Algorithm:       crypto.SHA256, // optional. if not set, default is crypto.SHA256
  Host:            requestURL.Host,
  Path:            requestURL.Path,
  Body:            message,
  KeyID:           "https://snippet.social/@juunini#main-key",
})
if err != nil {
  // handle error
}

resty.New().R().
  SetBody(message).
  SetHeader("Date", headers.Date).
  SetHeader("Digest", headers.Digest).
  SetHeader("Host", headers.Host).
  SetHeader("Signature", headers.Signature).
  SetHeader("Content-Type", "application/activity+json").
  Post(requestURL.String())
```

### Verify `Signature` header

```go
import (
  signature_header "github.com/cloudmatelabs/go-activitypub-signature-header"
)

verifier := signature_header.Verifier{
  Method: "POST",
  URL: "https://snippet.social/@juunini/inbox",
  Headers: map[string]string{
    "Signature": "...",
    "Host": "...",
    "Date": "...",
    "Digest": "...",
    "Authorization": "...",
    "...": "...",
  },
}

// Recommended
err := verifier.VerifyWithPublicKey(publicKey)
err := verifier.VerifyWithPublicKeyStr(publicKeyStr)

// You can use, but not recommended
err := verifier.VerifyWithActor("https://yodangang.express/@juunini")
err := verifier.VerifyWithBody([]byte("{...}"))
```

### Parse `Signature` header

```go
import (
  signature_header "github.com/cloudmatelabs/go-activitypub-signature-header"
)

// map[string]string
params := signature_header.ParseSignature(signature)
// or given Signature authorization header
// params := signature_header.ParseSignature(authorization)

params["keyId"]
params["algorithm"]
params["headers"]
params["signature"]
```

## License

[MIT](LICENSE)

But, this library use [httpsig].  
[httpsig] is licensed under the [BSD 3-Clause License](https://github.com/go-fed/httpsig/blob/master/LICENSE)

[httpsig]: https://github.com/go-fed/httpsig
