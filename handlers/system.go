package handlers

import (
    "bufio"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "strconv"
    "strings"

    "netron/models"
)

func GetSystemInfo(w http.ResponseWriter, r *http.Request) {
    speedTestMutex.Lock()
    speedTest := currentTest
    speedTestMutex.Unlock()

    info := models.SystemInfo{
        CPU:       getCPUInfoDetailed(),
        Memory:    getMemoryInfo(),
        Processes: getProcesses(),
        Network:   getNetworkInfo(),
        SpeedTest: speedTest,
        System:    getSystemDetails(),
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(info)
}

func getCPUInfo() models.CPUInfo {
    file, err := os.Open("/proc/stat")
    if err != nil {
        return models.CPUInfo{}
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    scanner.Scan()
    line := scanner.Text()
    fields := strings.Fields(line)

    if len(fields) < 8 {
        return models.CPUInfo{}
    }

    idle, _ := strconv.ParseUint(fields[4], 10, 64)
    total := uint64(0)
    for i := 1; i < len(fields); i++ {
        val, _ := strconv.ParseUint(fields[i], 10, 64)
        total += val
    }

    usage := float64(total-idle) / float64(total) * 100

    cores := 0
    coreFile, err := os.Open("/proc/cpuinfo")
    if err == nil {
        defer coreFile.Close()
        coreScanner := bufio.NewScanner(coreFile)
        for coreScanner.Scan() {
            if strings.Contains(coreScanner.Text(), "processor") {
                cores++
            }
        }
    }

    return models.CPUInfo{
        Usage: usage,
        Cores: cores,
    }
}

func getMemoryInfo() models.MemoryInfo {
    data, err := ioutil.ReadFile("/proc/meminfo")
    if err != nil {
        return models.MemoryInfo{}
    }

    lines := strings.Split(string(data), "\n")
    memInfo := make(map[string]uint64)

    for _, line := range lines {
        fields := strings.Fields(line)
        if len(fields) >= 2 {
            key := strings.TrimSuffix(fields[0], ":")
            val, _ := strconv.ParseUint(fields[1], 10, 64)
            memInfo[key] = val * 1024
        }
    }

    total := memInfo["MemTotal"]
    available := memInfo["MemAvailable"]
    used := total - available
    percent := float64(used) / float64(total) * 100

    return models.MemoryInfo{
        Total:     total,
        Used:      used,
        Available: available,
        Percent:   percent,
    }
}

func getProcesses() []models.ProcessInfo {
    dirs, err := ioutil.ReadDir("/proc")
    if err != nil {
        return []models.ProcessInfo{}
    }

    var processes []models.ProcessInfo
    for _, dir := range dirs {
        if !dir.IsDir() {
            continue
        }

        pid, err := strconv.Atoi(dir.Name())
        if err != nil {
            continue
        }

        statPath := fmt.Sprintf("/proc/%d/stat", pid)
        statData, err := ioutil.ReadFile(statPath)
        if err != nil {
            continue
        }

        fields := strings.Fields(string(statData))
        if len(fields) < 24 {
            continue
        }

        name := strings.Trim(fields[1], "()")
        status := fields[2]

        utime, _ := strconv.ParseUint(fields[13], 10, 64)
        stime, _ := strconv.ParseUint(fields[14], 10, 64)
        totalTime := utime + stime
        cpuUsage := float64(totalTime) / 100.0

        rss, _ := strconv.ParseUint(fields[23], 10, 64)
        memUsage := float64(rss * 4096)

        processes = append(processes, models.ProcessInfo{
            PID:    pid,
            Name:   name,
            CPU:    cpuUsage,
            Memory: memUsage,
            Status: status,
        })

        if len(processes) >= 20 {
            break
        }
    }

    return processes
}