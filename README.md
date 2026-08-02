# PingCraft

A simple, lightweight Go wrapper for the official Speedtest CLI.
It automatically executes a speedtest, parses the output, and displays network performance metrics—including latency, jitter, packet loss, download speed, and upload speed—in a clean, readable terminal format.

## Prerequisites

PingCraft relies on the official Speedtest CLI by Ookla. It expects the executable to be named `speedtest` or `speedtest.exe` and will look for it in either:
1. The same directory as the PingCraft executable.
2. Your system's `PATH`.

If it is not found, PingCraft will prompt you to download it from the [official Speedtest website](https://www.speedtest.net/apps/cli).

## Installation

### Option 1: Download Pre-compiled Binaries
Directly from the [Releases page](https://github.com/asong56/pingcraft/releases).

### Option 2: Build from Source
If you prefer to build PingCraft yourself, ensure you have Go 1.21 or later installed.

1. Clone the repository:
```bash
   git clone [https://github.com/asong56/pingcraft.git](https://github.com/asong56/pingcraft.git)
   cd pingcraft
```

2. Build the binary:
```bash
go build -o pingcraft main.go
```



## Usage

Simply run the executable in your terminal. It will automatically accept the license and GDPR agreements for the underlying speedtest CLI and execute the test:

```bash
./pingcraft
```
