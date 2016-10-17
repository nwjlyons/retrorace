package retroengine

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

const (
	alphaNumeric = "abcdefghijklmnopqrstuvwxyz0123456789"
	/* Key of length four gives 1,679,616 possible combinations.

	In Python
	>>> len("abcdefghijklmnopqrstuvwxyz0123456789") ** 4
	1,679,616
	*/
	keyLen = 4
)

func NewKey() string {
	key := make([]byte, keyLen)
	for i := range key {
		key[i] = alphaNumeric[rand.Intn(len(alphaNumeric))]
	}
	return string(key)
}
