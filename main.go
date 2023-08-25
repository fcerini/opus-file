package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
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
		packet := line[9:]
		//log.Print(packet)

		packet = strings.ReplaceAll(packet, " ", "")
		//log.Print(packet)

		bytes, _ := hex.DecodeString(packet)
		//log.Print(bytes)

		handleVoiceBroadcast(bytes)
	}
	readFile.Close()

	recordAudio()
}
