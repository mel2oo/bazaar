package cmd

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"time"
)

func Run(name string, args ...string) error {
	var stderr bytes.Buffer

	cmd := exec.Command(name, args...)
	cmd.Stderr = &stderr

	err := cmd.Start()
	if len(stderr.String()) > 0 {
		return errors.New(stderr.String())
	}

	return err
}

func RunOutput(name string, args ...string) (string, error) {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	cmd := exec.Command(name, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Start()
	if len(stderr.String()) > 0 {
		return stdout.String(), errors.New(stderr.String())
	}

	return stdout.String(), err
}

func RunOutputTimeout(name string, args ...string) (string, error) {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if len(stderr.String()) > 0 {
		return stdout.String(), errors.New(stderr.String())
	}

	return stdout.String(), err
}
