package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"io"
)

func generateID() string {
	buf := make([]byte, 32)
	io.ReadFull(rand.Reader, buf)
	return hex.EncodeToString(buf)
}

func newEncryptionKey() []byte {
	keyBuf := make([]byte, 32)
	io.ReadFull(rand.Reader, keyBuf)
	return keyBuf
}

func hashKey(key string) string {
	hash := sha1.Sum([]byte(key))
	return hex.EncodeToString(hash[:])
}

func copyStream(stream cipher.Stream, BlockSize int, src io.Reader, dst io.Writer) (int, error) {
	buf := make([]byte, 32*1024)
	nw := BlockSize

	for {
		n, err := src.Read(buf)
		if n > 0 {
			stream.XORKeyStream(buf, buf[:n])
			size, err := dst.Write(buf[:n])
			if err != nil {
				return 0, nil
			}
			nw += size
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return 0, err
		}
	}

	return nw, nil
}

func copyDecrypt(key []byte, src io.Reader, dst io.Writer) (int, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, err
	}

	// Read the IV from the given io.Reader which in our case should be
	// the block.BlockSize() bytes we read
	iv := make([]byte, block.BlockSize())
	if _, err := src.Read(iv); err != nil {
		return 0, err
	}

	stream := cipher.NewCTR(block, iv)
	return copyStream(stream, block.BlockSize(), src, dst)
}

func copyEncrypt(key []byte, src io.Reader, dst io.Writer) (int, error) {
	// Create a new AES cipher block using the provided key
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, err
	}

	// Create an initialization vector (IV) with the same block size as the cipher
	iv := make([]byte, block.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return 0, err
	}

	// Write the IV to the beginning of the destination file
	if _, err := dst.Write(iv); err != nil {
		return 0, err
	}

	// Create a new stream cipher using AES in CTR mode with the given block and IV
	stream := cipher.NewCTR(block, iv)
	return copyStream(stream, block.BlockSize(), src, dst)
}
