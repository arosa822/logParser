package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	//"log"
	"flag"
	"os"
	"strings"
)

// flag stuct which contains an array of flags
type Flags struct {
	Flags []Flag `json:"flags"`
}

// flag struct which contains user defined flags
type Flag struct {
	Start    string `json:"start"`
	End      string `json:"end"`
	FailFlag string `json:"failFlag"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (flags *Flags) loadFlags() {
	jsonFile, err := os.Open("flags.json")
	check(err)
	fmt.Println("Grabing flags from flags.json...")
	defer jsonFile.Close()

	// read in the json file as a byte array
	byteValue, _ := ioutil.ReadAll(jsonFile)

	//unmarshal the byte Array from above
	json.Unmarshal(byteValue, &flags)

}

func main() {
	var (
		flags             Flags
		status            bool
		errorLog          []string
		state, lineNumber byte
		errorLines        []byte
	)

	// first byte holds how many events
	// second byte holds number of passing events
	// third holds the number of failed events
	s := make([]int, 3)

	arg := os.Args[1]

	// load flags from json file
	flags.loadFlags()
	fmt.Println("json flag parser successfull.")

	fmt.Printf("\nFLAGS>>>>\nStart: %s\nStop: %s\nFail: %s\n",
		flags.Flags[0].Start, flags.Flags[0].End, flags.Flags[0].FailFlag)

	// parse the logFile
	flag.Parse()
	file, err := os.Open(arg)
	check(err)
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		lineNumber++

		if err == io.EOF {
			break
		}

		// check to see if we are starting a sequence
		if state == 0 && strings.Contains(line, flags.Flags[0].Start) {
			state = 1
			status = false
		}

		if state == 1 && strings.Contains(line, flags.Flags[0].FailFlag) {
			status = true
			errorLog = append(errorLog, line)
			errorLines = append(errorLines, lineNumber)
		}

		// check for end of sequence

		if state == 1 && strings.Contains(line, flags.Flags[0].End) {
			s[0]++
			state = 0
			// reset status if error was found
			if status != true {
				status = true
				// add to fail counter
				s[2]++
			} else {
				s[1]++
			}
		}

	}

	for n, i := range errorLog {
		fmt.Printf("Line: %v\n", errorLines[n])
		fmt.Printf(">> %s\n", i)
	}

	fmt.Printf("\nNumber of cycles: %v\n", s[0])
	fmt.Printf("\nFailed events: %v\n", s[1])
	fmt.Printf("\nSuccessfull events: %v\n", s[2])
	fmt.Println("\nFlagged lines>>>")

}
