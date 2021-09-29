rsrc:
	rsrc -ico icon.ico -manifest miniwarden.manifest -o rsrc.syso

debug: rsrc
	go build -o test/Miniwarden-test.exe .
	test/Miniwarden-test

prod: rsrc
	go build -ldflags -H=windowsgui -o bin/Miniwarden.exe .
