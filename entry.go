package main

type Uri struct {
	Uri string
}

type Login struct {
	Uris     []Uri
	Username string
	Password string
	Totp     string // ""
}

type Entry struct {
	Name     string
	Login    Login
	Type     int // 1
	Reprompt int // 0
}
