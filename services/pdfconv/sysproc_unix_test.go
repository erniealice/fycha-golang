//go:build unix

package pdfconv

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

// TestProcessGroupKillTakesDownGrandchild proves the timeout path kills the
// whole tree: sh (the "wrapper") backgrounds a long sleep (the "soffice.bin"
// grandchild). Default exec.CommandContext would kill only sh and orphan it.
func TestProcessGroupKillTakesDownGrandchild(t *testing.T) {
	t.Parallel()

	pidFile := filepath.Join(t.TempDir(), "grandchild.pid")
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	started := time.Now()
	_, err := runLibreOffice(ctx, commandSpec{
		binary: "/bin/sh",
		args:   []string{"-c", `sleep 60 & echo $! > "$1"; wait`, "sh", pidFile},
		env:    os.Environ(),
	})
	if err == nil {
		t.Fatal("run returned nil error, want a kill")
	}
	if took := time.Since(started); took > waitDelay/2 {
		t.Fatalf("run took %s; the group kill did not release Wait promptly", took)
	}

	raw, readErr := os.ReadFile(pidFile)
	if readErr != nil {
		t.Fatalf("grandchild pid not recorded: %v", readErr)
	}
	pid, _ := strconv.Atoi(strings.TrimSpace(string(raw)))
	if pid <= 0 {
		t.Fatalf("bad grandchild pid %q", raw)
	}
	// The killed grandchild reparents to init/launchd, which reaps it; poll
	// until the pid is gone (kill(pid, 0) succeeds on a live or unreaped pid).
	deadline := time.Now().Add(5 * time.Second)
	for {
		if errors.Is(syscall.Kill(pid, 0), syscall.ESRCH) {
			return
		}
		if time.Now().After(deadline) {
			_ = syscall.Kill(pid, syscall.SIGKILL)
			t.Fatalf("grandchild %d survived the process-group kill", pid)
		}
		time.Sleep(50 * time.Millisecond)
	}
}
