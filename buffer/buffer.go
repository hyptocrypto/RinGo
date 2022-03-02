package buffer

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type ReadChunk = [][]byte

// An IOT devices sending in data
type Client struct {
	id         uuid.UUID
	clientAddr string
}

type Buffer struct {
	mut          sync.Mutex
	readIdx      uint16
	writeIdx     uint16
	size         uint16
	readSize     uint16
	readInterval uint16
	overwrite    bool
	buff         [][]byte
}

func NewBuffer(size uint16, readInterval uint16, readSize uint16, overwriteBuffer bool) *Buffer {
	return &Buffer{
		readIdx:      0,
		writeIdx:     0,
		size:         size,
		readSize:     readSize,
		readInterval: readInterval,
		overwrite:    overwriteBuffer,
		buff:         make([][]byte, size),
	}
}

func (b *Buffer) IncWriteIdx() (uint16, error) {
	currIdx := b.writeIdx
	var nextIdx uint16
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
	ret := make(ReadChunk, b.readSize)
	if b.readIdx+b.readSize <= b.writeIdx {
		copy(ret, b.buff[b.readIdx:b.readIdx+b.readSize])
		b.readIdx += b.readSize
		return ret, nil
	} else {
		copy(ret, b.buff[b.readIdx:b.writeIdx])
		b.readIdx = b.writeIdx
		return ret, nil
	}

}

func BufferReader(b *Buffer, out chan ReadChunk) {
	for {
		time.Sleep(time.Duration(b.readInterval) * time.Second)
		data, err := b.Read()
		// fmt.Print(data)
		// fmt.Printf("Buffer size: %v\n", b.size)
		// fmt.Printf("Buffer readIdx: %v\n", b.readIdx)
		// fmt.Printf("Buffer writeIdx: %v\n", b.writeIdx)
		if err != nil {
			log.Fatal("Error reading from buffer!")
			close(out)
			panic(err)
		}
		out <- data
	}
}

func BufferChanConsumer(bufferChannel chan ReadChunk) {
	for {
		select {
		case d, ok := <-bufferChannel:
			if !ok {
				fmt.Println("Buffer channel closed. Exiting...")
				return
			}
			fmt.Println("Received:", len(d))
		default:
		}
		time.Sleep(100 * time.Millisecond) // Sleep to prevent busy-waiting
	}

}
