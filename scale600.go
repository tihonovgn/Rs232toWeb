package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
)

type Scale600 struct {
}

func (rs Scale600) SendCmd(port io.ReadWriteCloser, cmd string) {
	command := make([]byte, hex.DecodedLen(len(cmd)))
	hex.Decode(command, []byte(cmd))//todo обработать ошибку
	_, err := port.Write(command)
	if err != nil {
		log.Fatalf("port.Write: %v", err)
	}
}

func (rs Scale600) ReadWeight(port io.ReadWriteCloser) (string, error) {
	buf := make([]byte, 32)
	n, err := port.Read(buf)
	if err != nil {
		if err != io.EOF {
			fmt.Println("Error reading from serial port: ", err)
		}
		return "", errors.New("Read error")
	}
	buf = buf[:n]
	if len(buf) < 5 {
		return "", errors.New("Read error")
	}
	weight, _ := strconv.ParseFloat(fmt.Sprintf("%d%d", buf[3], buf[2]), 16)
	weight = weight * 0.1

	return fmt.Sprintf("%.2f", weight), nil
}
