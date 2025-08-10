package executer

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/bdoerfchen/webcmd/src/model/config"
)

type executor struct{}

func New() *executor {
	return &executor{}
}

func (e *executor) Execute(ctx context.Context, route config.Route) (result []byte, exitCode int, err error) {
	// Parse command and arguments
	parts := strings.SplitAfterN(route.Command, " ", 2)
	command := strings.TrimSpace(parts[0])

	// Initialize buffer to save command's output in
	outputBuffer := bytes.NewBuffer(make([]byte, 0))

	// Setup command with stdout buffer
	cmd := exec.CommandContext(ctx, command, route.Args...)
	cmd.Stdout = outputBuffer
	cmd.Stderr = cmd.Stdout

	// Add environment variables
	for key, value := range route.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	// Start and wait for command to finish
	if err := cmd.Run(); err != nil {
		if exitErr, isExitErr := err.(*exec.ExitError); isExitErr {
			return outputBuffer.Bytes(), exitErr.ExitCode(), nil
		} else {
			return nil, -1, fmt.Errorf("error during command execution: %w", err)
		}
	}

	// Return output
	return outputBuffer.Bytes(), 0, nil
}
