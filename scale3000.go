package main

import (
	"bufio"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"regexp"
)

type Scale3000 struct {
}

func (rs Scale3000) SendCmd(port io.ReadWriteCloser, cmd string) {
}

func (rs Scale3000) ReadWeight(port io.ReadWriteCloser) (string, error) {
	reader := bufio.NewReader(port)
	buf, _, err := reader.ReadLine()
	if err != nil || len(buf) == 0 {
		return "", errors.New("Errort null weight")
	}
	return parseWeight(buf, "7777(.*?)6b67"), nil
}

func parseWeight(buf []byte, pattern string) string {
	str := hex.EncodeToString(buf)
	re := regexp.MustCompile(pattern)
	found := re.FindStringSubmatch(str)
	src := []byte(found[1])
	dst := make([]byte, hex.DecodedLen(len(src)))
	n, err := hex.Decode(dst, src)
	if err != nil {
		log.Fatal(err)
	}
	return string(dst[:n])
}