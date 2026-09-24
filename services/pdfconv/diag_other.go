//go:build !linux

package pdfconv

import "runtime"

// zombieWarnThreshold is unused off Linux (no /proc scan); zombies report -1.
const zombieWarnThreshold = 20

// diagSnapshot is a stub off Linux: the resource probes read /proc and cgroup
// files. boot= is still emitted so log lines share one field set.
func diagSnapshot() (string, int) {
	return "boot=" + bootID + " diag=unavailable(" + runtime.GOOS + ")", -1
}
