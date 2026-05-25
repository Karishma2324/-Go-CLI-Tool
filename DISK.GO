package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
)

// DiskStats holds per-partition usage and I/O info
type DiskStats struct {
	Device    string
	MountPoint string
	Total     uint64
	Used      uint64
	Free      uint64
	UsedPct   float64
	ReadsKB   uint64
	WritesKB  uint64
}

// GetDiskStats reads mounted partitions and their usage via syscall.Statfs
func GetDiskStats() []DiskStats {
	mounts := getMounts()
	var result []DiskStats

	seen := make(map[string]bool)
	for _, m := range mounts {
		if seen[m.device] {
			continue
		}
		// Only real block devices (skip tmpfs, proc, sys, etc.)
		if !strings.HasPrefix(m.device, "/dev/") {
			continue
		}
		seen[m.device] = true

		var stat syscall.Statfs_t
		if err := syscall.Statfs(m.mountPoint, &stat); err != nil {
			continue
		}

		total := stat.Blocks * uint64(stat.Bsize)
		free  := stat.Bfree  * uint64(stat.Bsize)
		used  := total - free
		pct   := 0.0
		if total > 0 {
			pct = float64(used) / float64(total) * 100.0
		}

		reads, writes := getDiskIO(m.device)

		result = append(result, DiskStats{
			Device:     m.device,
			MountPoint: m.mountPoint,
			Total:      total,
			Used:       used,
			Free:       free,
			UsedPct:    pct,
			ReadsKB:    reads,
			WritesKB:   writes,
		})
	}
	return result
}

type mountEntry struct {
	device     string
	mountPoint string
}

// getMounts parses /proc/mounts
func getMounts() []mountEntry {
	f, err := os.Open("/proc/mounts")
	if err != nil {
		return nil
	}
	defer f.Close()

	var mounts []mountEntry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}
		mounts = append(mounts, mountEntry{
			device:     fields[0],
			mountPoint: fields[1],
		})
	}
	return mounts
}

// getDiskIO reads sectors read/written from /proc/diskstats
func getDiskIO(devicePath string) (readsKB, writesKB uint64) {
	devName := devicePath
	if idx := strings.LastIndex(devicePath, "/"); idx >= 0 {
		devName = devicePath[idx+1:]
	}

	f, err := os.Open("/proc/diskstats")
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 14 {
			continue
		}
		if fields[2] == devName {
			var reads, writes uint64
			fmt.Sscanf(fields[5], "%d", &reads)   // sectors read
			fmt.Sscanf(fields[9], "%d", &writes)  // sectors written
			// 1 sector = 512 bytes = 0.5 KB
			readsKB  = reads  / 2
			writesKB = writes / 2
			return
		}
	}
	return
}

func printDisk(disks []DiskStats) {
	fmt.Println("  ── Disk ─────────────────────────────────────────────")
	if len(disks) == 0 {
		fmt.Println("  No block devices found.")
		fmt.Println()
		return
	}
	for _, d := range disks {
		fmt.Printf("  %-12s  Mount: %-12s\n", d.Device, d.MountPoint)
		fmt.Printf("  Used    : %s / %s (%.1f%%)  %s\n",
			formatBytes(d.Used), formatBytes(d.Total), d.UsedPct,
			progressBar(d.UsedPct, 100, 30))
		fmt.Printf("  Reads   : %s read  |  %s written (lifetime)\n",
			formatKB(d.ReadsKB), formatKB(d.WritesKB))
	}
	fmt.Println()
}

func formatBytes(b uint64) string {
	const (
		_KB = 1024
		_MB = 1024 * _KB
		_GB = 1024 * _MB
		_TB = 1024 * _GB
	)
	switch {
	case b >= _TB:
		return fmt.Sprintf("%.2f TB", float64(b)/float64(_TB))
	case b >= _GB:
		return fmt.Sprintf("%.2f GB", float64(b)/float64(_GB))
	case b >= _MB:
		return fmt.Sprintf("%.1f MB", float64(b)/float64(_MB))
	default:
		return fmt.Sprintf("%.1f KB", float64(b)/float64(_KB))
	}
}
