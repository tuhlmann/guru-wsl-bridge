package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Configuration read from a file
type Configuration struct {
	GOPATHOnWindows string
	GOPATHOnLinux   string
}

func main() {

	config, err := readConfig()
	if err != nil {
		panic(err)
	}

	binary, err := exec.LookPath("wsl.exe")
	if err != nil {
		panic(err)
	}

	// Hey future self, if you need to debug
	// setLogOutput(fmt.Sprintf("%s\\guru-wsl-bridge.log", os.Getenv("USERPROFILE")))

	args := []string{"guru.sh"}

	for _, s := range os.Args[1:] {
		args = append(args, findAndReplaceGOPATH(config, s))
	}

	// log.Printf("modified args: %+v", args)

	command := exec.Command(binary, args...)
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout

	stdinInfo, _ := os.Stdin.Stat()

	if !(stdinInfo.Mode()&os.ModeCharDevice != 0) {
		inReader := bufio.NewReader(os.Stdin)
		var passedInputRaw []byte
		var guruInput []byte

		for {
			input, err := inReader.ReadByte()
			if err != nil && err == io.EOF {
				break
			}
			passedInputRaw = append(passedInputRaw, input)
		}

		passedInput := bufio.NewReader(bytes.NewReader(passedInputRaw))

		for {
			// file path
			line, err := passedInput.ReadString('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			if line == "\n" {
				guruInput = append(guruInput, []byte(line)...)
				break
			}

			newLine := findAndReplaceGOPATH(config, line)

			// log.Printf("Append line: %s", newLine)
			guruInput = append(guruInput, []byte(newLine)...)

			// buffer size
			line, err = passedInput.ReadString('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			size, err := strconv.Atoi(line[:len(line)-1])
			if err != nil {
				panic(err)
			}

			guruInput = append(guruInput, []byte(line)...)

			buffer := make([]byte, size)
			bytesRead, err := io.ReadFull(passedInput, buffer)
			if bytesRead != len(buffer) {
				panic(err)
			}

			guruInput = append(guruInput, buffer...)

		}

		command.Stdin = bytes.NewBuffer(guruInput)

		// log.Printf("\n\npiped into stdin:\n%v\n", string(guruInput))
	}

	if err := command.Start(); nil != err {
		log.Fatalf("Error starting program: %s, %s", command.Path, err.Error())
	}
	command.Wait()

}

func setLogOutput(logPath string) {
	f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v", err))
	}
	defer f.Close()
	log.SetOutput(f)
}

func readConfig() (*Configuration, error) {
	configFilePath := fmt.Sprintf("%s\\.guru-wsl-bridge.json", os.Getenv("USERPROFILE"))
	configFile, err := os.Open(configFilePath)
	if err != nil {
		return nil, err
	}
	var config Configuration
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	return &config, err
}

func findAndReplaceGOPATH(config *Configuration, line string) string {
	// log.Printf("Find %s in %s", strings.ToLower(line), strings.ToLower(config.GOPATHOnWindows))
	pos := strings.Index(strings.ToLower(line), strings.ToLower(config.GOPATHOnWindows))
	newLine := ""
	if pos > -1 {
		// log.Printf("Found at %v", pos)
		if pos > 0 {
			newLine = newLine + line[0:pos-1]
		}
		newLine = newLine + config.GOPATHOnLinux
		pos2 := pos + len(config.GOPATHOnWindows)
		newLine = newLine + line[pos2:]
		// line = strings.Replace(line, "c:\\entw\\go", config.GOPATHOnLinux, -1)
	} else {
		newLine = line
	}
	return strings.Replace(newLine, "\\", "/", -1)
}
