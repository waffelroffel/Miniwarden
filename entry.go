package main

type Uri struct {
	Uri string
}

type Login struct {
	Uris     []Uri
	Username string
	Password string
}

type Entry struct {
	Name  string
	Login Login
	Type  int
}
