package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func main() {
	key, err := Salt(39)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	} else {
		fmt.Printf("new key: %q\n", base64.URLEncoding.EncodeToString(key))
	}

}

func Salt(strength int) (k []byte, err error) {
	k = make([]byte, strength)
	_, err = io.ReadFull(rand.Reader, k)
	return
}
