package main

import (
	"encoding/hex"
	"errors"
	"log"
	"regexp"
	"strings"
)

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
