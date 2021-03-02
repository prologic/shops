package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	scp "github.com/hnakamur/go-scp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

type FileResult struct {
	err error

	Source string
	Target string
}

func (res *FileResult) Error() error {
	return res.err
}

func (res *FileResult) Ok() bool {
	return res.err == nil
}

func (res *FileResult) String() string {
	var sb strings.Builder

	if res.Ok() {
		sb.WriteString(fmt.Sprintf(" %s -> %s ✅", res.Source, res.Target))
	} else {
		sb.WriteString(fmt.Sprintf(" %s -> %s ❌ (%s)", res.Source, res.Target, res.Error()))
	}

	return sb.String()
}

type ItemResult struct {
	err error

	Name   string
	Output string

	Check  bool
	Action bool
}

func (res *ItemResult) Error() error {
	return res.err
}

func (res *ItemResult) Ok() bool {
	if res.err != nil {
		return false
	}
	return res.Check || (!res.Check && res.Action)
}

func (res *ItemResult) String() string {
	var sb strings.Builder

	if res.Ok() {
		sb.WriteString(fmt.Sprintf(" %s ✅ -> %s", res.Name, res.Output))
	} else {
		sb.WriteString(fmt.Sprintf(" %s ❌ -> %s", res.Name, res.Output))
	}

	return sb.String()
}

type HostResult struct {
	err error

	Addr  string
	Files []FileResult
	Items []ItemResult
}

func (res *HostResult) Error() error {
	return res.err
}

func (res *HostResult) Ok() bool {
	if res.err != nil {
		return false
	}
	for _, file := range res.Files {
		if !file.Ok() {
			return false
		}
	}
	for _, item := range res.Items {
		if !item.Ok() {
			return false
		}
	}
	return true
}

func (res *HostResult) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s:\n", res.Addr))
	if res.Error() != nil {
		sb.WriteString(fmt.Sprintf(" host failed: %s\n", res.Error()))
	} else {
		for _, file := range res.Files {
			sb.WriteString(fmt.Sprintf(" %s\n", file.String()))
		}
		for _, item := range res.Items {
			sb.WriteString(fmt.Sprintf(" %s\n", item.String()))
		}
	}

	return sb.String()
}

type Runner interface {
	Run() error
	Result() *HostResult
}

type SSHRunner struct {
	Addr string
	Conf Config
	User string
	Opts Options

	res *HostResult
}

func NewSSHRunner(addr string, conf Config, user string, opts Options) *SSHRunner {
	runner := &SSHRunner{Addr: addr, Conf: conf, User: user, Opts: opts}
	runner.res = &HostResult{Addr: addr}
	return runner
}

func (run *SSHRunner) Result() *HostResult {
	return run.res
}

func (run *SSHRunner) Run() error {
	failed := func(err error) error {
		run.res.err = err
		return err
	}

	client, _, err := connectToHost(run.User, run.Addr)
	if err != nil {
		return failed(fmt.Errorf("error connecting to host %s: %w", run.Addr, err))
	}

	scpClient := scp.NewSCP(client)
	for _, file := range run.Conf.Files {
		fileInfo, err := os.Stat(file.Source)
		if err != nil {
			run.res.Files = append(run.res.Files, FileResult{err: err, Source: file.Source, Target: file.Target})
			if run.Opts.ContinueOnError {
				continue
			} else {
				return failed(fmt.Errorf("error copying files (aborting)"))
			}
		}

		if fileInfo.IsDir() {
			err = scpClient.SendDir(file.Source, file.Target, nil)
		} else {
			err = scpClient.SendFile(file.Source, file.Target)
		}

		run.res.Files = append(run.res.Files, FileResult{err: err, Source: file.Source, Target: file.Target})
		if !run.Opts.ContinueOnError {
			return failed(fmt.Errorf("error copying files (aborting)"))
		}
	}

	for _, item := range run.Conf.Items {
		out, err := executeRemoteCommand(item.Check, run.Addr, client)
		if err == nil {
			run.res.Items = append(run.res.Items, ItemResult{
				err:    err,
				Name:   item.Name,
				Check:  true,
				Output: strings.TrimSpace(out),
			})
			continue
		}

		if exitError, ok := err.(*ssh.ExitError); ok && exitError.ExitStatus() != 0 {
			out, err := executeRemoteCommand(item.Action, run.Addr, client)
			if err == nil {
				run.res.Items = append(run.res.Items, ItemResult{
					err:    err,
					Name:   item.Name,
					Action: true,
					Output: strings.TrimSpace(out),
				})
				continue
			}

			if exitError, ok := err.(*ssh.ExitError); ok && exitError.ExitStatus() != 0 {
				out += fmt.Sprintf("\nExit status: %d\n", exitError.ExitStatus())
				run.res.Items = append(run.res.Items, ItemResult{
					err:    err,
					Name:   item.Name,
					Action: false,
					Output: strings.TrimSpace(out),
				})
				if !run.Opts.ContinueOnError {
					return failed(fmt.Errorf("error running item (aborting)"))
				}
			}
		} else {
			log.WithError(err).Errorf("error running check %s against %s", item, run.Addr)
			out += fmt.Sprintf("\nExit status: %d\n", exitError.ExitStatus())
			run.res.Items = append(run.res.Items, ItemResult{
				err:    err,
				Name:   item.Name,
				Output: strings.TrimSpace(out),
			})
		}
	}

	return nil
}

type LocalRunner struct {
	Conf Config
	Opts Options

	res *HostResult
}

func NewLocalRunner(conf Config, opts Options) *LocalRunner {
	runner := &LocalRunner{Conf: conf, Opts: opts}
	runner.res = &HostResult{Addr: "local://"}
	return runner
}

func (run *LocalRunner) Result() *HostResult {
	return run.res
}

func (run *LocalRunner) Run() error {
	failed := func(err error) error {
		run.res.err = err
		return err
	}

	for _, file := range run.Conf.Files {
		fileInfo, err := os.Stat(file.Source)
		if err != nil {
			run.res.Files = append(run.res.Files, FileResult{err: err, Source: file.Source, Target: file.Target})
			if run.Opts.ContinueOnError {
				continue
			} else {
				return failed(fmt.Errorf("error copying files (aborting)"))
			}
		}

		if fileInfo.IsDir() {
			err = CopyDirectory(file.Source, file.Target)
		} else {
			_, err = CopyFile(file.Source, file.Target)
		}

		run.res.Files = append(run.res.Files, FileResult{err: err, Source: file.Source, Target: file.Target})
		if !run.Opts.ContinueOnError {
			return failed(fmt.Errorf("error copying files (aborting)"))
		}
	}

	for _, item := range run.Conf.Items {
		out, err := executeLocalCommand(item.Check)
		if err == nil {
			run.res.Items = append(run.res.Items, ItemResult{
				err:    err,
				Name:   item.Name,
				Check:  true,
				Output: strings.TrimSpace(out),
			})
			continue
		}

		if exitError, ok := err.(*ssh.ExitError); ok && exitError.ExitStatus() != 0 {
			out, err := executeLocalCommand(item.Action)
			if err == nil {
				run.res.Items = append(run.res.Items, ItemResult{
					err:    err,
					Name:   item.Name,
					Action: true,
					Output: strings.TrimSpace(out),
				})
				continue
			}

			if exitError, ok := err.(ExitError); ok && exitError.ExitStatus() != 0 {
				out += fmt.Sprintf("\nExit status: %d\n", exitError.ExitStatus())
				run.res.Items = append(run.res.Items, ItemResult{
					err:    err,
					Name:   item.Name,
					Action: false,
					Output: strings.TrimSpace(out),
				})
				if !run.Opts.ContinueOnError {
					return failed(fmt.Errorf("error running item (aborting)"))
				}
			}
		} else {
			log.WithError(err).Errorf("error running check %s against local://", item)
			out += fmt.Sprintf("\nExit status: %d\n", exitError.ExitStatus())
			run.res.Items = append(run.res.Items, ItemResult{
				err:    err,
				Name:   item.Name,
				Output: strings.TrimSpace(out),
			})
		}
	}

	return nil
}

type GroupRunner struct {
	URIs []URI
	Conf Config
	Opts Options
}

func NewGroupRunner(uris []URI, conf Config, opts ...Option) (*GroupRunner, error) {
	options := NewOptions()
	for _, opt := range opts {
		if err := opt(options); err != nil {
			log.WithError(err).Error("error configuring runner")
			return nil, err
		}
	}

	return &GroupRunner{URIs: uris, Conf: conf, Opts: *options}, nil
}

func (run *GroupRunner) Run() {
	var wg sync.WaitGroup

	results := make(chan *HostResult)

	for _, u := range run.URIs {
		var runner Runner

		switch u.Type {
		case "local":
			runner = NewLocalRunner(run.Conf, run.Opts)
		case "ssh":
			runner = NewSSHRunner(u.HostAddr(), run.Conf, u.User, run.Opts)
		default:
			log.WithField("uri", u).Warn("invalid uri")
			continue
		}

		wg.Add(1)
		go func(runner Runner) {
			defer wg.Done()

			if err := runner.Run(); err != nil {
				log.WithError(err).Error("error running host")
			} else {
				results <- runner.Result()
			}
		}(runner)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		fmt.Printf("%s\n", res)
	}
}
