package buffer

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

var clientID = uuid.New()

func TestBufferLimit(t *testing.T) {
	buf := MakeFullBuffer()
	// Try to over flow buffer
	data := []byte("Test data 10000")
	err := buf.Write(clientID, data)
	if err == nil {
		t.Error("Buffer overflow allowed")
	}
}

func MakeFullBuffer() *Buffer {
	buf := NewBuffer(10)
	numIterations := 10
	// Fill buffer
	for i := 0; i < numIterations; i++ {
		data := []byte(fmt.Sprintf("Test data %d", i))
		err := buf.Write(clientID, data)
		fmt.Println(buf.writeIdx)
		if err != nil {
			panic("Error writing to buffer")
		}
	}
	return buf
}
