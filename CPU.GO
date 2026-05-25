package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// CPUStats holds raw CPU time values read from /proc/stat
type CPUStats struct {
	User    uint64
	Nice    uint64
	System  uint64
	Idle    uint64
	IOWait  uint64
	IRQ     uint64
	SoftIRQ uint64
	UsagePercent float64
	CoreCount    int
}

// readCPURaw reads a single snapshot from /proc/stat
func readCPURaw() (CPUStats, error) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return CPUStats{}, err
	}
	defer f.Close()

	var stats CPUStats
	scanner := bufio.NewScanner(f)
	coreCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}
		if fields[0] == "cpu" {
			stats.User,    _ = parseUint(fields[1])
			stats.Nice,    _ = parseUint(fields[2])
			stats.System,  _ = parseUint(fields[3])
			stats.Idle,    _ = parseUint(fields[4])
			if len(fields) > 5 { stats.IOWait,  _ = parseUint(fields[5]) }
			if len(fields) > 6 { stats.IRQ,     _ = parseUint(fields[6]) }
			if len(fields) > 7 { stats.SoftIRQ, _ = parseUint(fields[7]) }
		}
		// Count CPU cores (lines starting with cpu0, cpu1, ...)
		if len(fields[0]) > 3 && fields[0][:3] == "cpu" && fields[0] != "cpu" {
			coreCount++
		}
	}
	stats.CoreCount = coreCount
	return stats, scanner.Err()
}

// GetCPUStats calculates CPU usage % by taking two snapshots 200ms apart
func GetCPUStats() CPUStats {
	s1, err1 := readCPURaw()
	time.Sleep(200 * time.Millisecond)
	s2, err2 := readCPURaw()

	if err1 != nil || err2 != nil {
		return CPUStats{}
	}

	idle1 := s1.Idle + s1.IOWait
	idle2 := s2.Idle + s2.IOWait

	total1 := s1.User + s1.Nice + s1.System + idle1 + s1.IRQ + s1.SoftIRQ
	total2 := s2.User + s2.Nice + s2.System + idle2 + s2.IRQ + s2.SoftIRQ

	totalDiff := float64(total2 - total1)
	idleDiff  := float64(idle2  - idle1)

	usage := 0.0
	if totalDiff > 0 {
		usage = (totalDiff - idleDiff) / totalDiff * 100.0
	}

	s2.UsagePercent = usage
	return s2
}

func printCPU(s CPUStats) {
	fmt.Println("  ── CPU ─────────────────────────────────────────────")
	fmt.Printf("  Cores   : %d\n", s.CoreCount)
	fmt.Printf("  Usage   : %.1f%%  %s\n", s.UsagePercent, progressBar(s.UsagePercent, 100, 30))
	fmt.Printf("  User    : %.1f%%  |  System: %.1f%%  |  IOWait: %.1f%%\n",
		toPercent(s.User, s.User+s.Nice+s.System+s.Idle+s.IOWait+s.IRQ+s.SoftIRQ),
		toPercent(s.System, s.User+s.Nice+s.System+s.Idle+s.IOWait+s.IRQ+s.SoftIRQ),
		toPercent(s.IOWait, s.User+s.Nice+s.System+s.Idle+s.IOWait+s.IRQ+s.SoftIRQ),
	)
	fmt.Println()
}

func parseUint(s string) (uint64, error) {
	return strconv.ParseUint(strings.TrimSpace(s), 10, 64)
}

func toPercent(part, total uint64) float64 {
	if total == 0 {
		return 0
	}
	return float64(part) / float64(total) * 100.0
}
