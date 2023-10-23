package signature_header

import "time"

func Date() string {
	return time.Now().In(time.UTC).Format(time.RFC1123)
}
