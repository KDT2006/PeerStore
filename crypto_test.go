package main

import (
	"bytes"
	"testing"
)

func TestCopyEncryptDecrypt(t *testing.T) {
	payload := "Some text"
	src := bytes.NewReader([]byte(payload))
	dst := new(bytes.Buffer)
	key := newEncryptionKey()
	_, err := copyEncrypt(key, src, dst)
	if err != nil {
		t.Error(err)
	}

	out := new(bytes.Buffer)
	nw, err := copyDecrypt(key, dst, out)
	if err != nil {
		t.Error(err)
	}

	if nw != 16+len(payload) {
		t.Fail()
	}

	if out.String() != payload {
		t.Errorf("Expected %s, got %s. Decryption failed", payload, out.String())
	}
}
