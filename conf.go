package main

import (
	"errors"
	"os"
)

func (s *Session) SaveToDisk() {
	if err := os.Mkdir(confDir, 0600); !errors.Is(err, os.ErrExist) {
		fatal(err)
	}
	fatal(os.WriteFile(confFile, []byte(s.Key), 0600))
}

func (s *Session) LoadFromDisk() {
	if b, err := os.ReadFile(confFile); err == nil {
		s.Key = string(b)
	}
}

func (s *Session) ClearFromDisk() {
	if err := os.Mkdir(confDir, 0600); !errors.Is(err, os.ErrExist) {
		fatal(err)
	}
	fatal(os.WriteFile(confFile, []byte{}, 0600))
}
