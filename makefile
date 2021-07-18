rsrc:
	rsrc -ico icon.ico -manifest miniwarden.manifest -o rsrc.syso

debug: rsrc
	go build -o bin/Miniwarden.exe .
	bin/Miniwarden

prod: rsrc
	go build -ldflags -H=windowsgui -o bin/Miniwarden.exe .
