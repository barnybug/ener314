package ener314

import (
	"encoding/hex"
	"fmt"

	"github.com/quinte17/spi"
)

type HRF struct {
	dev *spi.SPI
}

const (
	SEED_PID               = 0x01
	MANUF_SENTEC           = 0x01
	PRODUCT_SENTEC_DEFAULT = 0x01
	MESSAGE_BUF_SIZE       = 66
	MAX_FIFO_SIZE          = 66
	TRUE                   = 1
	FALSE                  = 0

	ADDR_FIFO          = 0x00
	ADDR_OPMODE        = 0x01 // Operating modes
	ADDR_REGDATAMODUL  = 0x02
	ADDR_FDEVMSB       = 0x05
	ADDR_FDEVLSB       = 0x06
	ADDR_FRMSB         = 0x07
	ADDR_FRMID         = 0x08
	ADDR_FRLSB         = 0x09
	ADDR_VERSION       = 0x10
	ADDR_AFCCTRL       = 0x0B
	ADDR_LNA           = 0x18
	ADDR_RXBW          = 0x19
	ADDR_AFCFEI        = 0x1E
	ADDR_IRQFLAGS1     = 0x27
	ADDR_IRQFLAGS2     = 0x28
	ADDR_RSSITHRESH    = 0x29
	ADDR_PREAMBLELSB   = 0x2D
	ADDR_SYNCCONFIG    = 0x2E
	ADDR_SYNCVALUE1    = 0x2F
	ADDR_SYNCVALUE2    = 0X30
	ADDR_SYNCVALUE3    = 0x31
	ADDR_SYNCVALUE4    = 0X32
	ADDR_PACKETCONFIG1 = 0X37
	ADDR_PAYLOADLEN    = 0X38
	ADDR_NODEADDRESS   = 0X39
	ADDR_FIFOTHRESH    = 0X3C

	MASK_REGDATAMODUL_OOK = 0x08
	MASK_REGDATAMODUL_FSK = 0x00
	MASK_WRITE_DATA       = 0x80
	MASK_MODEREADY        = 0x80
	MASK_FIFONOTEMPTY     = 0x40
	MASK_FIFOLEVEL        = 0x20
	MASK_FIFOOVERRUN      = 0x10
	MASK_PACKETSENT       = 0x08
	MASK_TXREADY          = 0x20
	MASK_PACKETMODE       = 0x60
	MASK_MODULATION       = 0x18
	MASK_PAYLOADRDY       = 0x04

	/* Precise register description can be found on:
	 * www.hoperf.com/upload/rf/RFM69W-V1.3.pdf
	 * on page 63 - 74
	 */
	MODE_STANDBY         = 0x04        // Standby
	MODE_TRANSMITER      = 0x0C        // Transmiter
	MODE_RECEIVER        = 0x10        // Receiver
	VAL_REGDATAMODUL_FSK = 0x00        // Modulation scheme FSK
	VAL_REGDATAMODUL_OOK = 0x08        // Modulation scheme OOK
	VAL_FDEVMSB30        = 0x01        // frequency deviation 5kHz 0x0052 -> 30kHz 0x01EC
	VAL_FDEVLSB30        = 0xEC        // frequency deviation 5kHz 0x0052 -> 30kHz 0x01EC
	VAL_FRMSB434         = 0x6C        // carrier freq -> 434.3MHz 0x6C9333
	VAL_FRMID434         = 0x93        // carrier freq -> 434.3MHz 0x6C9333
	VAL_FRLSB434         = 0x33        // carrier freq -> 434.3MHz 0x6C9333
	VAL_FRMSB433         = 0x6C        // carrier freq -> 433.92MHz 0x6C7AE1
	VAL_FRMID433         = 0x7A        // carrier freq -> 433.92MHz 0x6C7AE1
	VAL_FRLSB433         = 0xE1        // carrier freq -> 433.92MHz 0x6C7AE1
	VAL_AFCCTRLS         = 0x00        // standard AFC routine
	VAL_AFCCTRLI         = 0x20        // improved AFC routine
	VAL_LNA50            = 0x08        // LNA input impedance 50 ohms
	VAL_LNA50G           = 0x0E        // LNA input impedance 50 ohms, LNA gain -> 48db
	VAL_LNA200           = 0x88        // LNA input impedance 200 ohms
	VAL_RXBW60           = 0x43        // channel filter bandwidth 10kHz -> 60kHz  page:26
	VAL_RXBW120          = 0x41        // channel filter bandwidth 120kHz
	VAL_AFCFEIRX         = 0x04        // AFC is performed each time RX mode is entered
	VAL_RSSITHRESH220    = 0xDC        // RSSI threshold 0xE4 -> 0xDC (220)
	VAL_PREAMBLELSB3     = 0x03        // preamble size LSB 3
	VAL_PREAMBLELSB5     = 0x05        // preamble size LSB 5
	VAL_SYNCCONFIG2      = 0x88        // Size of the Synch word = 2 (SyncSize + 1)
	VAL_SYNCCONFIG4      = 0x98        // Size of the Synch word = 4 (SyncSize + 1)
	VAL_SYNCVALUE1FSK    = 0x2D        // 1st byte of Sync word
	VAL_SYNCVALUE2FSK    = 0xD4        // 2nd byte of Sync word
	VAL_SYNCVALUE1OOK    = 0x80        // 1nd byte of Sync word
	VAL_PACKETCONFIG1FSK = 0xA2        // Variable length, Manchester coding, Addr must match NodeAddress
	VAL_PACKETCONFIG1OOK = 0           // Fixed length, no Manchester coding
	VAL_PAYLOADLEN255    = 0xFF        // max Length in RX, not used in Tx
	VAL_PAYLOADLEN64     = 0x40        // max Length in RX, not used in Tx
	VAL_PAYLOADLEN_OOK   = (13 + 8*17) // Payload Length
	VAL_NODEADDRESS01    = 0x04        // Node address used in address filtering
	VAL_FIFOTHRESH1      = 0x81        // Condition to start packet transmission: at least one byte in FIFO
	VAL_FIFOTHRESH30     = 0x1E        // Condition to start packet transmission: wait for 30 bytes in FIFO

	MSG_REMAINING_LEN = 0
	MSG_MANUF_ID      = 1
	MSG_PRODUCT_ID    = 2
	MSG_RESERVED_HI   = 3
	MSG_RESERVED_LO   = 4
	MSG_SENSOR_ID_2   = 5
	MSG_SENSOR_ID_1   = 6
	MSG_SENSOR_ID_0   = 7
	MSG_DATA_START    = 8
	MSG_ENCR_START    = MSG_SENSOR_ID_2
	MSG_OVERHEAD_LEN  = (MSG_DATA_START + 2)

	MAX_DATA_LENGTH = MESSAGE_BUF_SIZE

	/* OOK Message Parameters */
	OOK_BUF_SIZE           = 17
	OOK_MSG_ADDRESS_LENGTH = 10 /* 10 bytes in address */

	OT_JOIN_RESP = 0x6A
	OT_JOIN_CMD  = 0xEA

	OT_POWER      = 0x70
	OT_REACTIVE_P = 0x71

	OT_CURRENT    = 0x69
	OT_ACTUATE_SW = 0xF3
	OT_FREQUENCY  = 0x66
	OT_TEST       = 0xAA
	OT_SW_STATE   = 0x73

	OT_TEMP_SET    = 0xf4 /* Send new target temperature to driver board */
	OT_TEMP_REPORT = 0x74 /* Send externally read room temperature to motor board */

	OT_VOLTAGE = 0x76

	OT_EXERCISE_VALVE = 0xA3 /* Send exercise valve command to driver board.
	   Read diagnostic flags returned by driver board.
	   Send diagnostic flag acknowledgement to driver board.
	   Report diagnostic flags to the gateway.
	   Flash red LED once every 5 seconds if ‘battery dead’ flag
	   is set.
	     Unsigned Integer Length 0
	*/

	OT_REQUEST_VOLTAGE = 0xE2 /* Request battery voltage from driver board.
	   Report battery voltage to gateway.
	   Flash red LED 2 times every 5 seconds if voltage
	   is less than 2.4V
	     Unsigned Integer Length 0
	*/
	OT_REPORT_VOLTAGE = 0x62 /* Volts
	   Unsigned Integer Length 0
	*/

	OT_REQUEST_DIAGNOTICS = 0xA6 /*   Read diagnostic flags from driver board and report
	     these to gateway Flash red LED once every 5 seconds
	     if ‘battery dead’ flag is set
	     Unsigned Integer Length 0
	*/

	OT_REPORT_DIAGNOSTICS = 0x26

	OT_SET_VALVE_STATE = 0xA5 /*
	   Send a message to the driver board
	   0 = Set Valve Fully Open
	   1=Set Valve Fully Closed
	   2 = Set Normal Operation
	   Valve remains either fully open or fully closed until
	   valve state is set to ‘normal operation’.
	   Red LED flashes continuously while motor is running
	   terminated by three long green LED flashes when valve
	   fully open or three long red LED flashes when valve is
	   closed

	   Unsigned Integer Length 1
	*/

	OT_SET_LOW_POWER_MODE = 0xA4 /*
	   0=Low power mode off
	   1=Low power mode on

	   Unsigned Integer Length 1
	*/
	OT_IDENTIFY = 0xBF

	OT_SET_REPORTING_INTERVAL = 0xD2 /*
	      Update reporting interval to requested value

	   Unsigned Integer Length 2
	*/

	OT_CRC = 0x00

	encryptId uint16 = 0xf2
)

func NewHRF() (*HRF, error) {
	dev, err := spi.New(0, 1, spi.SPIMode0, 9600000)
	if err != nil {
		return nil, err
	}
	inst := &HRF{dev: dev}
	return inst, err
}

type Cmd struct {
	addr byte
	val  byte
}

func (self *HRF) ConfigFSK() error {
	regSetup := []Cmd{
		{ADDR_REGDATAMODUL, VAL_REGDATAMODUL_FSK}, // modulation scheme FSK
		{ADDR_FDEVMSB, VAL_FDEVMSB30},             // frequency deviation 5kHz 0x0052 -> 30kHz 0x01EC
		{ADDR_FDEVLSB, VAL_FDEVLSB30},             // frequency deviation 5kHz 0x0052 -> 30kHz 0x01EC
		{ADDR_FRMSB, VAL_FRMSB434},                // carrier freq -> 434.3MHz 0x6C9333
		{ADDR_FRMID, VAL_FRMID434},                // carrier freq -> 434.3MHz 0x6C9333
		{ADDR_FRLSB, VAL_FRLSB434},                // carrier freq -> 434.3MHz 0x6C9333
		{ADDR_AFCCTRL, VAL_AFCCTRLS},              // standard AFC routine
		{ADDR_LNA, VAL_LNA50},                     // 200ohms, gain by AGC loop -> 50ohms
		{ADDR_RXBW, VAL_RXBW60},                   // channel filter bandwidth 10kHz -> 60kHz  page:26
		//{ADDR_AFCFEI, 		VAL_AFCFEIRX},		// AFC is performed each time rx mode is entered
		//{ADDR_RSSITHRESH, 	VAL_RSSITHRESH220},	// RSSI threshold 0xE4 -> 0xDC (220)
		{ADDR_PREAMBLELSB, VAL_PREAMBLELSB3},       // preamble size LSB -> 3
		{ADDR_SYNCCONFIG, VAL_SYNCCONFIG2},         // Size of the Synch word = 2 (SyncSize + 1)
		{ADDR_SYNCVALUE1, VAL_SYNCVALUE1FSK},       // 1st byte of Sync word
		{ADDR_SYNCVALUE2, VAL_SYNCVALUE2FSK},       // 2nd byte of Sync word
		{ADDR_PACKETCONFIG1, VAL_PACKETCONFIG1FSK}, // Variable length, Manchester coding, Addr must match NodeAddress
		{ADDR_PAYLOADLEN, VAL_PAYLOADLEN64},        // max Length in RX, not used in Tx
		{ADDR_NODEADDRESS, VAL_NODEADDRESS01},      // Node address used in address filtering
		{ADDR_FIFOTHRESH, VAL_FIFOTHRESH1},         // Condition to start packet transmission: at least one byte in FIFO
		{ADDR_OPMODE, MODE_RECEIVER},               // Operating mode to Receiver
	}
	for _, cmd := range regSetup {
		err := self.regW(cmd.addr, cmd.val)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *HRF) WaitFor(addr byte, mask byte, val bool) {
	cnt := 0
	for {
		cnt += 1 // Uncomment to wait in a loop finite amount of time
		if cnt > 100000 {
			panic("timeout inside a while for addr")
			// log4c_category_warn(hrflog, "timeout inside a while for addr %02x\n", addr);
			// break
		}
		ret := self.regR(addr)
		if val {
			if (ret & mask) == mask {
				break
			}
		} else {
			if (ret & mask) == 0 {
				break
			}
		}
	}
}

func (self *HRF) ClearFifo() {
	for {
		if self.regR(ADDR_IRQFLAGS2)&MASK_FIFONOTEMPTY == 0 {
			break
		}
		self.regR(ADDR_FIFO)
	}
}

func (self *HRF) GetVersion() byte {
	return self.regR(ADDR_VERSION)
}

type Message struct {
	ManuId   byte
	ProdId   byte
	SensorId uint32
}

func decrypt(pid, pip uint16, data []byte) {
	ran := (pid << 8) ^ pip
	for i := range data {
		for j := 0; j < 5; j += 1 {
			if ran&1 == 1 {
				ran = (ran >> 1) ^ 62965
			} else {
				ran = ran >> 1
			}
		}
		data[i] = (byte(ran) ^ data[i] ^ 90)
	}
}

func (self *HRF) ReceiveFSKMessage(encryptionId byte, productId byte, manufacturerId byte) *Message {
	if self.regR(ADDR_IRQFLAGS2)&MASK_PAYLOADRDY == MASK_PAYLOADRDY {
		length := self.regR(ADDR_FIFO)
		data := make([]byte, length)
		for i := 0; i < int(length); i += 1 {
			data[i] = self.regR(ADDR_FIFO)
		}
		pip := uint16(data[2])<<8 | uint16(data[3])
		decrypt(encryptId, pip, data[4:])
		message := Message{
			ManuId:   data[0],
			ProdId:   data[1],
			SensorId: uint32(data[4])<<24 | uint32(data[5])<<16 | uint32(data[6]),
		}
		for i := 7; i < len(data); i += 2 {
			paramId := data[i]
			if paramId == 0 {
				// end of parameterss
				break
			}

			fmt.Println("Parameter:", paramId)
			typeDesc := data[i+1]
			dlen := typeDesc & 0x0f
			value := data[i+2 : i+2+int(dlen)]
			fmt.Println("Value:", value)
		}
		fmt.Println(hex.Dump(data))
		return &message
	}

	return nil
}

func (self *HRF) regR(addr byte) byte {
	buf := []byte{addr, 0}
	self.dev.Write(buf)
	self.dev.Read(buf)
	return buf[1]
}

func (self *HRF) regW(addr byte, val byte) error {
	buf := []byte{addr, val}
	_, err := self.dev.Write(buf)
	self.dev.Read(buf) // ignored
	return err
}
