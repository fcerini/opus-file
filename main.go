package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)

var RecBuffMux sync.Mutex
var RecorderBuffer = make(map[uint32]*PacketBuffer) // Key: session

func main() {

	filePath := "debug.txt"
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
		log.Print(line)

		buf := []byte{0}
		handleVoiceBroadcast(buf)
	}

	readFile.Close()

}
