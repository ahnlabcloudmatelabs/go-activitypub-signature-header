package signature_header

import "crypto"

func hashing(algorithm crypto.Hash, message []byte) []byte {
	sha := algorithm.New()
	sha.Write(message)
	return sha.Sum(nil)
}
