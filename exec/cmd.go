package exec

import (
	"bytes"
	"os/exec"
)

// execCmd runs a shell command and returns its stdout, stderr, and exit code.
func Exec(workdir, command string) (string, string, int, error) {
	// Use 'bash -c' to allow complex shell syntax like pipes, redirects, etc.
	cmd := exec.Command("bash", "-c", command)

	// Set working directory if provided
	if workdir != "" {
		cmd.Dir = workdir
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
