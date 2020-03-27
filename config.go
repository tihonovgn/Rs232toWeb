package main

import (
	"encoding/json"
	"github.com/jacobsa/go-serial/serial"
	"log"
	"os"
)

type ScaleUnit struct {
	Name string
	Cmd string
	Pattern string
	Type string
	Options serial.OpenOptions
}

type configuration struct {
	Port string
	Units []ScaleUnit
}

func getConfig() configuration {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	conf := configuration{}
	err := decoder.Decode(&conf)
	if err != nil {
		log.Fatalf("Config read error: %v", err)
	}
	return conf
}
