package main

import (
	"fmt"
	"time"

	"github.com/barnybug/ener314"
)

func fatalIfErr(err error) {
	if err != nil {
		panic(fmt.Sprint("Error:", err))
	}
}

const (
	/* OpenThings definitions */
	engManufacturerId = 0x04 // Energenie Manufacturer Id
	eTRVProductId     = 0x3  // Product ID for eTRV
	encryptId         = 0xf2 // Encryption ID for eTRV
)

func main() {
	hrf, err := ener314.NewHRF()
	fatalIfErr(err)

	fmt.Println(hrf.GetVersion())

	fmt.Println("Configuring FSK")
	err = hrf.ConfigFSK()
	fatalIfErr(err)

	fmt.Println("Wait for ready...")
	hrf.WaitFor(ener314.ADDR_IRQFLAGS1, ener314.MASK_MODEREADY, true)

	fmt.Println("Clearing FIFO...")
	hrf.ClearFifo()

	for {
		msg := hrf.ReceiveFSKMessage(encryptId, eTRVProductId, engManufacturerId)
		if msg == nil {
			time.Sleep(100 * time.Millisecond)
		} else {
			fmt.Printf("%+v\n", msg)
		}
	}
}
