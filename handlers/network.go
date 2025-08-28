package handlers

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "os"
    "strconv"
    "strings"

    "netron/models"
)

func getNetworkInfo() models.NetworkInfo {
    tcp := getTCPConnections()
    udp := getUDPConnections()
    
    return models.NetworkInfo{
        Interfaces: getInterfaces(),
        TCP:        tcp,
        UDP:        udp,
        TCPCount:   len(tcp),
        UDPCount:   len(udp),
    }
}

func getInterfaces() []models.InterfaceInfo {
    data, err := ioutil.ReadFile("/proc/net/dev")
    if err != nil {
        return []models.InterfaceInfo{}
    }

    lines := strings.Split(string(data), "\n")
    var interfaces []models.InterfaceInfo

    for i, line := range lines {
        if i < 2 {
            continue
        }

        fields := strings.Fields(line)
        if len(fields) < 10 {
            continue
        }

        name := strings.TrimSuffix(fields[0], ":")
        if name == "lo" {
            continue
        }

        bytesRecv, _ := strconv.ParseUint(fields[1], 10, 64)
        bytesSent, _ := strconv.ParseUint(fields[9], 10, 64)

        speed := uint64(1000000000)
        speedFile := fmt.Sprintf("/sys/class/net/%s/speed", name)
        if speedData, err := ioutil.ReadFile(speedFile); err == nil {
            if s, err := strconv.ParseUint(strings.TrimSpace(string(speedData)), 10, 64); err == nil {
                speed = s * 1000000
            }
        }

        interfaces = append(interfaces, models.InterfaceInfo{
            Name:      name,
            BytesSent: bytesSent,
            BytesRecv: bytesRecv,
            Speed:     speed,
        })
    }

    return interfaces
}

func getTCPConnections() []models.Connection {
    return getConnections("/proc/net/tcp")
}

func getUDPConnections() []models.Connection {
    return getConnections("/proc/net/udp")
}

func getConnections(path string) []models.Connection {
    file, err := os.Open(path)
    if err != nil {
        return []models.Connection{}
    }
    defer file.Close()

    var connections []models.Connection
    scanner := bufio.NewScanner(file)
    scanner.Scan()

    for scanner.Scan() {
        line := scanner.Text()
        fields := strings.Fields(line)
        if len(fields) < 10 {
            continue
        }

        localAddr := parseAddr(fields[1])
        remoteAddr := parseAddr(fields[2])
        status := parseStatus(fields[3])

        pid := 0
        if len(fields) > 7 {
            pid, _ = strconv.Atoi(fields[7])
        }

        connections = append(connections, models.Connection{
            LocalAddr:  localAddr,
            RemoteAddr: remoteAddr,
            Status:     status,
            PID:        pid,
        })
    }

    return connections
}

func parseAddr(addr string) string {
    parts := strings.Split(addr, ":")
    if len(parts) != 2 {
        return addr
    }

    ipHex := parts[0]
    portHex := parts[1]

    if len(ipHex) == 8 {
        ip1, _ := strconv.ParseUint(ipHex[6:8], 16, 8)
        ip2, _ := strconv.ParseUint(ipHex[4:6], 16, 8)
        ip3, _ := strconv.ParseUint(ipHex[2:4], 16, 8)
        ip4, _ := strconv.ParseUint(ipHex[0:2], 16, 8)
        port, _ := strconv.ParseUint(portHex, 16, 16)
        return fmt.Sprintf("%d.%d.%d.%d:%d", ip1, ip2, ip3, ip4, port)
    }

    return addr
}

func parseStatus(status string) string {
    statusMap := map[string]string{
        "01": "ESTABLISHED",
        "02": "SYN_SENT",
        "03": "SYN_RECV",
        "04": "FIN_WAIT1",
        "05": "FIN_WAIT2",
        "06": "TIME_WAIT",
        "07": "CLOSE",
        "08": "CLOSE_WAIT",
        "09": "LAST_ACK",
        "0A": "LISTEN",
        "0B": "CLOSING",
    }

    if s, exists := statusMap[status]; exists {
        return s
    }
    return status
}