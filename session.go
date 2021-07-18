package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"io"
)

type Session struct {
	UserEmail string
	Key       string
	entries   Entries
	sec       [32]byte
}

func (s *Session) Clear() {
	s.ClearFromDisk()
	s.UserEmail = ""
	s.Key = ""
	s.entries = Entries{}
	s.sec = [32]byte{}
}

func (s *Session) FetchUserEmail() {
	out, err := cmdStatus()
	fatal(err)
	fatal(json.NewDecoder(&out).Decode(&session))
}

func (s *Session) FetchAllEntries() {
	if s.Key == "" {
		gwl := GowardLoginWindow{}
		gwl.Start()
	}

	out, err := cmdListItems(s.Key)
	if err != nil {
		s.Key = "" // add logic for incorrect s.key
		return
	}
	fatal(json.NewDecoder(&out).Decode(&s.entries))
	session.InitSec()
	session.EncryptAll()
}

func (s *Session) InitSec() {
	_, err := rand.Read(s.sec[:])
	fatal(err)
}

func (s *Session) EncryptAll() {
	for i := 0; i < s.entries.Len(); i++ {
		pt := []byte(s.entries[i].Login.Password)

		block, err := aes.NewCipher(s.sec[:])
		fatal(err)

		aesGCM, err := cipher.NewGCM(block)
		fatal(err)

		nonce := make([]byte, aesGCM.NonceSize())
		_, err = io.ReadFull(rand.Reader, nonce)
		fatal(err)

		ciphertext := aesGCM.Seal(nonce, nonce, pt, nil)
		s.entries[i].Login.Password = string(ciphertext)
	}
}

func (s *Session) Decrypt(ncts string) string {
	nct := []byte(ncts)

	block, err := aes.NewCipher(s.sec[:])
	fatal(err)

	aesGCM, err := cipher.NewGCM(block)
	fatal(err)

	nonce, ct := nct[:aesGCM.NonceSize()], nct[aesGCM.NonceSize():]

	pt, err := aesGCM.Open(nil, nonce, ct, nil)
	fatal(err)

	return string(pt)
}
