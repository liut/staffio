package random

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	idxBits = 6              // 6 bits to represent a letter index
	idxMask = 1<<idxBits - 1 // All 1-bits, as many as idxBits
	idxMax  = 63 / idxBits   // # of letter indices fitting in 63 bits
)

// GenString generate string without number
func GenString(n int) string {
	b := make([]byte, n)
	src := rand.NewSource(time.Now().UnixNano())
	// A src.Int63() generates 63 random bits, enough for idxMax characters!
	for i, cache, remain := n-1, src.Int63(), idxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), idxMax
		}
		if idx := int(cache & idxMask); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= idxBits
		remain--
	}

	return string(b)
}

// GenCode generate string with number
func GenCode() string {
	r := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	d := r.Intn(999999)
	return fmt.Sprintf("%06d", d)
}
