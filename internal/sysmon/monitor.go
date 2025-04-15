package sysmon

import (
	"context"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Monitor struct {
	conf   Config
	status Status
	mux    sync.Mutex

	ctx       context.Context
	ctxCancel context.CancelFunc
}

// NewMonitor instantiates a new Monitor instance
func NewMonitor(conf Config) *Monitor {
	ctx, cancel := context.WithCancel(context.Background())
	m := &Monitor{
		conf:      conf,
		ctx:       ctx,
		ctxCancel: cancel,
		status:    newStatus(),
	}

	m.run()

	return m
}

// Read returns the current stats of the machine
func (m *Monitor) Read() Status {
	m.mux.Lock()
	defer m.mux.Unlock()

	return m.status
}

// Close the Monitor's dependencies
func (m *Monitor) Close() {
	m.ctxCancel()
}

func (m *Monitor) run() error {
	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(m.conf.PollRate))
		defer ticker.Stop()
		defer func() {
			if m.mux.TryLock() {
				m.mux.Unlock()
			}
		}()

		for {
			select {
			case <-ticker.C:
				m.mux.Lock()

				m.status = Status{
					Uptime:  uptime(),
					Cpu:     cpu(m.status.Cpu),
					Memory:  memory(m.status.Memory),
					Mounts:  mounts(m.conf.Mounts, m.status.Mounts),
					Network: network(uint64(m.conf.PollRate), m.conf.NetInterfaces, m.status.Network),
				}

				m.mux.Unlock()
			case <-m.ctx.Done():
				return
			}
		}
	}()

	return nil
}

// uptime parses the /proc/uptime file and returns uptime in seconds
func uptime() uint64 {
	line, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0
	}

	parts := strings.Split(string(line), ".")
	uptime, _ := strconv.ParseUint(parts[0], 10, 64)

	return uptime
}

// cpu parses the /proc/cpuinfo file and returns usage staticstics for each core
func cpu(history map[string]CpuCore) map[string]CpuCore {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return history
	}

	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 || "cpu" == fields[0] || !strings.HasPrefix(fields[0], "cpu") {
			continue
		}

		key := fields[0]
		if _, ok := history[key]; !ok {
			history[key] = CpuCore{}
		}

		var total, idle uint64
		for i, field := range fields[1:] {
			val, _ := strconv.ParseUint(field, 10, 64)

			total += val

			if i == 3 {
				idle = val
			}
		}

		if core, ok := history[key]; !ok {
			history[key] = CpuCore{}
		} else {
			core.Total = total - history[key].Total
			core.Idle = idle - history[key].Idle
			history[key] = core
		}
	}

	return history
}

// memory parses the /proc/meminfo file and returns usage staticstics in kb
func memory(history Usage) Usage {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return history
	}

	var free uint64
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "MemTotal") {
			fields := strings.Fields(line)
			history.Total, _ = strconv.ParseUint(fields[1], 10, 64)
			history.Total *= 1024
		} else if strings.HasPrefix(line, "MemAvailable") {
			fields := strings.Fields(line)
			free, _ = strconv.ParseUint(fields[1], 10, 64)
			free *= 1024
		}
	}

	history.Used = history.Total - free

	return history
}

// mounts scans drive mounts based on the provided tracked slice
// and returns statistics on their utilization
func mounts(tracked []string, history map[string]Usage) map[string]Usage {
	if len(tracked) == 0 {
		return nil
	}

	for key := range history {
		if !slices.Contains(tracked, key) {
			delete(history, key)
		}
	}

	for _, mount := range tracked {
		entry, ok := history[mount]
		if !ok {
			entry = Usage{}
		}

		fs := syscall.Statfs_t{}
		if err := syscall.Statfs(mount, &fs); err != nil {
			continue
		}

		entry.Total = fs.Blocks * uint64(fs.Bsize)
		entry.Total = entry.Total - fs.Bfree*uint64(fs.Bsize)
		history[mount] = entry
	}

	return history
}

// network scans the /proc/net/dev file and returns traffic statistics on tracked interfaces in bytes
func network(interval uint64, tracked []string, history map[string]NetInterface) map[string]NetInterface {
	if len(tracked) == 0 {
		return nil
	}

	for key := range history {
		if !slices.Contains(tracked, key) {
			delete(history, key)
		}
	}

	data, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		return history
	}

	for _, line := range strings.Split(string(data), "\n") {
		if !strings.Contains(line, ":") {
			continue
		}

		parts := strings.Split(line, ":")
		iface := strings.Trim(parts[0], " ")
		if !slices.Contains(tracked, iface) {
			continue
		}

		entry, ok := history[iface]
		if !ok {
			entry = NetInterface{}
		}

		fields := strings.Fields(strings.Trim(parts[1], " "))

		rx, _ := strconv.ParseUint(fields[0], 10, 64)
		tx, _ := strconv.ParseUint(fields[8], 10, 64)

		entry.Rx = (rx - entry.Rx) / uint64(interval)
		entry.Tx = (tx - entry.Tx) / uint64(interval)

		history[iface] = entry
	}

	return history
}
