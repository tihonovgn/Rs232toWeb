package main

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/jacobsa/go-serial/serial"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

var weights = make(map[string]string)
var config configuration

func main() {
	config = getConfig()
	for _ , unit := range config.Units {
		go rsHandler(unit)
	}

	http.HandleFunc("/", webHandler)
	http.ListenAndServe(":" + config.Port, nil)
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	result := make(map[string]string)

	unit, unitErr := getUnitFromPath(r.URL.Path, config.Units)
	if unitErr != nil {
		result["error"] = unitErr.Error();
	} else if weight, ok := weights[unit]; ok {
		result["weight"] = weight
	} else {
		result["weight"] = "0"
		result["error"] = "Can not get weight"
	}
	js, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func rsHandler(unit ScaleUnit) {
	port, err := serial.Open(unit.Options)
	if err != nil {
		fmt.Println("Can't open port: %v", err)
		return
	}
	defer port.Close()

	for {
		if len(unit.Cmd) > 0 {
			cmd := make([]byte, hex.DecodedLen(len(unit.Cmd)))
			hex.Decode(cmd, []byte(unit.Cmd))//todo обработать ошибку
			_, err := port.Write(cmd)
			if err != nil {
				log.Fatalf("port.Write: %v", err)
			}
			fmt.Println("Send ", unit.Cmd, " command.")
		}

		if unit.BinaryResponse {
			buf := make([]byte, 32)
			n, err := port.Read(buf)
			if err != nil {
				if err != io.EOF {
					fmt.Println("Error reading from serial port: ", err)
				}
				time.Sleep(1 * time.Second)
				continue
			}
			buf = buf[:n]
			if len(buf) < 5 {
				time.Sleep(1 * time.Second)
				continue
			}
			weight, _ := strconv.ParseFloat(fmt.Sprintf("%d%d", buf[3], buf[2]), 16)
			weight = weight * 0.1

			weights[unit.Name] = fmt.Sprintf("%.2f", weight)
		} else {
			reader := bufio.NewReader(port)
			buf, _, err := reader.ReadLine()
			if err != nil || len(buf) == 0 {
				time.Sleep(1 * time.Second)
				continue
			}

			weights[unit.Name] = parseWeight(buf, unit.Pattern)
		}
		fmt.Println("Rx ", unit.Name, " ", weights[unit.Name])
		time.Sleep(1 * time.Second)
	}
}
