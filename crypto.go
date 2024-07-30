package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

func newEncryptionKey() []byte {
	keyBuf := make([]byte, 32)
	io.ReadFull(rand.Reader, keyBuf)
	return keyBuf
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

	buf := make([]byte, 32*1024)
	stream := cipher.NewCTR(block, iv)
	nw := block.BlockSize()

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

	buf := make([]byte, 32*1024)
	// Create a new stream cipher using AES in CTR mode with the given block and IV
	stream := cipher.NewCTR(block, iv)
	nw := block.BlockSize()

	for {
		// Read data from the source into the buffer
		n, err := src.Read(buf)
		if n > 0 {
			// Encrypt the data in the buffer using the stream cipher
			stream.XORKeyStream(buf, buf[:n])
			// Write the encrypted data to the destination
			size, err := dst.Write(buf[:n])
			if err != nil {
				return 0, err
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
