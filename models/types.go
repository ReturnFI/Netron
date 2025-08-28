package models

type SystemInfo struct {
    CPU       CPUInfo       `json:"cpu"`
    Memory    MemoryInfo    `json:"memory"`
    Processes []ProcessInfo `json:"processes"`
    Network   NetworkInfo   `json:"network"`
    SpeedTest SpeedTestInfo `json:"speedtest"`
    System    SystemDetails `json:"system"`
}

type CPUInfo struct {
    Model     string  `json:"model"`
    Cores     int     `json:"cores"`
    Frequency string  `json:"frequency"`
    Cache     string  `json:"cache"`
    Usage     float64 `json:"usage"`
    AES       bool    `json:"aes"`
    VMX       bool    `json:"vmx"`
}

type SystemDetails struct {
    OS             string `json:"os"`
    Kernel         string `json:"kernel"`
    Arch           string `json:"arch"`
    Uptime         string `json:"uptime"`
    LoadAverage    string `json:"load_average"`
    TCPCongestion  string `json:"tcp_cc"`
    Virtualization string `json:"virtualization"`
    IPv4Status     string `json:"ipv4_status"`
    IPv6Status     string `json:"ipv6_status"`
    Organization   string `json:"organization"`
    Location       string `json:"location"`
    Region         string `json:"region"`
    TotalDisk      string `json:"total_disk"`
    UsedDisk       string `json:"used_disk"`
}

type MemoryInfo struct {
    Total     uint64  `json:"total"`
    Used      uint64  `json:"used"`
    Available uint64  `json:"available"`
    Percent   float64 `json:"percent"`
}

type ProcessInfo struct {
    PID    int     `json:"pid"`
    Name   string  `json:"name"`
    CPU    float64 `json:"cpu"`
    Memory float64 `json:"memory"`
    Status string  `json:"status"`
}

type NetworkInfo struct {
    Interfaces []InterfaceInfo `json:"interfaces"`
    TCP        []Connection    `json:"tcp"`
    UDP        []Connection    `json:"udp"`
    TCPCount   int             `json:"tcp_count"`
    UDPCount   int             `json:"udp_count"`
}

type InterfaceInfo struct {
    Name      string `json:"name"`
    BytesSent uint64 `json:"bytes_sent"`
    BytesRecv uint64 `json:"bytes_recv"`
    Speed     uint64 `json:"speed"`
}

type Connection struct {
    LocalAddr  string `json:"local_addr"`
    RemoteAddr string `json:"remote_addr"`
    Status     string `json:"status"`
    PID        int    `json:"pid"`
}

type SpeedTestInfo struct {
    Running      bool    `json:"running"`
    Download     float64 `json:"download"`
    Upload       float64 `json:"upload"`
    Ping         float64 `json:"ping"`
    Server       string  `json:"server"`
    LastUpdated  string  `json:"last_updated"`
    Error        string  `json:"error,omitempty"`
}