//go:build !unix

package pdfconv

import "os/exec"

// Non-unix builds (Windows dev boxes) have no process groups via SysProcAttr;
// fall back to killing the direct child so the package still compiles.
func setProcessGroup(cmd *exec.Cmd) {}

func killProcessGroup(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return nil
	}
	return cmd.Process.Kill()
}
