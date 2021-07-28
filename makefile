rsrc:
	rsrc -ico icon.ico -manifest miniwarden.manifest -o rsrc.syso

ahk:
	ahk2exe /in autotype.ahk /out "bin\autotype.exe"

debug: rsrc ahk
	go build -o bin/Miniwarden.exe .
	bin/Miniwarden

prod: rsrc ahk
	go build -ldflags -H=windowsgui -o bin/Miniwarden.exe .
