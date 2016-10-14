package ener314

import (
	"fmt"
	"log"
	"time"
)

type Device struct {
	hrf *HRF
}

func NewDevice() *Device {
	return &Device{}
}

func (d *Device) Start() error {
	var err error

	log.Println("Resetting...")
	Reset()

	d.hrf, err = NewHRF()
	if err != nil {
		return err
	}

	version := d.hrf.GetVersion()
	if version != 36 {
		return fmt.Errorf("Unexpected version: %d", version)
	}

	log.Println("Configuring FSK")
	err = d.hrf.ConfigFSK()
	if err != nil {
		return err
	}

	log.Println("Wait for ready...")
	d.hrf.WaitFor(ADDR_IRQFLAGS1, MASK_MODEREADY, true)

	log.Println("Clearing FIFO...")
	d.hrf.ClearFifo()

	for {
		msg := d.hrf.ReceiveFSKMessage()
		if msg == nil {
			time.Sleep(100 * time.Millisecond)
		} else {
			log.Println("Message:", msg)
		}
	}
	return nil
}
