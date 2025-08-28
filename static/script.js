class Dashboard {
    constructor() {
        this.init();
        this.startUpdates();
        this.initSpeedTest();
    }

    init() {
        this.updateData();
    }

    initSpeedTest() {
        const btn = document.getElementById('speedtest-btn');
        btn.addEventListener('click', () => this.startSpeedTest());
    }

    async updateData() {
        try {
            const response = await fetch('/api/system');
            const data = await response.json();
            this.updateUI(data);
        } catch (error) {
            console.error('Failed to fetch system data:', error);
        }
    }

    updateUI(data) {
        this.updateCPU(data.cpu);
        this.updateMemory(data.memory);
        this.updateProcesses(data.processes);
        this.updateNetwork(data.network);
        this.updateSpeedTest(data.speedtest);
        this.updateSystemInfo(data.system, data.cpu);
    }

    updateCPU(cpu) {
        document.getElementById('cpu-usage').textContent = `${cpu.usage.toFixed(1)}%`;
    }

    updateMemory(memory) {
        document.getElementById('mem-percent').textContent = `${memory.percent.toFixed(1)}%`;
        document.getElementById('mem-used').textContent = this.formatBytes(memory.used);
        document.getElementById('mem-total').textContent = this.formatBytes(memory.total);
    }

    updateProcesses(processes) {
        const tbody = document.getElementById('processes-table');
        tbody.innerHTML = '';

        processes.forEach(proc => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${proc.pid}</td>
                <td>${proc.name}</td>
                <td>${proc.cpu.toFixed(1)}%</td>
                <td>${this.formatBytes(proc.memory)}</td>
                <td>${proc.status}</td>
            `;
            tbody.appendChild(row);
        });
    }

    updateNetwork(network) {
        this.updateInterfaces(network.interfaces);
        this.updateConnections(network.tcp, 'tcp-table');
        this.updateConnections(network.udp, 'udp-table');
        
        document.getElementById('tcp-count').textContent = network.tcp_count || 0;
        document.getElementById('udp-count').textContent = network.udp_count || 0;
    }

    updateInterfaces(interfaces) {
        const tbody = document.getElementById('interfaces-table');
        tbody.innerHTML = '';

        interfaces.forEach(iface => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${iface.name}</td>
                <td>${this.formatBytes(iface.bytes_sent)}</td>
                <td>${this.formatBytes(iface.bytes_recv)}</td>
                <td>${this.formatSpeed(iface.speed)}</td>
            `;
            tbody.appendChild(row);
        });
    }

    updateConnections(connections, tableId) {
        const tbody = document.getElementById(tableId);
        tbody.innerHTML = '';

        connections.slice(0, 10).forEach(conn => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${conn.local_addr}</td>
                <td>${conn.remote_addr}</td>
                <td>${conn.status}</td>
                <td>${conn.pid || '-'}</td>
            `;
            tbody.appendChild(row);
        });
    }

    formatBytes(bytes) {
        if (bytes === 0) return '0 B';
        const k = 1024;
        const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    }

    formatSpeed(speed) {
        if (speed >= 1000000000) {
            return (speed / 1000000000).toFixed(1) + ' Gbps';
        }
        if (speed >= 1000000) {
            return (speed / 1000000).toFixed(1) + ' Mbps';
        }
        return (speed / 1000).toFixed(1) + ' Kbps';
    }

    startUpdates() {
        setInterval(() => {
            this.updateData();
        }, 3000);
    }

    async startSpeedTest() {
        const btn = document.getElementById('speedtest-btn');
        btn.disabled = true;
        btn.textContent = 'Running...';

        try {
            await fetch('/api/speedtest/start', { method: 'POST' });
        } catch (error) {
            console.error('Failed to start speed test:', error);
            btn.disabled = false;
            btn.textContent = 'Start Speed Test';
        }
    }

    updateSpeedTest(speedtest) {
        const btn = document.getElementById('speedtest-btn');
        
        if (speedtest.running) {
            btn.disabled = true;
            btn.textContent = 'Running...';
        } else {
            btn.disabled = false;
            btn.textContent = 'Start Speed Test';
        }

        if (speedtest.error) {
            document.getElementById('download-speed').textContent = 'Error';
            document.getElementById('upload-speed').textContent = 'Error';
            document.getElementById('ping-time').textContent = 'Error';
            document.getElementById('server-info').textContent = speedtest.error;
            document.getElementById('last-updated').textContent = '';
            return;
        }

        document.getElementById('download-speed').textContent = 
            speedtest.download ? speedtest.download.toFixed(2) : '-';
        document.getElementById('upload-speed').textContent = 
            speedtest.upload ? speedtest.upload.toFixed(2) : '-';
        document.getElementById('ping-time').textContent = 
            speedtest.ping ? speedtest.ping.toFixed(1) : '-';
        document.getElementById('server-info').textContent = 
            speedtest.server || '';
        document.getElementById('last-updated').textContent = 
            speedtest.last_updated ? `Last updated: ${speedtest.last_updated}` : '';
    }

    updateSystemInfo(system, cpu) {
        document.getElementById('cpu-model').textContent = cpu.model || 'Unknown';
        document.getElementById('cpu-cores-detailed').textContent = 
            cpu.frequency ? `${cpu.cores} @ ${cpu.frequency}` : `${cpu.cores}`;
        document.getElementById('cpu-cache').textContent = cpu.cache || 'Unknown';
        document.getElementById('cpu-aes').textContent = cpu.aes ? '✓ Enabled' : '✗ Disabled';
        document.getElementById('cpu-vmx').textContent = cpu.vmx ? '✓ Enabled' : '✗ Disabled';
        
        document.getElementById('total-disk').textContent = 
            `${system.total_disk} (${system.used_disk} Used)`;
        document.getElementById('os-info').textContent = system.os || 'Unknown';
        document.getElementById('kernel-info').textContent = system.kernel || 'Unknown';
        document.getElementById('arch-info').textContent = system.arch || 'Unknown';
        document.getElementById('uptime-info').textContent = system.uptime || 'Unknown';
        document.getElementById('load-avg').textContent = system.load_average || 'Unknown';
        document.getElementById('tcp-cc').textContent = system.tcp_cc || 'Unknown';
        document.getElementById('virt-info').textContent = system.virtualization || 'Unknown';
        document.getElementById('ip-status').textContent = 
            `${system.ipv4_status} / ${system.ipv6_status}`;
        document.getElementById('organization').textContent = system.organization || 'Unknown';
        document.getElementById('location').textContent = system.location || 'Unknown';
    }
}

document.addEventListener('DOMContentLoaded', () => {
    new Dashboard();
});