package main

import (
	"fmt"

	"github.com/hyptocrypto/RinGo/buffer"
)

func main() {
	b := buffer.NewBuffer(100)
	fmt.Println(b)
}
