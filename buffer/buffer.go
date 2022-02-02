package buffer

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

const (
	readSize = 2 // Number of chunks of data to read from buffer at a time
)

type ReadChunk = [readSize][]byte

// An IOT devices sending in data
type Client struct {
	id         uuid.UUID
	clientAddr string
}

type Buffer struct {
	mut      sync.Mutex
	readIdx  uint8
	writeIdx uint8
	size     uint8
	readSize uint8
	buff     [][]byte
}

func NewBuffer(size uint8) *Buffer {
	return &Buffer{
		readIdx:  0,
		writeIdx: 0,
		readSize: readSize,
		size:     size,
		buff:     make([][]byte, size),
	}
}

func (b *Buffer) IncWriteIdx() (uint8, error) {
	currIdx := b.writeIdx
	var nextIdx uint8
	if b.writeIdx+1 == b.size {
		nextIdx = 0
	} else {
		nextIdx = b.writeIdx + 1
	}
	if currIdx == b.readIdx && len(b.buff[b.readIdx]) != 0 {
		return 0, errors.New("Buffer full")
	}

	b.writeIdx = nextIdx
	return currIdx, nil
}

func (b *Buffer) Write(clientId uuid.UUID, data []byte) error {
	b.mut.Lock()
	defer b.mut.Unlock()
	writeIdx, err := b.IncWriteIdx()
	if err != nil {
		return err
	}
	b.buff[writeIdx] = data
	return nil
}

func (b *Buffer) Read() (ReadChunk, error) {
	b.mut.Lock()
	defer b.mut.Unlock()
	var ret ReadChunk

	if b.readIdx+b.readSize <= b.writeIdx {
		copy(ret[:], b.buff[b.readIdx:b.readIdx+b.readSize])
		b.readIdx += b.readSize
		return ret, nil
	} else {
		copy(ret[:], b.buff[b.readIdx:b.writeIdx])
		b.readIdx = b.writeIdx
		return ret, nil
	}

}
