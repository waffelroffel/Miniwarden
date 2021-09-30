package main

import (
	"bytes"
	"os/exec"
	"syscall"
)

func cmdWithOutput(arg ...string) (bytes.Buffer, error) {
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

func cmdWithoutOutput(arg ...string) error {
	cmd := exec.Command("bw", arg...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}

func cmdLogin(un, pw string) (bytes.Buffer, error) {
	return cmdWithOutput("login", un, pw)
}

func cmdUnlock(pw string) (bytes.Buffer, error) {
	return cmdWithOutput("unlock", pw)
}

func cmdLogout() error {
	return cmdWithoutOutput("logout")
}

func cmdStatus() (bytes.Buffer, error) {
	return cmdWithOutput("status")
}

func cmdListItems(key string) (bytes.Buffer, error) {
	return cmdWithOutput("list", "items", "--session", key)
}

func cmdSync(key string) error {
	return cmdWithoutOutput("sync", "--session", key)
}
