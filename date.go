package signature_header

import (
	"time"
)

func Date() string {
	return time.Now().In(time.FixedZone("GMT", 0)).Format(time.RFC1123)
}
