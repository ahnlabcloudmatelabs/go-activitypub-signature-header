<div align="center">

![cloudmate logo](https://avatars.githubusercontent.com/u/69299682?s=200&v=4)

# ActivityPub Signature header

<small style="opacity: 0.7;">by Cloudmate</small>

---

![Golang](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)

</div>

## Install

```bash
go get -u github.com/cloudmatelabs/go-activitypub-signature-header
```

## Introduce

This library is generate `Signature` header for the connect with ActivityPub federations  

## Usage

```go
import (
  "crypto"

  "github.com/go-resty/resty/v2"
  signature_header "github.com/cloudmatelabs/go-activitypub-signature-header"
)

const privateKeyBytes = []byte("-----BEGIN RSA PRIVATE KEY-----...")
const message = []byte(`{
  "@context": "https://www.w3.org/ns/activitystreams",
  "id": "https://snippet.cloudmt.co.kr/@juunini",
  "type": "Follow",
  "actor": "https://snippet.cloudmt.co.kr/@juunini",
  "object": "https://yodangang.express/users/9iffvxhojp"
}`)
const host := "yodangang.express"
const path := "/users/9iffvxhojp/inbox"
const keyID := "https://snippet.cloudmt.co.kr/@juunini#main-key"

privateKey, err := signature_header.PrivateKeyFromBytes(privateKeyBytes)
if err != nil {
  // handle error
}

algorithm := crypto.SHA256
date := signature_header.Date()
digest := signature_header.Digest(algorithm, message)
signature, err := src.Signature{
  PrivateKey: privateKey,
  Algorithm:  algorithm,
  Date:       date,
  Digest:     digest,
  Host:       host,
  Path:       path,
  KeyID:      keyID,
}.String()
if err != nil {
  // handle error
}

resty.New().R().
  SetBody(message).
  SetHeader("Date", date).
  SetHeader("Digest", digest).
  SetHeader("Host", host).
  SetHeader("Signature", signature).
  SetHeader("Content-Type", "application/activity+json").
  Post("https://" + host + path)
```

## License

[MIT](LICENSE)
