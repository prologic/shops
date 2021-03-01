package main

import (
	"bytes"
	"os/exec"
	"syscall"
)

func executeLocalCommand(command string) (string, error) {
	var buf bytes.Buffer

	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	if err := cmd.Start(); err != nil {
		return "", err
	}

	var err error

	if err := cmd.Wait(); err != nil {
		err = exitStatus{err: err, status: err.(*exec.ExitError).Sys().(syscall.WaitStatus).ExitStatus()}
	}

	return buf.String(), err
}
