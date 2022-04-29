package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/mmcloughlin/geohash"
)

func main() {
	if len(os.Args) <= 2 {
		log.Printf("Usage: %s <lng> <lat>", os.Args[0])
		os.Exit(1)
	}

	lng, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		panic(err)
	}
	lat, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		panic(err)
	}

	hash := geohash.Encode(lat, lng)
	fmt.Print(hash)
}
