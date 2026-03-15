package cmd

import "os/exec"

func runShell(command string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}
