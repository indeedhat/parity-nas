package logging

import (
	"io"
	"sync"
)

var bufferMux sync.Mutex

type LogBuffer struct {
	// data held by the buffer
	data [][]byte
	// cursor is the current position in the buffer (next cell to write to)
	cursor int
	// cap is the capacity of the buffer total
	cap int
	// len is the current entry count in the buffer
	len int

	channels map[chan []byte]bool
}

func newBuffer(size int) *LogBuffer {
	return &LogBuffer{
		data:     make([][]byte, size),
		cap:      size,
		channels: make(map[chan []byte]bool),
	}
}

func (b *LogBuffer) Connect(ch chan []byte, preload int) int {
	b.channels[ch] = true

	if preload == 0 {
		return 0
	}

	var logs [][]byte
	if preload < 0 {
		logs, _ = b.ReadAll()
	} else {
		logs, _ = b.ReadN(preload)
	}

	go func() {
		for _, l := range logs {
			ch <- l
		}
	}()

	return len(logs)
}

func (b *LogBuffer) Disconnect(ch chan []byte) {
	if _, ok := b.channels[ch]; ok {
		delete(b.channels, ch)
	}
}

func (b *LogBuffer) Write(data []byte) (int, error) {
	bufferMux.Lock()
	defer bufferMux.Unlock()

	if b.len < b.cap {
		b.cap++
	}

	b.data[b.cursor] = data
	for ch, _ := range b.channels {
		ch <- data
	}

	b.cursor++
	if b.cursor >= b.cap {
		b.cursor = 0
	}

	return len(data), nil
}

func (b *LogBuffer) ReadAll() ([][]byte, int) {
	if b.len < b.cap {
		return b.data[:b.cursor], b.len
	}

	return append(b.data[b.cursor+1:b.cap-1], b.data[:b.cursor]...), b.len
}

func (b *LogBuffer) ReadN(n int) ([][]byte, int) {
	if b.cursor >= n {
		return b.data[b.cursor-n : b.cursor], n
	}

	if b.len < b.cap {
		if n > b.len {
			return b.data[:b.len], b.len
		} else {
			return b.data[:n], b.len
		}
	}

	if n > b.cap {
		n = b.cap
	}

	return append(b.data[b.cursor+1:b.cap-1], b.data[:b.cursor]...)[:n], n
}

var _ io.Writer = (*LogBuffer)(nil)
