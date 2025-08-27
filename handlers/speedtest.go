package handlers

import (
    "encoding/json"
    "net/http"
    "os/exec"
    "strconv"
    "strings"
    "sync"
    "time"

    "netron/models"
)

var (
    speedTestMutex sync.Mutex
    currentTest    models.SpeedTestInfo
    isRunning      bool
)

func GetSpeedTest(w http.ResponseWriter, r *http.Request) {
    speedTestMutex.Lock()
    defer speedTestMutex.Unlock()

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(currentTest)
}

func StartSpeedTest(w http.ResponseWriter, r *http.Request) {
    speedTestMutex.Lock()
    if isRunning {
        speedTestMutex.Unlock()
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{"error": "Speed test already running"})
        return
    }
    isRunning = true
    currentTest.Running = true
    currentTest.Error = ""
    speedTestMutex.Unlock()

    go runSpeedTest()

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "started"})
}

func runSpeedTest() {
    defer func() {
        speedTestMutex.Lock()
        isRunning = false
        currentTest.Running = false
        speedTestMutex.Unlock()
    }()

    cmd := exec.Command("speedtest-cli", "--simple")
    output, err := cmd.Output()
    
    speedTestMutex.Lock()
    defer speedTestMutex.Unlock()

    if err != nil {
        currentTest.Error = "Failed to run speedtest-cli: " + err.Error()
        return
    }

    lines := strings.Split(string(output), "\n")
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if strings.HasPrefix(line, "Ping:") {
            parts := strings.Fields(line)
            if len(parts) >= 2 {
                pingStr := strings.TrimSpace(parts[1])
                if ping, err := strconv.ParseFloat(pingStr, 64); err == nil {
                    currentTest.Ping = ping
                }
            }
        } else if strings.HasPrefix(line, "Download:") {
            parts := strings.Fields(line)
            if len(parts) >= 2 {
                downloadStr := strings.TrimSpace(parts[1])
                if download, err := strconv.ParseFloat(downloadStr, 64); err == nil {
                    currentTest.Download = download
                }
            }
        } else if strings.HasPrefix(line, "Upload:") {
            parts := strings.Fields(line)
            if len(parts) >= 2 {
                uploadStr := strings.TrimSpace(parts[1])
                if upload, err := strconv.ParseFloat(uploadStr, 64); err == nil {
                    currentTest.Upload = upload
                }
            }
        }
    }

    serverCmd := exec.Command("speedtest-cli", "--list")
    if serverOutput, err := serverCmd.Output(); err == nil {
        serverLines := strings.Split(string(serverOutput), "\n")
        for _, serverLine := range serverLines {
            if strings.Contains(serverLine, ")") && len(serverLine) > 10 {
                currentTest.Server = strings.TrimSpace(serverLine)
                break
            }
        }
    }

    currentTest.LastUpdated = time.Now().Format("2006-01-02 15:04:05")
}