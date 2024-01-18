package main

import (
	"fmt"
	"log"

	"go.bug.st/serial"
)

func main() {
	fmt.Println("Gcode Macro by Npi")

	ports, err := serial.GetPortsList()

	if err != nil {
		log.Fatal(err)
	}

	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}

	for _, port := range ports {
		fmt.Printf("Found port: %v\n", port)
	}

	mode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open("COM3", mode)
	if err != nil {
		log.Fatal(err)
	}

	err = port.SetMode(mode)
	if err != nil {
		log.Fatal(err)
	}

	n, err := port.Write([]byte("G28 X Y\n"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Sent %v bytes\n", n)
}
