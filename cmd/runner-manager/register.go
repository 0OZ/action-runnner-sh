package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// runConfigSh invokes the GitHub Actions runner registration script (./config.sh)
// with the provided URL and token, along with any optional configuration.
func runConfigSh(url, token, name, labels, work string, replace, unattended bool) error {
	configScript := "./config.sh"
	if _, err := os.Stat(configScript); os.IsNotExist(err) {
		return fmt.Errorf("%s not found in current directory (run from an extracted runner package)", configScript)
	}

	args := []string{"--url", url, "--token", token}
	if name != "" {
		args = append(args, "--name", name)
	}
	if labels != "" {
		args = append(args, "--labels", labels)
	}
	if work != "" {
		args = append(args, "--work", work)
	}
	if replace {
		args = append(args, "--replace")
	}
	if unattended {
		args = append(args, "--unattended")
	}

	log.Printf("Running: %s %s", configScript, redactTokenArg(args))

	cmd := exec.Command(configScript, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// redactTokenArg produces a printable version of the config.sh argument list
// with the registration token masked.
func redactTokenArg(args []string) string {
	masked := make([]string, len(args))
	copy(masked, args)
	for i, a := range masked {
		if a == "--token" && i+1 < len(masked) {
			masked[i+1] = "***REDACTED***"
		}
	}
	return strings.Join(masked, " ")
}
