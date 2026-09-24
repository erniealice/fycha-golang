//go:build unix

package pdfconv

import (
	"os/exec"
	"syscall"
)

// setProcessGroup puts the LibreOffice child in its OWN process group (pgid ==
// child pid) so the whole tree can be killed together. On Linux the `soffice`
// on PATH is a shell wrapper that execs oosplash, which forks soffice.bin plus
// helpers; exec.CommandContext's default cancel kills only the direct child and
// orphans the rest, which then pile up under PID 1.
func setProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

// killProcessGroup SIGKILLs the child's entire process group (negative pid).
// It is the cmd.Cancel hook, so a timed-out or cancelled conversion takes its
// whole soffice tree down. A helper that calls setsid() escapes the group; the
// image's init (tini) reaps it once it exits.
func killProcessGroup(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return nil
	}
	return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
}
