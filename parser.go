package main

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
