package ener314

import (
	"fmt"
	"io"
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

type Identify struct{}

func (i Identify) String() string {
	return "Identify"
}

func (i Identify) Encode(buf io.ByteWriter) {
	buf.WriteByte(OT_IDENTIFY)
	buf.WriteByte(0)
}

type TargetTemperature struct{}

func (i TargetTemperature) String() string {
	return "TargetTemperature"
}

func (i TargetTemperature) Encode(buf io.ByteWriter) {
	buf.WriteByte(OT_TEMP_SET)
	buf.WriteByte(0)
}
