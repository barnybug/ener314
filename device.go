package ener314

import (
	"fmt"
	"log"
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
	err = Reset()
	if err != nil {
		return err
	}

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
	return nil
}

func (d *Device) Receive() *Message {
	msg := d.hrf.ReceiveFSKMessage()
	if msg == nil {
		return nil
	}
	if msg.ManuId != energenieManuId {
		log.Printf("Warning: ignored message from manufacturer %d", msg.ManuId)
		return nil
	}
	if msg.ProdId != eTRVProdId {
		log.Printf("Warning: ignored message from product %d", msg.ProdId)
		return nil
	}
	if len(msg.Records) == 0 {
		log.Println("Warning: ignoring message with 0 records")
		return nil
	}
	return msg
}

func (d *Device) Respond(sensorId uint32, record Record) {
	message := &Message{
		ManuId:   energenieManuId,
		ProdId:   eTRVProdId,
		SensorId: sensorId,
		Records:  []Record{record},
	}
	err := d.hrf.SendFSKMessage(message)
	if err != nil {
		log.Println("Error sending", err)
	}
}

func (d *Device) Identify(sensorId uint32) {
	log.Printf("Requesting identify from device %06x", sensorId)
	d.Respond(sensorId, Identify{})
}

func (d *Device) Join(sensorId uint32) {
	log.Printf("Responding to Join from device %06x", sensorId)
	d.Respond(sensorId, JoinReport{})
}

func (d *Device) Voltage(sensorId uint32) {
	log.Printf("Requesting Voltage from device %06x", sensorId)
	d.Respond(sensorId, Voltage{})
}
