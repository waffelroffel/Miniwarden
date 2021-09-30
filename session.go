package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"os"
)

type Session struct {
	UserEmail string
	Key       string
	entries   []Entry
	sec       [32]byte
}

func (s *Session) Clear() {
	warning(s.ClearSession())
	s.UserEmail = ""
	s.Key = ""
	s.entries = []Entry{}
	s.sec = [32]byte{}
}

func (s *Session) FetchUserEmail() {
	out, err := cmdStatus()
	fatal(err)
	fatal(json.NewDecoder(&out).Decode(&session))
}

func (s *Session) Sync() {
	warning(cmdSync(s.Key))
}

func (s *Session) FetchAllEntries() {
	if s.Key == "" {
		gwl := GowardLoginWindow{}
		gwl.Start()
	}

	out, err := cmdListItems(s.Key)
	if err != nil {
		s.Key = ""
		s.FetchAllEntries()
		return
	}

	allEntries := []Entry{}
	fatal(json.NewDecoder(&out).Decode(&allEntries))

	for _, e := range allEntries {
		if e.Type == 1 || e.Login.Totp == "" {
			s.entries = append(s.entries, e)
		}
	}

	warning(session.SaveSessionKey())
	session.InitSec()
	session.EncryptAll()
}

func (s *Session) InitSec() {
	_, err := rand.Read(s.sec[:])
	fatal(err)
}

func (s *Session) EncryptAll() {
	for i, entry := range s.entries {
		pt := []byte(entry.Login.Password)

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

func (s *Session) Decrypt(epw string) string {
	nct := []byte(epw)

	block, err := aes.NewCipher(s.sec[:])
	fatal(err)

	aesGCM, err := cipher.NewGCM(block)
	fatal(err)

	nonce, ct := nct[:aesGCM.NonceSize()], nct[aesGCM.NonceSize():]

	pt, err := aesGCM.Open(nil, nonce, ct, nil)
	fatal(err)

	return string(pt)
}

func (s *Session) SaveSessionKey() error {
	if err := os.Mkdir(confDir, 0600); !errors.Is(err, os.ErrExist) {
		return err
	}
	return os.WriteFile(confFile, []byte(s.Key), 0600)
}

func (s *Session) LoadSessionKey() error {
	b, err := os.ReadFile(confFile)
	if err != nil {
		return err
	}
	s.Key = string(b)
	return nil
}

func (s *Session) ClearSession() error {
	if err := os.Mkdir(confDir, 0600); !errors.Is(err, os.ErrExist) {
		return err
	}
	return os.WriteFile(confFile, []byte{}, 0600)
}
