package ping

import (
	"fmt"
	"strconv"
	"strings"
)

func validateHostFlag() []byte {
	toReturn := make([]byte, 4)
	splitted := strings.Split(*flagip, ".")
	if len(splitted) < 4 {
		panic("Error! usage: -h a.b.c.d (example: -h 127.0.0.1)")
	}
	for i, e := range splitted {
		res, err := strconv.Atoi(e)
		if err != nil {
			panic(fmt.Sprintf("Error! usage: -h a.b.c.d (example: -h 127.0.0.1): %v", err))
		}
		temp := byte(res)
		if temp < 0 || temp > 255 {
			panic("Couldn't parse -h parameter")
		}

		toReturn[i] = temp
	}

	fmt.Printf("Sending ping to: %s\n---------------------------------\n", *flagip)
	return toReturn
}

func getOsFromTTL(TTL uint8) string {
	if TTL <= 64 && TTL > 0 {
		return "Linux"
	} else if TTL <= 128 {
		return "Windows"
	} else {
		return "Unknown"
	}
}

func csum(b []byte) uint16 {
	var s uint32
	for i := 0; i < len(b); i += 2 {
		s += uint32(b[i+1])<<8 | uint32(b[i])
	}
	// add back the carry
	s = s>>16 + s&0xffff
	s = s + s>>16
	return uint16(^s)
}
