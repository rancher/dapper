package file

import (
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var ran *rand.Rand

func init() {
	ran = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func randString() string {
	b := make([]byte, 7)
	for i := range b {
		b[i] = letters[ran.Intn(len(letters))]
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

func (d *Dapperfile) tempfile(content []byte) (string, error) {
	tempfile, err := os.CreateTemp(".", d.File)
	if err != nil {
		return "", err
	}
	defer tempfile.Close()

	logrus.Debugf("Created tempfile %s", tempfile.Name())

	if _, err := tempfile.Write(content); err != nil {
		return "", err
	}

	return tempfile.Name(), nil
}
