package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"go.bug.st/serial"
)

var ComPort string
var Speed int = 115200

var MacroSet struct {
	ComPort    string   `json:"port"`
	Speed      int      `json:"speed"`
	Commands   []string `json:"commands"`
	StartDelay int      `json:"startDelay"`
	Delay      int      `json:"commandDelay"`
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

	mode := &serial.Mode{
		BaudRate: Speed,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	ComPort = MacroSet.ComPort
	port, err := serial.Open(ComPort, mode)
	if err != nil {
		fmt.Println("error opening Serial")
		log.Fatal(err)
	}
	err = port.SetMode(mode)
	if err != nil {
		fmt.Println("error setting Serial mode")
		log.Fatal(err)
	}

	fmt.Println("Waiting for startup")
	time.Sleep(time.Duration(MacroSet.StartDelay) * time.Second)
	/*
		n, err := port.Write([]byte("G28 X Y\n"))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Sent %v bytes\n", n)
	*/

	for _, command := range MacroSet.Commands {
		fmt.Println(command)
		_, err := port.Write([]byte(command + "\n"))
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Duration(MacroSet.Delay) * time.Millisecond)
	}

	fmt.Println("Macro ran succesfully")
	fmt.Println("Thank you for using Gcode Macro")
}

func getInputs() {

	jsonFile, err := os.Open("macro.json")
	if err != nil {
		fmt.Println("error opening Json")
		log.Fatal(err)
	}
	fmt.Println("Json opened")
	defer jsonFile.Close()

	jsonParser := json.NewDecoder(jsonFile)
	err = jsonParser.Decode(&MacroSet)
	if err != nil {
		fmt.Println("error parsing json")
		log.Fatal(err)
	}
	fmt.Println("Json Parsed")
}
