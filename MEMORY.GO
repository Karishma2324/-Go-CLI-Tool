package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// MemStats holds memory information in KB
type MemStats struct {
	Total     uint64
	Free      uint64
	Available uint64
	Buffers   uint64
	Cached    uint64
	SwapTotal uint64
	SwapFree  uint64
	Used      uint64
	SwapUsed  uint64
}

// GetMemStats reads memory info from /proc/meminfo
func GetMemStats() MemStats {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return MemStats{}
	}
	defer f.Close()

	fields := make(map[string]uint64)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		key := strings.TrimSuffix(parts[0], ":")
		val, _ := strconv.ParseUint(parts[1], 10, 64)
		fields[key] = val
	}

	mem := MemStats{
		Total:     fields["MemTotal"],
		Free:      fields["MemFree"],
		Available: fields["MemAvailable"],
		Buffers:   fields["Buffers"],
		Cached:    fields["Cached"],
		SwapTotal: fields["SwapTotal"],
		SwapFree:  fields["SwapFree"],
	}
	mem.Used = mem.Total - mem.Available
	if mem.SwapTotal > mem.SwapFree {
		mem.SwapUsed = mem.SwapTotal - mem.SwapFree
	}
	return mem
}

func printMemory(m MemStats) {
	usedPct := 0.0
	if m.Total > 0 {
		usedPct = float64(m.Used) / float64(m.Total) * 100.0
	}
	swapPct := 0.0
	if m.SwapTotal > 0 {
		swapPct = float64(m.SwapUsed) / float64(m.SwapTotal) * 100.0
	}

	fmt.Println("  ── Memory ───────────────────────────────────────────")
	fmt.Printf("  RAM     : %s / %s (%.1f%%)  %s\n",
		formatKB(m.Used), formatKB(m.Total), usedPct,
		progressBar(usedPct, 100, 30))
	fmt.Printf("  Buffers : %s   Cached: %s   Free: %s\n",
		formatKB(m.Buffers), formatKB(m.Cached), formatKB(m.Free))
	fmt.Printf("  Swap    : %s / %s (%.1f%%)  %s\n",
		formatKB(m.SwapUsed), formatKB(m.SwapTotal), swapPct,
		progressBar(swapPct, 100, 30))
	fmt.Println()
}

// formatKB converts KB to human-readable string
func formatKB(kb uint64) string {
	const (
		_MB = 1024
		_GB = 1024 * 1024
	)
	switch {
	case kb >= _GB:
		return fmt.Sprintf("%.2f GB", float64(kb)/float64(_GB))
	case kb >= _MB:
		return fmt.Sprintf("%.1f MB", float64(kb)/float64(_MB))
	default:
		return fmt.Sprintf("%d KB", kb)
	}
}
