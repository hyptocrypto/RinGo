package buffer

import (
	"fmt"
	"testing"

	"github.com/hyptocrypto/RinGo/config"
)

var conf *config.Config

func init() {
	conf = config.LoadConfig("../config.yaml")
}

func TestBufferOverflowLimit(t *testing.T) {
	buf := MakeFullBuffer()
	// Try to overflow buffer
	data := []byte("Test data 10000")
	err := buf.Write(data)
	if err == nil {
		t.Error("buffer overflow allowed")
	}
}

func TestBufferRead(t *testing.T) {
	buf := MakeFullBuffer()
	data := buf.Read()
	if len(data) != int(buf.readSize) {
		t.Errorf("data length discrepancy. expected: %v received: %v", buf.readSize, len(data))
	}
}

func MakeFullBuffer() *Buffer {
	// Fill buffer
	buff := NewBuffer(conf.BufferSize, conf.ReadInterval, conf.ReadSize, conf.OverwriteBuffer)
	for i := 0; i < int(conf.BufferSize); i++ {
		data := []byte(fmt.Sprintf("Test data %d", i))
		err := buff.Write(data)
		if err != nil {
			panic("error writing to buffer")
		}
	}
	return buff
}
