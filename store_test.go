package main

import (
	"bytes"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "bestpics"
	pathKey := CASPathTransformFunc(key)
	// fmt.Println(pathname)
	expectedOriginalKey := "c565996f77ccab3a98f55f6546faa5b311ea674b"
	expectedPathName := "c5659/96f77/ccab3/a98f5/5f654/6faa5/b311e/a674b"
	if pathKey.PathName != expectedPathName {
		t.Errorf("have %s, want %s", pathKey.PathName, expectedPathName)
	}
	if pathKey.FileName != expectedOriginalKey {
		t.Errorf("have %s, want %s", pathKey.FileName, expectedOriginalKey)
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "specialpics"
	data := []byte("some image")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	r, err := s.ReadStream(key)
	if err != nil {
		t.Error(err)
	}

	b, err := io.ReadAll(r)
	if err != nil {
		t.Error(err)
	}

	if string(b) != string(data) {
		t.Errorf("want %s, have %s", data, b)
	}

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}

}

func TestDelete(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "specialpics"
	data := []byte("some image")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}
