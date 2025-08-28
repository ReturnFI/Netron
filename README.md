# Netron - System Monitor Dashboard

A lightweight web-based system monitoring dashboard for Linux servers.

## Quick Start

### 1. Download

**For AMD64/x86_64 (most common):**
```bash
wget https://github.com/ReturnFI/Netron/releases/latest/download/netron-0.0.4-linux-amd64.tar.gz
```

**For ARM64 (Raspberry Pi 4, Apple Silicon servers):**
```bash
wget https://github.com/ReturnFI/Netron/releases/latest/download/netron-0.0.4-linux-arm64.tar.gz
```

**For ARMv7 (Raspberry Pi 3):**
```bash
wget https://github.com/ReturnFI/Netron/releases/latest/download/netron-0.0.4-linux-armv7.tar.gz
```

**For ARMv6 (Raspberry Pi Zero):**
```bash
wget https://github.com/ReturnFI/Netron/releases/latest/download/netron-0.0.4-linux-armv6.tar.gz
```

**For 32-bit x86:**
```bash
wget https://github.com/ReturnFI/Netron/releases/latest/download/netron-0.0.4-linux-i386.tar.gz
```

### 2. Extract

```bash
tar -xzf netron-*-linux-*.tar.gz
```

### 3. Run

```bash
./netron-*-linux-* --run
```

### 4. Access

Open your browser: **http://your-server-ip:8080**

## Custom Port

```bash
# Use port 3000
./netron-*-linux-* --run -p 3000

# Use port 9090
./netron-*-linux-* --run --port 9090
```

## Features

- 📊 **Real-time System Stats** - CPU, Memory, Processes
- 🌐 **Network Monitoring** - Interfaces, TCP/UDP connections
- 🚀 **Speed Test** - Built-in internet speed testing
- 💻 **System Information** - Hardware details, OS info, virtualization
- 🎯 **Single Binary** - No dependencies, just download and run

## One-line Install

```bash
wget https://github.com/ReturnFI/Netron/releases/latest/download/netron-0.0.4-linux-amd64.tar.gz && tar -xzf netron-*-linux-*.tar.gz && ./netron-*-linux-* --run
```

## Architecture Guide

Not sure which version to download?

```bash
# Check your architecture
uname -m
```

- `x86_64` → use `linux-amd64`
- `aarch64` or `arm64` → use `linux-arm64`
- `armv7l` → use `linux-armv7`
- `armv6l` → use `linux-armv6`
- `i386` or `i686` → use `linux-i386`

## Optional: Install Speed Test

For internet speed testing feature:
```bash
# Ubuntu/Debian
apt update && apt install speedtest-cli

# CentOS/RHEL
yum install speedtest-cli
```

## That's it! 🎉

No configuration needed. No external files. Just one binary.
