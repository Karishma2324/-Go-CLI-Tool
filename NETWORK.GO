package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// NetStats holds per-interface network I/O counters
type NetStats struct {
	Interface string
	RxBytes   uint64
	TxBytes   uint64
	RxPackets uint64
	TxPackets uint64
	RxErrors  uint64
	TxErrors  uint64
}

// GetNetStats reads network interface stats from /proc/net/dev
func GetNetStats() []NetStats {
	f, err := os.Open("/proc/net/dev")
	if err != nil {
		return nil
	}
	defer f.Close()

	var result []NetStats
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		if lineNum <= 2 {
			// Skip header lines
			continue
		}
		line := scanner.Text()
		// Format: iface: rx_bytes rx_packets rx_errs ... tx_bytes tx_packets tx_errs ...
		colonIdx := strings.Index(line, ":")
		if colonIdx < 0 {
			continue
		}
		iface := strings.TrimSpace(line[:colonIdx])
		// Skip loopback
		if iface == "lo" {
			continue
		}

		var rxBytes, rxPackets, rxErrs, rxDrop, rxFifo, rxFrame, rxComp, rxMulti uint64
		var txBytes, txPackets, txErrs, txDrop, txFifo, txColls, txCarr, txComp uint64

		rest := line[colonIdx+1:]
		fmt.Sscanf(strings.TrimSpace(rest),
			"%d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d",
			&rxBytes, &rxPackets, &rxErrs, &rxDrop, &rxFifo, &rxFrame, &rxComp, &rxMulti,
			&txBytes, &txPackets, &txErrs, &txDrop, &txFifo, &txColls, &txCarr, &txComp,
		)

		result = append(result, NetStats{
			Interface: iface,
			RxBytes:   rxBytes,
			TxBytes:   txBytes,
			RxPackets: rxPackets,
			TxPackets: txPackets,
			RxErrors:  rxErrs,
			TxErrors:  txErrs,
		})
	}
	return result
}

func printNetwork(nets []NetStats) {
	fmt.Println("  ── Network ──────────────────────────────────────────")
	if len(nets) == 0 {
		fmt.Println("  No active network interfaces found.")
		fmt.Println()
		return
	}
	for _, n := range nets {
		fmt.Printf("  %-10s  RX: %-12s  TX: %-12s\n",
			n.Interface,
			formatBytes(n.RxBytes),
			formatBytes(n.TxBytes))
		fmt.Printf("             Packets  RX: %-10d  TX: %-10d\n",
			n.RxPackets, n.TxPackets)
		if n.RxErrors > 0 || n.TxErrors > 0 {
			fmt.Printf("             Errors   RX: %-10d  TX: %-10d\n",
				n.RxErrors, n.TxErrors)
		}
	}
	fmt.Println()
}
