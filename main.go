package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

var RecBuffMux sync.Mutex
var RecorderBuffer = make(map[uint32]*PacketBuffer) // Key: session

func main() {

	filePath := "audio.txt"
	if len(os.Args) > 1 {
		filePath = os.Args[1]
	}
	readFile, err := os.Open(filePath)

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		values := strings.Split(line, ",")

		bytes := make([]byte, len(values))
		for i, value := range values {
			num, _ := strconv.Atoi(value)
			bytes[i] = byte(num)
		}

		handleVoiceBroadcast(bytes)
	}
	readFile.Close()

	recordAudio()
}
