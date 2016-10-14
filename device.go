package ener314

import (
	"fmt"
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

	fmt.Println("Resetting...")
	Reset()

	d.hrf, err = NewHRF()
	if err != nil {
		return err
	}

	fmt.Println(d.hrf.GetVersion())

	fmt.Println("Configuring FSK")
	err = d.hrf.ConfigFSK()
	if err != nil {
		return err
	}

	fmt.Println("Wait for ready...")
	d.hrf.WaitFor(ADDR_IRQFLAGS1, MASK_MODEREADY, true)

	fmt.Println("Clearing FIFO...")
	d.hrf.ClearFifo()

	for {
		msg := d.hrf.ReceiveFSKMessage()
		if msg == nil {
			time.Sleep(100 * time.Millisecond)
		} else {
			fmt.Println("Message:", msg)
		}
	}
	return nil
}
