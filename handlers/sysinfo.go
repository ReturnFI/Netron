package handlers

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "os"
    "os/exec"
    "strconv"
    "strings"
    "syscall"

    "netron/models"
)

func getSystemDetails() models.SystemDetails {
    return models.SystemDetails{
        OS:             getOS(),
        Kernel:         getKernel(),
        Arch:           getArch(),
        Uptime:         getUptime(),
        LoadAverage:    getLoadAverage(),
        TCPCongestion:  getTCPCongestion(),
        Virtualization: getVirtualization(),
        IPv4Status:     getIPv4Status(),
        IPv6Status:     getIPv6Status(),
        Organization:   getOrganization(),
        Location:       getLocation(),
        Region:         getRegion(),
        TotalDisk:      getTotalDisk(),
        UsedDisk:       getUsedDisk(),
    }
}

func getOS() string {
    if data, err := ioutil.ReadFile("/etc/os-release"); err == nil {
        lines := strings.Split(string(data), "\n")
        for _, line := range lines {
            if strings.HasPrefix(line, "PRETTY_NAME=") {
                return strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), "\"")
            }
        }
    }
    return "Unknown"
}

func getKernel() string {
    if output, err := exec.Command("uname", "-r").Output(); err == nil {
        return strings.TrimSpace(string(output))
    }
    return "Unknown"
}

func getArch() string {
    if output, err := exec.Command("uname", "-m").Output(); err == nil {
        arch := strings.TrimSpace(string(output))
        bit := "32"
        if strings.Contains(arch, "64") {
            bit = "64"
        }
        return fmt.Sprintf("%s (%s Bit)", arch, bit)
    }
    return "Unknown"
}

func getUptime() string {
    if data, err := ioutil.ReadFile("/proc/uptime"); err == nil {
        fields := strings.Fields(string(data))
        if len(fields) > 0 {
            if uptime, err := strconv.ParseFloat(fields[0], 64); err == nil {
                days := int(uptime) / 86400
                hours := (int(uptime) % 86400) / 3600
                minutes := (int(uptime) % 3600) / 60
                return fmt.Sprintf("%d days, %d hour %d min", days, hours, minutes)
            }
        }
    }
    return "Unknown"
}

func getLoadAverage() string {
    if data, err := ioutil.ReadFile("/proc/loadavg"); err == nil {
        fields := strings.Fields(string(data))
        if len(fields) >= 3 {
            return fmt.Sprintf("%s, %s, %s", fields[0], fields[1], fields[2])
        }
    }
    return "Unknown"
}

func getTCPCongestion() string {
    if output, err := exec.Command("sysctl", "-n", "net.ipv4.tcp_congestion_control").Output(); err == nil {
        return strings.TrimSpace(string(output))
    }
    return "Unknown"
}

func getVirtualization() string {
    if data, err := ioutil.ReadFile("/proc/cpuinfo"); err == nil {
        content := strings.ToLower(string(data))
        if strings.Contains(content, "vmware") {
            return "VMware"
        }
        if strings.Contains(content, "kvm") {
            return "KVM"
        }
    }
    
    if output, err := exec.Command("dmidecode", "-s", "system-product-name").Output(); err == nil {
        product := strings.ToLower(strings.TrimSpace(string(output)))
        if strings.Contains(product, "vmware") {
            return "VMware"
        }
        if strings.Contains(product, "kvm") {
            return "KVM"
        }
        if strings.Contains(product, "virtualbox") {
            return "VirtualBox"
        }
    }
    
    if _, err := os.Stat("/proc/xen"); err == nil {
        return "Xen"
    }
    
    if data, err := ioutil.ReadFile("/proc/1/cgroup"); err == nil {
        content := string(data)
        if strings.Contains(content, "docker") {
            return "Docker"
        }
        if strings.Contains(content, "lxc") {
            return "LXC"
        }
    }
    
    return "Dedicated"
}

func getIPv4Status() string {
    if err := exec.Command("ping", "-4", "-c", "1", "-W", "4", "8.8.8.8").Run(); err == nil {
        return "✓ Online"
    }
    return "✗ Offline"
}

func getIPv6Status() string {
    if err := exec.Command("ping", "-6", "-c", "1", "-W", "4", "2001:4860:4860::8888").Run(); err == nil {
        return "✓ Online"
    }
    return "✗ Offline"
}

func getOrganization() string {
    if output, err := exec.Command("wget", "-q", "-T10", "-O-", "http://ipinfo.io/org").Output(); err == nil {
        return strings.TrimSpace(string(output))
    }
    return "Unknown"
}

func getLocation() string {
    cmd1 := exec.Command("wget", "-q", "-T10", "-O-", "http://ipinfo.io/city")
    cmd2 := exec.Command("wget", "-q", "-T10", "-O-", "http://ipinfo.io/country")
    
    city, err1 := cmd1.Output()
    country, err2 := cmd2.Output()
    
    if err1 == nil && err2 == nil {
        return fmt.Sprintf("%s / %s", strings.TrimSpace(string(city)), strings.TrimSpace(string(country)))
    }
    return "Unknown"
}

func getRegion() string {
    if output, err := exec.Command("wget", "-q", "-T10", "-O-", "http://ipinfo.io/region").Output(); err == nil {
        return strings.TrimSpace(string(output))
    }
    return "Unknown"
}

func getTotalDisk() string {
    var stat syscall.Statfs_t
    if err := syscall.Statfs("/", &stat); err == nil {
        total := stat.Blocks * uint64(stat.Bsize)
        return formatBytes(total)
    }
    return "Unknown"
}

func getUsedDisk() string {
    var stat syscall.Statfs_t
    if err := syscall.Statfs("/", &stat); err == nil {
        total := stat.Blocks * uint64(stat.Bsize)
        available := stat.Bavail * uint64(stat.Bsize)
        used := total - available
        return formatBytes(used)
    }
    return "Unknown"
}

func formatBytes(bytes uint64) string {
    const unit = 1024
    if bytes < unit {
        return fmt.Sprintf("%d B", bytes)
    }
    div, exp := uint64(unit), 0
    for n := bytes / unit; n >= unit; n /= unit {
        div *= unit
        exp++
    }
    return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func getCPUInfoDetailed() models.CPUInfo {
    cpuInfo := models.CPUInfo{
        Usage: getCPUUsage(),
        Cores: getCoreCount(),
    }
    
    if data, err := ioutil.ReadFile("/proc/cpuinfo"); err == nil {
        lines := strings.Split(string(data), "\n")
        for _, line := range lines {
            if strings.Contains(line, "model name") {
                parts := strings.Split(line, ":")
                if len(parts) > 1 {
                    cpuInfo.Model = strings.TrimSpace(parts[1])
                }
            } else if strings.Contains(line, "cpu MHz") {
                parts := strings.Split(line, ":")
                if len(parts) > 1 {
                    if freq := strings.TrimSpace(parts[1]); freq != "" {
                        cpuInfo.Frequency = freq + " MHz"
                    }
                }
            } else if strings.Contains(line, "cache size") {
                parts := strings.Split(line, ":")
                if len(parts) > 1 {
                    cpuInfo.Cache = strings.TrimSpace(parts[1])
                }
            } else if strings.Contains(line, "flags") {
                flags := strings.ToLower(line)
                cpuInfo.AES = strings.Contains(flags, "aes")
                cpuInfo.VMX = strings.Contains(flags, "vmx") || strings.Contains(flags, "svm")
                break
            }
        }
    }
    
    if cpuInfo.Model == "" {
        cpuInfo.Model = "Unknown CPU"
    }
    
    return cpuInfo
}

func getCPUUsage() float64 {
    file, err := os.Open("/proc/stat")
    if err != nil {
        return 0
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    scanner.Scan()
    line := scanner.Text()
    fields := strings.Fields(line)

    if len(fields) < 8 {
        return 0
    }

    idle, _ := strconv.ParseUint(fields[4], 10, 64)
    total := uint64(0)
    for i := 1; i < len(fields); i++ {
        val, _ := strconv.ParseUint(fields[i], 10, 64)
        total += val
    }

    if total == 0 {
        return 0
    }
    
    return float64(total-idle) / float64(total) * 100
}

func getCoreCount() int {
    if data, err := ioutil.ReadFile("/proc/cpuinfo"); err == nil {
        return strings.Count(string(data), "processor")
    }
    return 0
}