package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jacobsa/go-serial/serial"
	"net/http"
	"strings"
	"time"
)

var weights = make(map[string]string)
var config configuration

func main() {
	config = getConfig()
	for _ , unit := range config.Units {
		device := createDevice(unit.Type)
		go rsHandler(unit, device)
	}

	http.HandleFunc("/", webHandler)
	http.ListenAndServe(":" + config.Port, nil)
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	result := make(map[string]string)
	result["weight"] = "0"

	unit, unitErr := getUnitFromPath(r.URL.Path, config.Units)
	if unitErr != nil {
		result["error"] = unitErr.Error();
	} else if weight, ok := weights[unit]; ok {
		result["weight"] = weight
	} else {
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

func rsHandler(unit ScaleUnit, device Device) {
	port, err := serial.Open(unit.Options)
	if err != nil {
		fmt.Println("Can't open port: %v", err)
		return
	}
	defer port.Close()

	for {
		device.SendCmd(port, unit.Cmd)
		weight, err := device.ReadWeight(port)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		weights[unit.Name] = weight
		fmt.Println("Rx ", unit.Name, " ", weights[unit.Name])
		time.Sleep(1 * time.Second)
	}
}

func getUnitFromPath(path string, units []ScaleUnit ) (string, error) {
	unitName := strings.Trim(path, "/")
	found := false
	for _ , unit := range units {
		if unit.Name == unitName {
			found = true
		}
	}
	if ! found {
		return "", errors.New("Unit not found " + unitName)
	}
	return unitName, nil
}
