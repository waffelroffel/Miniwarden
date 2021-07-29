package main

import (
	"bytes"
	"os/exec"
	"syscall"
)

func cmdAndOutput(arg ...string) (bytes.Buffer, error) {
	cmd := exec.Command("bw", arg...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stderr, err
	}
	return stdout, nil
}

func cmdLogin(un, pw string) (bytes.Buffer, error) {
	return cmdAndOutput("login", un, pw)
}

func cmdUnlock(pw string) (bytes.Buffer, error) {
	return cmdAndOutput("unlock", pw)
}

func cmdLogout() error {
	cmd := exec.Command("bw", "logout")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}

func cmdStatus() (bytes.Buffer, error) {
	return cmdAndOutput("status")
}

func cmdListItems(key string) (bytes.Buffer, error) {
	return cmdAndOutput("list", "items", "--session", key) // bundle with bw.exe
}

/*
func AutoType(un, pw string) {
	cmd := exec.Command("autotype.exe", un, pw)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Start()
}
*/
