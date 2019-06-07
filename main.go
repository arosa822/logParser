package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// flags strict which contains an array of flags
type Flags struct {
	Flags []Flag `json:"flags"`
}

// Flag struct which contains all the flags we are looking for
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
func main() {

	var (
		flags             Flags
		status            bool
		errorLog          []string
		state, lineNumber byte
		errLines          []byte
	)

	// first byte holds the number of cycles
	// second byte holds the number of successfull attempts
	// third byte holds the number of flagged events
	s := make([]int, 4)

	arg := os.Args[1]

	// open our jsonFile
	jsonFile, err := os.Open("flags.json")

	check(err)

	fmt.Println("Reading flags from .json file...")

	defer jsonFile.Close()

	//read our opened jsonFile as a bytearray
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// unmarshal our byteArray which contains our json's file content into
	// flags defined above
	json.Unmarshal(byteValue, &flags)

	file, err := os.Open(arg)

	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++

		// check for cycle start event
		if state == 0 && strings.Contains(line, flags.Flags[0].Start) {
			state = 1
			status = false
		}

		// check for failure
		if state == 1 && strings.Contains(line, flags.Flags[0].FailFlag) {
			status = true
			errorLog = append(errorLog, line)
			errLines = append(errLines, lineNumber)

		}
		// check for end of line
		if state == 1 && strings.Contains(line, flags.Flags[0].End) {
			s[0]++
			state = 0
			if status != true {
				// reset status
				status = true
				s[1]++
				// add a count to the last bit
			} else {
				s[2]++
			}
		}
		// reset states status indicates a failure, state tells us if we are in
		// the middle of a process

	}

	fmt.Printf("\nNumber of cycles: %v\n", s[0])
	fmt.Printf("Successfull Attempts %v\n", s[1])
	fmt.Printf("Flagged Attempts: %v\n\n", s[2])

	fmt.Println("Flagged lines:")
	fmt.Println("########################################")

	for n, i := range errorLog {
		fmt.Printf("Line #: %v\n", errLines[n])
		fmt.Printf(">> %s\n", i)
	}

}
