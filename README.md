# Netron - System Monitor Dashboard

A lightweight web-based system monitoring dashboard for Linux servers.

<div align="center">
<img width="680" height="692" alt="image" src="https://github.com/user-attachments/assets/37833e6f-856b-49d9-8ae3-3c34d55571c7" />
</div>

## Quick Start

### One-line Install

```bash
wget https://github.com/ReturnFI/Netron/releases/latest/download/netron-0.0.6-linux-amd64.tar.gz && tar -xzf netron-*-linux-*.tar.gz && ./netron --run
```


## Manual 

### 1. Download

**For AMD64/x86_64 (most common):**
```bash
wget https://github.com/ReturnFI/Netron/releases/latest/download/netron-0.0.6-linux-amd64.tar.gz
```

### 2. Extract

```bash
tar -xzf netron-*-linux-*.tar.gz
```

### 3. Run

```bash
./netron --run
```

### 4. Access

Open your browser: **http://your-server-ip:8080**

## Custom Port

```bash
# Use port 3000
./netron --run -p 3000

# Use port 9090
./netron --run --port 9090
```

## Features

- 📊 **Real-time System Stats** - CPU, Memory, Processes
- 🌐 **Network Monitoring** - Interfaces, TCP/UDP connections
- 🚀 **Speed Test** - Built-in internet speed testing
- 💻 **System Information** - Hardware details, OS info, virtualization
- 🎯 **Single Binary** - No dependencies, just download and run

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



## That's it! 🎉

No configuration needed. No external files. Just one binary.
