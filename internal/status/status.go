package status

import "encoding/json"

type Config struct {
	PollRate      uint8
	Mounts        []string
	NetInterfaces []string
}

type Status struct {
	Uptime  uint64                  `json:"uptime"`
	Cpu     map[string]CpuCore      `json:"cpu,omitempty"`
	Mounts  map[string]Usage        `json:"mounts,omitempty"`
	Memory  Usage                   `json:"memory,omitempty"`
	Network map[string]NetInterface `json:"network,omitempty"`
}

func newStatus() Status {
	return Status{
		Cpu:     make(map[string]CpuCore),
		Mounts:  make(map[string]Usage),
		Network: make(map[string]NetInterface),
	}
}

// Json marshal the message struct to json ready for transit
func (s Status) Json() (string, error) {
	data, err := json.Marshal(&s)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

type CpuCore struct {
	Total uint64 `json:"total"`
	Idle  uint64 `json:"idle"`
}

type Usage struct {
	Used  uint64 `json:"used"`
	Total uint64 `json:"total"`
}

type NetInterface struct {
	Rx uint64 `json:"rx"`
	Tx uint64 `json:"tx"`
}
