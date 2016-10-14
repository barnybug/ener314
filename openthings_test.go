package ener314

import "fmt"

func ExampleDecodePacketJoin() {
	packet := []byte{0x04, 0x03, 0x65, 0xce, 0xa0, 0x97, 0x51, 0xac, 0xc2, 0xf4, 0xa2, 0x19}
	message := DecodePacket(packet)
	fmt.Printf("%s\n", message)
	// Output:
	// {ManuId:4 ProdId:3 SensorId:09007f Records:[Join]}
}

func ExampleDecodePacketVoltage() {
	packet := []byte{0x04, 0x03, 0x13, 0x04, 0x20, 0x3b, 0x19, 0xd5, 0x8c, 0xf1, 0x5f, 0xf1, 0xd3, 0x7b}
	message := DecodePacket(packet)
	fmt.Printf("%s\n", message)
	// Output:
	// {ManuId:4 ProdId:3 SensorId:09007f Records:[Voltage{3.121569}]}
}

func ExampleDecodePacketTemp() {
	packet := []byte{0x04, 0x03, 0x0f, 0x42, 0x89, 0x00, 0x3a, 0x46, 0x9c, 0xa6, 0xe2, 0x35, 0x1f, 0xdc}
	message := DecodePacket(packet)
	fmt.Printf("%s\n", message)
	// Output:
	// {ManuId:4 ProdId:3 SensorId:09007f Records:[Temperature{17.701962}]}
}
