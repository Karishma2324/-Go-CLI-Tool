package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const version = "1.0.0"

func main() {
	// Handle Ctrl+C gracefully
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// Parse interval argument (default 2 seconds)
	interval := 2
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--help", "-h":
			printHelp()
			return
		case "--version", "-v":
			fmt.Printf("go-system-monitor v%s\n", version)
			return
		default:
			if n, err := strconv.Atoi(os.Args[1]); err == nil && n > 0 {
				interval = n
			}
		}
	}

	printBanner()
	fmt.Printf("  Refreshing every %d second(s). Press Ctrl+C to exit.\n\n", interval)

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	// Channels for concurrent metric collection
	cpuCh  := make(chan CPUStats, 1)
	memCh  := make(chan MemStats, 1)
	diskCh := make(chan []DiskStats, 1)
	netCh  := make(chan []NetStats, 1)

	// Initial render
	renderAll(cpuCh, memCh, diskCh, netCh)

	for {
		select {
		case <-ticker.C:
			renderAll(cpuCh, memCh, diskCh, netCh)
		case <-sig:
			fmt.Println("\n\n  Goodbye!")
			return
		}
	}
}

// renderAll launches goroutines concurrently and prints results
func renderAll(cpuCh chan CPUStats, memCh chan MemStats, diskCh chan []DiskStats, netCh chan []NetStats) {
	// Launch all collectors concurrently
	go func() { cpuCh <- GetCPUStats() }()
	go func() { memCh <- GetMemStats() }()
	go func() { diskCh <- GetDiskStats() }()
	go func() { netCh <- GetNetStats() }()

	cpu  := <-cpuCh
	mem  := <-memCh
	disk := <-diskCh
	net  := <-netCh

	// Clear screen
	fmt.Print("\033[H\033[2J")

	printBanner()
	printTimestamp()
	printCPU(cpu)
	printMemory(mem)
	printDisk(disk)
	printNetwork(net)
}

func printBanner() {
	fmt.Println("╔══════════════════════════════════════════════════════╗")
	fmt.Println("║         Go System Resource Monitor v" + version + "          ║")
	fmt.Println("║         github.com/Karishma2324                      ║")
	fmt.Println("╚══════════════════════════════════════════════════════╝")
}

func printTimestamp() {
	fmt.Printf("  Time: %s\n\n", time.Now().Format("2006-01-02  15:04:05"))
}

func printHelp() {
	fmt.Println("Go System Resource Monitor")
	fmt.Println("Usage: go-system-monitor [interval_seconds]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -h, --help      Show this help message")
	fmt.Println("  -v, --version   Show version")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  go-system-monitor        # refresh every 2 seconds (default)")
	fmt.Println("  go-system-monitor 5      # refresh every 5 seconds")
}
