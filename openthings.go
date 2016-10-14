package ener314

import "fmt"

const (
	/* OpenThings definitions */
	engManufacturerId = 0x04 // Energenie Manufacturer Id
	eTRVProductId     = 0x3  // Product ID for eTRV
	encryptId         = 0xf2 // Encryption ID for eTRV

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
)

type Record interface {
	String() string
}

type Join struct{}

func (j Join) String() string {
	return "Join"
}

type Temperature struct {
	Value float32
}

func (t Temperature) String() string {
	return fmt.Sprintf("Temperature{%f}", t.Value)
}

type Voltage struct {
	Value float32
}

func (v Voltage) String() string {
	return fmt.Sprintf("Voltage{%f}", v.Value)
}

type UnhandledRecord struct {
	ID    byte
	Type  byte
	Value []byte
}

func (t UnhandledRecord) String() string {
	return fmt.Sprintf("Unhandled{%x,%x,%v}", t.ID, t.Type, t.Value)
}

type Message struct {
	ManuId   byte
	ProdId   byte
	SensorId uint32
	Records  []Record
}

func (m *Message) String() string {
	records := ""
	for _, record := range m.Records {
		if len(records) > 0 {
			records += ","
		}
		records += fmt.Sprint(record)
	}
	return fmt.Sprintf("{ManuId:%d ProdId:%d SensorId:%06x Records:[%s]}", m.ManuId, m.ProdId, m.SensorId, records)
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

func decodeFloat32(typeDesc byte, value []byte) float32 {
	return float32(value[0]) + float32(value[1])/0xff
}

func DecodePacket(data []byte) *Message {
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

		typeDesc := data[i+1]
		dlen := typeDesc & 0x0f
		value := data[i+2 : i+2+int(dlen)]
		i += int(dlen)

		var record Record
		switch paramId {
		case OT_JOIN_CMD:
			record = Join{}
		case OT_TEMP_REPORT:
			record = Temperature{decodeFloat32(typeDesc, value)}
		case OT_VOLTAGE:
			record = Voltage{decodeFloat32(typeDesc, value)}
		default:
			record = UnhandledRecord{paramId, typeDesc, value}
		}
		message.Records = append(message.Records, record)
	}
	return &message
}
