package exec

import (
	"bytes"
	"context"
	"os/exec"
	"time"
)

// Exec runs a shell command and returns its stdout, stderr, and exit code.
// If timeoutSeconds > 0, the command will be killed after the specified duration.
func Exec(workdir string, env []map[string]string, command string, timeoutSeconds int) (string, string, int, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	if timeoutSeconds > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	// Use 'bash -c' to allow complex shell syntax like pipes, redirects, etc.
	cmd := exec.CommandContext(ctx, "bash", "-c", command)

	// Set working directory if provided
	if workdir != "" {
		cmd.Dir = workdir
	}

	// add env
	if len(env) > 0 {
		cmd.Env = append([]string{}, cmd.Environ()...)

		for _, e := range env {
			for k, v := range e {
				cmd.Env = append(cmd.Env, k+"="+v)
			}
		}
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	exitCode := 0
	if err != nil {
		// If the command fails, try to extract exit code
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1 // Unexpected error
		}
	}

	return stdout.String(), stderr.String(), exitCode, err
}
