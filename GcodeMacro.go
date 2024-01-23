package main

import (
	"fmt"
	"log"
	"os"

	"go.bug.st/serial"
)

var ComPort string
var Speed int = 115200

type MacroSet struct {
	ComPort  string   `json:"port"`
	Speed    int      `json:"speed"`
	Commands []string `json:"commands"`
}

func main() {
	fmt.Println("Gcode Macro by Npi")
	getInputs()

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
	ComPort = "COM3"

	mode := &serial.Mode{
		BaudRate: Speed,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(ComPort, mode)
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

func getInputs() {
	jsonFile, err := os.Open("macro.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Json opened")
	defer jsonFile.Close()
}
