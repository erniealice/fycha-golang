//go:build linux

package pdfconv

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

// zombieWarnThreshold is the instance-wide zombie count at which a conversion
// log line warns that nothing is reaping orphaned LibreOffice helpers. A healthy
// container under tini stays at ~0; the rc-docx-pdf-1 wedge grew ~5 per
// conversion until fork() failed.
const zombieWarnThreshold = 20

var pid1WarnOnce sync.Once

// diagSnapshot reports the instance resources that actually gate fork() for a
// LibreOffice run: live soffice processes, total threads, zombies, the cgroup
// pids ceiling, cgroup memory, and commit (Committed_AS/CommitLimit). It returns
// the rendered fields plus the zombie count so the caller can warn on it.
func diagSnapshot() (string, int) {
	pid1WarnOnce.Do(func() {
		if os.Getpid() == 1 {
			log.Printf("pdfconv: WARNING running as PID 1 without an init; orphaned LibreOffice helpers will never be reaped (run under tini)")
		}
	})
	soffice, threads, zombies := procScan()
	return fmt.Sprintf("boot=%s soffice_procs=%d threads=%d zombies=%d pids=%s mem=%s committed=%s goroutines=%d",
		bootID, soffice, threads, zombies, pidsUsage(), memUsage(), commitUsage(), runtime.NumGoroutine()), zombies
}

// procScan walks /proc once and counts soffice/oosplash processes, total
// threads, and zombies (State Z). Per-pid read errors are ignored because
// processes come and go between the readdir and the read.
func procScan() (soffice, threads, zombies int) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return -1, -1, -1
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if _, err := strconv.Atoi(e.Name()); err != nil {
			continue
		}
		b, err := os.ReadFile("/proc/" + e.Name() + "/status")
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(b), "\n") {
			switch {
			case strings.HasPrefix(line, "Name:"):
				name := strings.TrimSpace(line[len("Name:"):])
				if strings.HasPrefix(name, "soffice") || strings.HasPrefix(name, "oosplash") {
					soffice++
				}
			case strings.HasPrefix(line, "State:"):
				if s := strings.TrimSpace(line[len("State:"):]); len(s) > 0 && s[0] == 'Z' {
					zombies++
				}
			case strings.HasPrefix(line, "Threads:"):
				t, _ := strconv.Atoi(strings.TrimSpace(line[len("Threads:"):]))
				threads += t
			}
		}
	}
	return soffice, threads, zombies
}

// pidsUsage reports cgroup pids.current/max (v2, then v1).
func pidsUsage() string {
	cur := readUint("/sys/fs/cgroup/pids.current")
	max := readCgroupMax("/sys/fs/cgroup/pids.max")
	if cur == 0 {
		cur = readUint("/sys/fs/cgroup/pids/pids.current")
		max = readCgroupMax("/sys/fs/cgroup/pids/pids.max")
	}
	if max == 0 {
		return fmt.Sprintf("%d/max", cur)
	}
	return fmt.Sprintf("%d/%d", cur, max)
}

// memUsage reports memory used vs limit: cgroup v2, then v1, then /proc/meminfo.
func memUsage() string {
	if cur := readUint("/sys/fs/cgroup/memory.current"); cur > 0 {
		return formatUsed(cur, readCgroupMax("/sys/fs/cgroup/memory.max"))
	}
	if cur := readUint("/sys/fs/cgroup/memory/memory.usage_in_bytes"); cur > 0 {
		max := readUint("/sys/fs/cgroup/memory/memory.limit_in_bytes")
		if max >= 1<<62 {
			max = 0
		}
		return formatUsed(cur, max)
	}
	fields := meminfo("MemTotal:", "MemAvailable:")
	if total := fields["MemTotal:"]; total > 0 {
		return formatUsed(total-fields["MemAvailable:"], total) + "(meminfo)"
	}
	return "n/a"
}

// commitUsage reports Committed_AS/CommitLimit; commit (not RSS) gates fork().
func commitUsage() string {
	fields := meminfo("Committed_AS:", "CommitLimit:")
	if fields["CommitLimit:"] == 0 {
		return human(fields["Committed_AS:"]) + "/?"
	}
	return human(fields["Committed_AS:"]) + "/" + human(fields["CommitLimit:"])
}

// meminfo reads the requested /proc/meminfo keys (kB) as bytes.
func meminfo(keys ...string) map[string]uint64 {
	out := make(map[string]uint64, len(keys))
	b, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return out
	}
	for _, line := range strings.Split(string(b), "\n") {
		f := strings.Fields(line)
		if len(f) < 2 {
			continue
		}
		for _, k := range keys {
			if f[0] == k {
				n, _ := strconv.ParseUint(f[1], 10, 64)
				out[k] = n * 1024
			}
		}
	}
	return out
}

func formatUsed(cur, max uint64) string {
	if max == 0 {
		return human(cur)
	}
	return human(cur) + "/" + human(max)
}

func readUint(path string) uint64 {
	b, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	n, _ := strconv.ParseUint(strings.TrimSpace(string(b)), 10, 64)
	return n
}

// readCgroupMax maps the literal "max" (no limit) to 0.
func readCgroupMax(path string) uint64 {
	b, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	s := strings.TrimSpace(string(b))
	if s == "max" {
		return 0
	}
	n, _ := strconv.ParseUint(s, 10, 64)
	return n
}

// human formats a byte count in binary units (e.g. 1.5GB).
func human(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%dB", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(b)/float64(div), "KMGTPE"[exp])
}
