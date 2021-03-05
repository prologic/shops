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

	var exitError error

	if waitErr := cmd.Wait(); waitErr != nil {
		exitError = exitStatus{
			err:    waitErr,
			status: waitErr.(*exec.ExitError).Sys().(syscall.WaitStatus).ExitStatus(),
		}
	}

	return buf.String(), exitError
}
