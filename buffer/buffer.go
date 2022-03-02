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

// Len of data left in buffer
func (b *Buffer) Len() uint16 {
	if b.writeIdx >= b.readIdx {
		return b.writeIdx - b.readIdx
	}
	return b.size - (b.readIdx - b.writeIdx)
}

func (b *Buffer) IncWriteIdx() (uint16, error) {
	currIdx := b.writeIdx
	nextIdx := (b.writeIdx + 1) % b.size
	if currIdx == b.readIdx && len(b.buff[b.readIdx]) != 0 {
		// If read and write index are equal and there is data at that index, then the buffer is full
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

// TODO This is ugly and feels hacky. Find a better way to do this
// When we read a chunk of data from the buffer, we replace it with null bytes.
func (b *Buffer) Read() (ReadChunk, error) {
	b.mut.Lock()
	defer b.mut.Unlock()
	var ret ReadChunk // Initialize ret as nil slice

	if b.readIdx == b.writeIdx {
		if len(b.buff[b.readIdx]) > 0 {
			// Full buffer that has data
			ret = append(ret, b.buff[b.readIdx:b.readIdx+b.readSize]...)
			copy(b.buff[b.readIdx:b.readIdx+b.readSize], make([][]byte, b.readSize))
			b.readIdx += b.readSize
			return ret, nil
		}
		// Empty buffer
		return ret, nil
	}

	if b.readIdx > b.writeIdx {
		// Read to end of buffer, and wrap around if needed
		if (b.readIdx + b.readSize) >= b.size {
			diff := (b.readIdx + b.readSize) - b.size
			ret = append(ret, b.buff[b.readIdx:]...)
			copy(b.buff[b.readIdx:], make([][]byte, len(b.buff[b.readIdx:])))
			ret = append(ret, b.buff[:diff]...)
			copy(b.buff[:diff], make([][]byte, len(b.buff[:diff])))
			b.readIdx = diff
		} else {
			ret = append(ret, b.buff[b.readIdx:b.readIdx+b.readSize]...)
			copy(b.buff[b.readIdx:b.readIdx+b.readSize], make([][]byte, b.readSize))
			b.readIdx += b.readSize
		}
		return ret, nil
	}
	// Catch up to writeIdx
	if b.readIdx < b.writeIdx {
		if b.readIdx+b.readSize >= b.writeIdx {
			ret = append(ret, b.buff[b.readIdx:b.writeIdx]...)
			copy(b.buff[b.readIdx:b.writeIdx], make([][]byte, len(b.buff[b.readIdx:b.writeIdx])))
			b.readIdx = b.writeIdx
		} else {
			ret = append(ret, b.buff[b.readIdx:b.readIdx+b.readSize]...)
			copy(b.buff[b.readIdx:b.readIdx+b.readSize], make([][]byte, len(b.buff[b.readIdx:b.readIdx+b.readSize])))
			b.readIdx += b.readSize
		}
		return ret, nil
	}
	return ret, errors.New("un-handled buffer read case")
}

func BufferReader(b *Buffer, out chan ReadChunk) {
	for {
		time.Sleep(time.Duration(b.readInterval) * time.Second)
		data, err := b.Read()
		// fmt.Print(data)
		fmt.Printf("Buffer size: %v\n", b.size)
		fmt.Printf("Buffer readIdx: %v\n", b.readIdx)
		fmt.Printf("Buffer writeIdx: %v\n", b.writeIdx)
		fmt.Printf("Data left to read in buffer: %v\n\n", b.Len())
		// fmt.Print(b.buff)
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
			fmt.Printf("Pulled: %v chunks from buffer\n", len(d))
		default:
		}
		time.Sleep(100 * time.Millisecond) // Sleep to prevent busy-waiting
	}

}
