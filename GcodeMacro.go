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

var Completed bool = false
var FirstRun bool = true

func main() {
	fmt.Print("Gcode Macro by Npi\n\n")
	var port serial.Port
	fileList := getAvailableFiles()
	var input string
	for !Completed {
		fmt.Print("Please select a macro to run:\n\n")
		for i, file := range fileList {
			fmt.Println(i, ": ", file)
		}
		fmt.Println()
		fmt.Scanln(&input)
		j, err := strconv.Atoi(input)
		if err != nil || j > len(fileList) {
			log.Fatal("invalid input")
		}
		if j == 0 {
			Completed = true
			break
		}
		fmt.Println("running " + fileList[j])

		getInputs(fileList[j])
		if FirstRun {
			port = SetupSerial()
			defer port.Close()
			FirstRun = false
		}

		//DummyWrite()
		WriteSerial(port)

	}

	fmt.Println("Thank you for using Gcode Macro")
}

func ReadSerial(Port serial.Port) {
	// Read and print the response

	buff := make([]byte, 100)
	for {
		// Reads up to 100 bytes
		n, err := Port.Read(buff)
		if err != nil {
			log.Fatal(err)
			break
		}
		if n == 0 {
			fmt.Println("\nEOF")
			break
		}
		fmt.Printf("%v", string(buff[:n]))
	}

}

func WriteSerial(port serial.Port) {
	for _, command := range MacroSet.Commands {
		fmt.Println(command)
		_, err := port.Write([]byte(command + "\n"))
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Duration(MacroSet.Delay) * time.Millisecond)
	}
	fmt.Println("Macro completed")

}

func SetupSerial() serial.Port {

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

	Port, err := serial.Open(MacroSet.ComPort, mode)
	if err != nil {
		fmt.Println("error opening Serial")
		log.Fatal(err)
	}
	err = Port.SetMode(mode)
	if err != nil {
		fmt.Println("error setting Serial mode")
		log.Fatal(err)
	}

	fmt.Println("Waiting for startup")
	time.Sleep(time.Duration(MacroSet.StartDelay) * time.Second)
	return Port

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
	list = append(list, "Exit Program")
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

func DummyWrite() {
	for _, command := range MacroSet.Commands {
		fmt.Println(command)
		time.Sleep(time.Duration(MacroSet.Delay) * time.Millisecond)
	}
	fmt.Println("Macro completed")
}
