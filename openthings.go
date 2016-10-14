package ener314

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strconv"
)

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
	Encode(buf io.ByteWriter)
}

type Join struct{}

func (j Join) String() string {
	return "Join"
}

func (j Join) Encode(buf io.ByteWriter) {
	buf.WriteByte(OT_JOIN_CMD)
	buf.WriteByte(0)
}

type Temperature struct {
	Value float64
}

func (t Temperature) String() string {
	return fmt.Sprintf("Temperature{%f}", t.Value)
}

func (t Temperature) Encode(buf io.ByteWriter) {
	buf.WriteByte(OT_TEMP_REPORT)
	// TODO - encode signed fixed .8
}

type Voltage struct {
	Value float64
}

func (v Voltage) String() string {
	return fmt.Sprintf("Voltage{%f}", v.Value)
}

func (v Voltage) Encode(buf io.ByteWriter) {
	buf.WriteByte(OT_VOLTAGE)
	// TODO - encode signed fixed .8
}

type UnhandledRecord struct {
	ID    byte
	Type  byte
	Value []byte
}

func (t UnhandledRecord) String() string {
	return fmt.Sprintf("Unhandled{%x,%x,%v}", t.ID, t.Type, t.Value)
}

func (t UnhandledRecord) Encode(buf io.ByteWriter) {
	// Unhandled
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

func decodeFixedPoint(value []byte, mantissa uint, signed bool) float64 {
	var ret float64
	sign := false
	if signed && len(value) > 0 && (value[0]&0x80 != 0) {
		value[0] = value[0] & 0x7f
		sign = true
	}

	for _, b := range value {
		ret = ret*256 + float64(b)
	}
	div := 1 << mantissa
	if sign {
		ret = -ret
	}
	return ret / float64(div)
}

func decodeFloat64(typeDesc byte, value []byte) float64 {
	switch typeDesc >> 4 {
	case 0x0: // Unsigned x.0 normal integer
		return decodeFixedPoint(value, 0, false)
	case 0x1: // Unsigned x.4 fixed point integer
		return decodeFixedPoint(value, 4, false)
	case 0x2: // Unsigned x.8 fixed point integer
		return decodeFixedPoint(value, 8, false)
	case 0x3: // Unsigned x.12 fixed point integer
		return decodeFixedPoint(value, 12, false)
	case 0x4: // Unsigned x.16 fixed point integer
		return decodeFixedPoint(value, 16, false)
	case 0x5: // Unsigned x.20 fixed point integer
		return decodeFixedPoint(value, 20, false)
	case 0x6: // Unsigned x.24 fixed point integer
		return decodeFixedPoint(value, 24, false)
	case 0x7: // Characters
		f64, _ := strconv.ParseFloat(string(value), 32)
		return f64
	case 0x8: // Signed x.0 normal integer
		return decodeFixedPoint(value, 0, true)
	case 0x9: // Signed x.8 fixed point integer
		return decodeFixedPoint(value, 8, true)
	case 0xa: // Signed x.16 fixed point integer
		return decodeFixedPoint(value, 16, true)
	case 0xb: // Signed x.24 fixed point integer
		return decodeFixedPoint(value, 24, true)
	case 0xc: // Enumeration
		// Just treat as unsigned integer
		return decodeFixedPoint(value, 0, false)
	case 0xd, 0xe: // Reserved
	case 0xf: // IEEE754-2008 floating point
		// untesed - 32 or 64?
		var ret float64
		buf := bytes.NewReader(value)
		binary.Read(buf, binary.LittleEndian, &ret)
		return ret
	}
	return 0
}

func DecodePacket(data []byte) (*Message, error) {
	pip := uint16(data[2])<<8 | uint16(data[3])
	decrypt(encryptId, pip, data[4:])
	return DecodeUnencryptedPacket(data)
}

var ErrShortPacket = errors.New("Short or corrupt packet")
var ErrCRCFail = errors.New("CRC fail")

func DecodeUnencryptedPacket(data []byte) (*Message, error) {
	ln := len(data)
	if ln < 10 {
		// absolute minimum:
		// 2 manufacturer, product
		// 2 encryption pip
		// 3 sensor id
		// 1 no records
		// 2 crc
		return nil, ErrShortPacket
	}

	// check CRC
	crc := uint16(data[ln-2])<<8 + uint16(data[ln-1])
	expected := calculateCRC(data[4 : ln-2])
	if crc != expected {
		return nil, ErrCRCFail
	}

	message := Message{
		ManuId:   data[0],
		ProdId:   data[1],
		SensorId: uint32(data[4])<<16 | uint32(data[5])<<8 | uint32(data[6]),
	}
	// i + one byte + crc
	for i := 7; true; i += 2 {
		paramId := data[i]
		if paramId == 0 {
			// end of parameterss
			break
		}
		if i >= ln-4 {
			// at least [code] [typedesc] [crc] [crc]
			return nil, ErrShortPacket
		}

		typeDesc := data[i+1]
		dlen := typeDesc & 0x0f
		if i+2+int(dlen)+2 >= ln {
			// at least [code] [typedesc] [..variable..] [crc] [crc]
			return nil, ErrShortPacket
		}

		value := data[i+2 : i+2+int(dlen)]
		i += int(dlen)

		// value length check
		switch paramId {
		case OT_TEMP_REPORT, OT_VOLTAGE:
			if dlen == 0 {
				return nil, ErrShortPacket
			}
		}

		var record Record
		switch paramId {
		case OT_JOIN_CMD:
			record = Join{}
		case OT_TEMP_REPORT:
			record = Temperature{decodeFloat64(typeDesc, value)}
		case OT_VOLTAGE:
			record = Voltage{decodeFloat64(typeDesc, value)}
		default:
			record = UnhandledRecord{paramId, typeDesc, value}
		}
		message.Records = append(message.Records, record)
	}
	return &message, nil
}

func calculateCRC(data []byte) uint16 {
	var rem uint16

	for _, d := range data {
		rem = rem ^ uint16(d)<<8
		for bit := 8; bit > 0; bit -= 1 {
			if rem&(1<<15) == 0 {
				rem = rem << 1
			} else {
				rem = (rem << 1) ^ 0x1021
			}
		}
	}
	return rem
}

func encrypt(pid, pip uint16, data []byte) {
	// reversable: encrypt is decrypt
	decrypt(pid, pip, data)
}

func EncodeMessage(message *Message) []byte {
	pip := uint16(rand.Uint32())
	data := EncodeData(message, pip)
	encrypt(encryptId, pip, data[4:])
	return data
}

func EncodeData(message *Message, pip uint16) []byte {
	var buf bytes.Buffer
	buf.WriteByte(message.ManuId)
	buf.WriteByte(message.ProdId)
	buf.WriteByte(byte(pip >> 8))
	buf.WriteByte(byte(pip))

	buf.WriteByte(byte(message.SensorId >> 16))
	buf.WriteByte(byte(message.SensorId >> 8))
	buf.WriteByte(byte(message.SensorId))

	for _, record := range message.Records {
		record.Encode(&buf)
	}
	buf.WriteByte(0) // end of params

	// only encrypted data is CRCed
	crc := calculateCRC(buf.Bytes()[5:])
	buf.WriteByte(byte(crc >> 8))
	buf.WriteByte(byte(crc))
	return buf.Bytes()
}
