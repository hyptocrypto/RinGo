package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hyptocrypto/RinGo/buffer"
)

type RequestBody struct {
	Data     []byte `json:"data"`
	DeviceId string `json:"deviceId"`
}

var buff *buffer.Buffer

func init() {
	buff = buffer.NewBuffer(10) // Initialize buffer with size 10
}

func handleBuff(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var body RequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	fmt.Println(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	fmt.Println(body)
	deviceId, err := uuid.Parse(body.DeviceId)
	if err != nil {
		http.Error(w, "Invalid deviceId format", http.StatusBadRequest)
		return
	}

	if err := buff.Write(deviceId, body.Data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Data stored successfully")
}

func main() {
	http.HandleFunc("/buff", handleBuff)
	fmt.Printf("RinGo server started \n")
	if err := http.ListenAndServe(":5555", nil); err != nil {
		panic(err)
	}
}
