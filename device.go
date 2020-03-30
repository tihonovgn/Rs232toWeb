package main

import (
	"fmt"
	"io"
	"log"
)

type Device interface {
	SendCmd(port io.ReadWriteCloser, cmd string)
	ReadWeight(port io.ReadWriteCloser) (string, error)
}

func createDevice(unitType string) Device{
	switch unitType {
	case "scale600":
		return Scale600{}
	case "scale3000":
		return Scale3000{}
	case "scale3001":
		fmt.Println("Device Scale3001")
		return Scale3001{}
	default:
		log.Fatalln("Device type error")
		return nil
	}
}