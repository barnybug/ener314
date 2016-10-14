package main

import (
	"fmt"

	"github.com/barnybug/ener314"
)

func fatalIfErr(err error) {
	if err != nil {
		panic(fmt.Sprint("Error:", err))
	}
}

func main() {
	dev := ener314.NewDevice()
	dev.Start()
}
