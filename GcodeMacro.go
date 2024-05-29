package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"go.bug.st/serial"
)

var MacroSet struct {
	ComPort    string   `json:"port"`
	Speed      int      `json:"speed"`
	Commands   []string `json:"commands"`
	StartDelay int      `json:"startDelay"`
	Delay      int      `json:"commandDelay"`
}

func main() {
	fmt.Print("Gcode Macro by Npi\n\n")
	fileList := getAvailableFiles()
	var input string
	fmt.Print("Please select a macro to run:\n\n")
	for i, file := range fileList {
		fmt.Println(i, ": ", file)
	}
	fmt.Println()
	fmt.Scanln(&input)
	j, err := strconv.Atoi(input)
	if err != nil || j > len(fileList)-1 {
		log.Fatal("invalid input")
	}
	fmt.Println("running " + fileList[j])
	os.Exit(1)

	getInputs(fileList[j])

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
		BaudRate: MacroSet.Speed,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(MacroSet.ComPort, mode)
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

func getInputs(fileName string) {

	jsonFile, err := os.Open(fileName)
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

func getAvailableFiles() []string {
	var list []string

	dir, err := os.Open(".")
	if err != nil {
		log.Fatal(err)
	}
	defer dir.Close()
	files, err := dir.Readdirnames(50)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {

		if len(f) > 5 {
			if f[len(f)-5:] == ".json" {
				list = append(list, f)
			}
		}
	}
	return list
}
