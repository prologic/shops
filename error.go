package main

type ExitError interface {
	Error() string
	ExitStatus() int
}

type exitStatus struct {
	err    error
	status int
}

func (e exitStatus) Error() string {
	return e.err.Error()
}

func (e exitStatus) ExitStatus() int {
	return e.status
}
