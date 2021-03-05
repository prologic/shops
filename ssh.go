package main

import (
	"bytes"
	"fmt"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func connectToHost(user, hostaddr string) (*ssh.Client, *ssh.Session, error) {
	socket := os.Getenv("SSH_AUTH_SOCK")

	conn, err := net.Dial("unix", socket)
	if err != nil {
		log.WithError(err).Error("error connecting to ssh agent")
		return nil, nil, fmt.Errorf("error connecting to ssh agent: %w", err)
	}

	agentClient := agent.NewClient(conn)
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(agentClient.Signers),
		},
		// TODO: This is probably a security risk :/
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", hostaddr, config)
	if err != nil {
		log.WithError(err).Error("error conencting to remote host")
		return nil, nil, fmt.Errorf("error connecting to remote host: %w", err)
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		log.WithError(err).Error("error creating session")
		return nil, nil, fmt.Errorf("error creating session: %w", err)
	}

	return client, session, nil
}

func executeRemoteCommand(command, hostaddr string, client *ssh.Client) (string, error) {
	var (
		err     error
		session *ssh.Session
	)

	if client == nil {
		client, session, err = connectToHost(user, hostaddr)
		if err != nil {
			log.WithError(err).Error("error connecting to host")
			return "", fmt.Errorf("error connecting to host %s: %w", hostaddr, err)
		}
		defer client.Close()
	} else {
		session, err = client.NewSession()
		if err != nil {
			log.WithError(err).Error("error creating new session to host")
			return "", fmt.Errorf("error creating new session to host %s: %w", hostaddr, err)
		}
	}

	var stdout bytes.Buffer

	session.Stdout = &stdout

	var exitError error

	if runErr := session.Run(command); runErr != nil {
		exitError = exitStatus{
			err:    err,
			status: runErr.(*ssh.ExitError).ExitStatus(),
		}
	}

	return stdout.String(), exitError
}
