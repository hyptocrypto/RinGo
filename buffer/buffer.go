package buffer

import (
	"fmt"
	"sync"
	"time"
)

type ReadChunk = [][]byte

type Buffer struct {
	lock         sync.Mutex
	read         uint16
	write        uint16
	count        uint16
	size         uint16
	readSize     uint16
	readInterval uint16
	overwrite    bool
	buffer       [][]byte
}

func NewBuffer(size uint16, readInterval uint16, readSize uint16, overwriteBuffer bool) *Buffer {
	return &Buffer{
		read:         0,
		write:        0,
		count:        0,
		size:         size,
		readSize:     readSize,
		readInterval: readInterval,
		overwrite:    overwriteBuffer,
		buffer:       make([][]byte, size),
	}
}

// Write adds an item to the buffer if there is space.
func (b *Buffer) Write(data []byte) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.count == b.size {
		// buffer is full, cannot write
		return fmt.Errorf("Buffer full")
	}

	b.buffer[b.write] = data
	b.write = (b.write + 1) % b.size
	b.count++
	return nil
}

// Read returns up to 15 items from the buffer.
func (r *Buffer) Read() ReadChunk {
	r.lock.Lock()
	defer r.lock.Unlock()

	var items ReadChunk
	itemsToRead := min(r.count, r.readSize)
	for i := 0; uint16(i) < itemsToRead; i++ {
		items = append(items, r.buffer[r.read])
		r.read = (r.read + 1) % r.size
	}
	r.count -= itemsToRead
	return items
}

// Read from channel at set interval
func BufferReader(b *Buffer, out chan ReadChunk) {
	for {
		time.Sleep(time.Duration(b.readInterval) * time.Second)
		data := b.Read()
		fmt.Printf("Buffer size: %v\n", b.size)
		fmt.Printf("Buffer read: %v\n", b.read)
		fmt.Printf("Buffer writeIdx: %v\n", b.write)
		fmt.Printf("Data left to read in buffer: %v\n\n", b.count)
		out <- data
	}
}

// Consume/process data that has been pulled from buffer
func BufferConsumer(bufferChannel chan ReadChunk) {
	for {
		select {
		case d, ok := <-bufferChannel:
			if !ok {
				fmt.Println("Buffer channel closed. Exiting...")
				return
			}
			fmt.Printf("Pulled: %v chunks from buffer\n", len(d))
		default:
		}
		time.Sleep(100 * time.Millisecond) // Sleep to prevent busy-waiting
	}

}
