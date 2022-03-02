package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/hyptocrypto/RinGo/buffer"
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

	deviceId, err := uuid.Parse(body.DeviceId)
	if err != nil {
		msg := "Invalid deviceId format"
		http.Error(w, msg, http.StatusBadRequest)
		log.Println(msg)
		return
	}

	if err := buff.Write(deviceId, body.Data); err != nil {
		msg := err.Error()
		http.Error(w, msg, http.StatusInternalServerError)
		log.Println(msg)
		return
	}

	msg := fmt.Sprintf("Client(%v) data buffered", deviceId)
	fmt.Fprint(w, msg)
	log.Println(msg)
}

func main() {
	http.HandleFunc("/buff", handleBuff)
	fmt.Printf("RinGo server started \n")
	go buffer.BufferReader(buff, dataChannel)
	go buffer.BufferChanConsumer(dataChannel)
	if err := http.ListenAndServe(conf.Port, nil); err != nil {
		panic(err)
	}
}
