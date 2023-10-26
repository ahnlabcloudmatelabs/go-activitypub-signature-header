package signature_header_test

import (
	"fmt"
	"strings"
	"testing"

	signature_header "github.com/cloudmatelabs/go-activitypub-signature-header"
)

var (
	keyId         = `https://snippet.social/@juunini#main-key`
	algorithm     = `rsa-sha256`
	_headers      = `(request-target) date host digest`
	signatureHash = "g99dh20dygasd=="
)

var signature = strings.Join([]string{
	fmt.Sprintf(`keyId="%s"`, keyId),
	fmt.Sprintf(`algorithm="%s"`, algorithm),
	fmt.Sprintf(`headers="%s"`, headers),
	fmt.Sprintf(`signature="%s"`, signatureHash),
}, ",")

func Test_ParseSignature(t *testing.T) {
	params := signature_header.ParseSignature(signature)

	if params["algorithm"] != algorithm {
		t.Errorf("Algorithm is not correct: %s", params["algorithm"])
	}

	if params["keyId"] != keyId {
		t.Errorf("KeyID is not correct: %s", params["keyId"])
	}

	if params["signature"] != signatureHash {
		t.Errorf("Signature is not correct: %s", params["signature"])
	}

	if params["headers"] != _headers {
		t.Errorf("Headers is not correct: %s", params["headers"])
	}
}

func Test_ParseSignature_WithAuthorization(t *testing.T) {
	params := signature_header.ParseSignature("Signature " + signature)

	if params["algorithm"] != algorithm {
		t.Errorf("Algorithm is not correct: %s", params["algorithm"])
	}

	if params["keyId"] != keyId {
		t.Errorf("KeyID is not correct: %s", params["keyId"])
	}

	if params["signature"] != signatureHash {
		t.Errorf("Signature is not correct: %s", params["signature"])
	}

	if params["headers"] != _headers {
		t.Errorf("Headers is not correct: %s", params["headers"])
	}
}
