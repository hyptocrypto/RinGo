package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hyptocrypto/RinGo/buffer"
	"github.com/hyptocrypto/RinGo/client"
	"github.com/hyptocrypto/RinGo/config"
)

type RequestBody struct {
	Data     []byte `json:"data"`
	DeviceId string `json:"deviceId"`
}

var buff *buffer.Buffer
var conf *config.Config
var dataChannel chan buffer.ReadChunk

func init() {
	conf = config.LoadConfig("config.yaml")
	buff = buffer.NewBuffer(conf.BufferSize, conf.ReadInterval, conf.ReadSize, conf.OverwriteBuffer)
	dataChannel = make(chan buffer.ReadChunk)
}

func handleBuff(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		msg := "Invalid request method"
		http.Error(w, msg, http.StatusMethodNotAllowed)
		log.Println(msg)
		return
	}

	var body RequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		msg := "Error reading request body"
		http.Error(w, msg, http.StatusInternalServerError)
		log.Println(msg)
		return
	}
	client := client.ClientFromRequest(r)
	if len(body.Data) == 0 {
		msg := "Body must contain data"
		http.Error(w, msg, http.StatusBadRequest)
		log.Println(msg)
	}

	// If all checks pass, write to the buffer
	if err := buff.Write(body.Data); err != nil {
		msg := err.Error()
		http.Error(w, msg, http.StatusInternalServerError)
		log.Println(msg)
		return
	}

	msg := fmt.Sprintf("Client(%v) data buffered", client.ID)
	fmt.Fprint(w, msg)
	log.Println(msg)
}

func main() {
	http.HandleFunc("/buff", handleBuff)
	fmt.Printf("RinGo server started \n")
	go buffer.BufferReader(buff, dataChannel)
	go buffer.BufferConsumer(dataChannel)
	if err := http.ListenAndServe(conf.Port, nil); err != nil {
		panic(err)
	}
}
