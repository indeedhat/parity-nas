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
}

func newBuffer(size int) *LogBuffer {
	return &LogBuffer{
		data: make([][]byte, size),
		cap:  size,
	}
}

func (b *LogBuffer) Write(data []byte) (int, error) {
	bufferMux.Lock()
	defer bufferMux.Unlock()

	if b.len < b.cap {
		b.cap++
	}

	b.data[b.cursor] = data

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

var _ io.Writer = (*LogBuffer)(nil)
