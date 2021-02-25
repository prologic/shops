package main

import (
	"context"
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

type HostRunner struct {
	Addr string
	Conf Config
	User string

	res *HostResult
}

func NewHostRunner(addr string, conf Config, user string) *HostRunner {
	runner := &HostRunner{Addr: addr, Conf: conf, User: user}
	runner.res = &HostResult{Addr: addr}
	return runner
}

func (run *HostRunner) Result() *HostResult {
	return run.res
}

func (run *HostRunner) Run() error {
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
			continue
		}

		if fileInfo.IsDir() {
			err = scpClient.SendDir(file.Source, file.Target, nil)
		} else {
			err = scpClient.SendFile(file.Source, file.Target)
		}

		run.res.Files = append(run.res.Files, FileResult{err: err, Source: file.Source, Target: file.Target})
	}

	for _, item := range run.Conf.Items {
		out, err := executeCommand(item.Check, run.Addr, client)
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
			out, err := executeCommand(item.Action, run.Addr, client)
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

type GroupRunner struct {
	Addrs []string
	Conf  Config
	User  string

	Debug bool
}

func NewGroupRunner(addrs []string, conf Config, user string, debug bool) *GroupRunner {
	return &GroupRunner{Addrs: addrs, Conf: conf, User: user, Debug: debug}
}

func (run *GroupRunner) Run() {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	results := make(chan *HostResult)

	go func() {
		for {
			select {
			case res, ok := <-results:
				if !ok {
					return
				}
				fmt.Printf("%s\n", res)
			case <-ctx.Done():
				return
			}
		}
	}()

	for _, addr := range run.Addrs {
		runner := NewHostRunner(addr, run.Conf, run.User)
		if debug {
			log.Debugf("created runner for %s", addr)
		}

		wg.Add(1)
		go func(runner *HostRunner) {
			defer wg.Done()

			if debug {
				log.Debugf("running runner for %s", runner.Addr)
			}
			if err := runner.Run(); err != nil {
				log.WithError(err).Error("error running host")
			} else {
				res := runner.Result()
				if debug {
					log.Debugf("result for %s: %s", runner.Addr, res)
				}
				results <- res
			}
		}(runner)
		if debug {
			log.Debugf("started runner for %s", runner.Addr)
		}
	}

	wg.Wait()
	close(results)
}
