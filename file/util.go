package file

import (
	"math/rand"
	"strings"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randString() string {
	b := make([]byte, 7)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func toMap(str string) map[string]string {
	kv := map[string]string{}

	for _, part := range strings.Fields(str) {
		kvs := strings.SplitN(part, "=", 2)
		if len(kvs) != 2 {
			continue
		}
		kv[kvs[0]] = kvs[1]
	}

	return kv
}
