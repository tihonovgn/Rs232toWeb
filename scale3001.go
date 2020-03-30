package main

import (
	"errors"
	"fmt"
	"io"
)

type Scale3001 struct {
}

func (rs Scale3001) SendCmd(port io.ReadWriteCloser, cmd string) {
	fmt.Println("method SendCmd")
}

func (rs Scale3001) ReadWeight(port io.ReadWriteCloser) (string, error) {
	fmt.Println("method ReadWeight")
	startByte := make([]byte, 1)
	for {
		port.Read(startByte)
		if string(startByte) == "=" {
			break
		}
	}
	buf := make([]byte, 7)
	n, err := port.Read(buf)
	if err != nil {
		if err != io.EOF {
			fmt.Println("Error reading from serial port: ", err)
		}
		return "", errors.New("Read error")
	}
	buf = buf[0:n]

	return string(rs.reverse(buf)), nil
}

func (rs Scale3001) reverse(weight []byte) []byte {
	for i, j := 0, len(weight)-1; i < j; i, j = i+1, j-1 {
		weight[i], weight[j] = weight[j], weight[i]
	}
	return weight
}
