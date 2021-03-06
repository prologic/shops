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
	if e.err != nil {
		return e.err.Error()
	}
	return ""
}

func (e exitStatus) ExitStatus() int {
	return e.status
}
