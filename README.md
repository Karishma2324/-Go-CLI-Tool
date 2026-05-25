# 🖥️ Go System Resource Monitor

A real-time command-line system resource monitoring tool built in **Go**, reading directly from the Linux `/proc` filesystem — no external dependencies.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![Linux](https://img.shields.io/badge/Platform-Linux-FCC624?style=flat&logo=linux)
![License](https://img.shields.io/badge/License-MIT-green)

---

## 📸 Screenshot

```
╔══════════════════════════════════════════════════════╗
║         Go System Resource Monitor v1.0.0            ║
║         github.com/Karishma2324                      ║
╚══════════════════════════════════════════════════════╝
  Time: 2025-06-01  14:32:10

  ── CPU ─────────────────────────────────────────────
  Cores   : 4
  Usage   : 23.4%  [███████░░░░░░░░░░░░░░░░░░░░░░░]
  User    : 18.1%  |  System: 5.3%  |  IOWait: 0.2%

  ── Memory ───────────────────────────────────────────
  RAM     : 4.21 GB / 15.54 GB (27.1%)  [████████░░░░░░░░░░░░░░░░░░░░░░]
  Buffers : 312.0 MB   Cached: 2.10 GB   Free: 8.90 GB
  Swap    : 0.00 KB / 2.00 GB (0.0%)    [░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░]

  ── Disk ─────────────────────────────────────────────
  /dev/sda1    Mount: /
  Used    : 42.30 GB / 230.00 GB (18.4%)  [█████░░░░░░░░░░░░░░░░░░░░░░░░░]
  Reads   : 12.30 GB read  |  8.50 GB written (lifetime)

  ── Network ──────────────────────────────────────────
  eth0        RX: 1.23 GB       TX: 456.7 MB
              Packets  RX: 1234567     TX: 987654
```

---

## ✨ Features

- **CPU monitoring** — real-time usage %, user/system/iowait breakdown
- **Memory monitoring** — RAM and swap usage with human-readable sizes
- **Disk I/O monitoring** — per-partition usage + lifetime reads/writes from `/proc/diskstats`
- **Network monitoring** — per-interface RX/TX bytes, packets, and errors from `/proc/net/dev`
- **Go concurrency** — all four metrics collected simultaneously via goroutines and channels
- **Zero dependencies** — reads directly from Linux `/proc` filesystem
- **ASCII progress bars** — instant visual feedback
- **Configurable refresh interval** — pass any interval in seconds as argument

---

## 🏗️ Architecture

```
main.go          →  entry point, signal handling, ticker loop
cpu.go           →  reads /proc/stat, calculates usage via two snapshots
memory.go        →  reads /proc/meminfo, parses RAM and swap
disk.go          →  reads /proc/mounts + /proc/diskstats via syscall.Statfs
network.go       →  reads /proc/net/dev
utils.go         →  shared helpers (progress bar, formatters)
```

All collectors run as **goroutines** and send results through **channels**, so CPU, memory, disk, and network data are all fetched concurrently on each refresh tick.

---

## 🚀 Getting Started

### Prerequisites
- Go 1.21+
- Linux (reads from `/proc` filesystem)

### Install & Run

```bash
# Clone the repo
git clone https://github.com/Karishma2324/go-system-monitor.git
cd go-system-monitor

# Run directly
go run .

# Or build a binary
go build -o go-system-monitor .
./go-system-monitor

# Custom refresh interval (e.g. every 5 seconds)
./go-system-monitor 5
```

### Options

```
go-system-monitor [interval_seconds]

  -h, --help      Show help
  -v, --version   Show version
  (no args)       Refresh every 2 seconds (default)
```

---

## 📁 Project Structure

```
go-system-monitor/
├── main.go       # Entry point, goroutine orchestration
├── cpu.go        # CPU stats from /proc/stat
├── memory.go     # Memory stats from /proc/meminfo
├── disk.go       # Disk stats from /proc/diskstats
├── network.go    # Network stats from /proc/net/dev
├── utils.go      # Progress bar & formatters
├── go.mod        # Go module definition
└── README.md
```

---

## 🔧 How It Works

This tool reads raw data directly from the Linux **`/proc` virtual filesystem** — the same source used by tools like `top`, `htop`, and `vmstat`.

| Metric  | Source                  |
|---------|-------------------------|
| CPU     | `/proc/stat`            |
| Memory  | `/proc/meminfo`         |
| Disk    | `/proc/mounts` + `/proc/diskstats` + `syscall.Statfs` |
| Network | `/proc/net/dev`         |

CPU usage is calculated by taking **two snapshots 200ms apart** and computing the ratio of non-idle ticks to total ticks — the standard method used by system monitors.

---

## 🤝 Author

**Mohammad Karishma** — [github.com/Karishma2324](https://github.com/Karishma2324) | [LinkedIn](https://www.linkedin.com/in/karishma-mohammad2324/)

---

## 📄 License

MIT License — free to use and modify.
