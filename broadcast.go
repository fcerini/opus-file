package main

import (
	"fmt"
	"log"
	"time"
)

func getOpusPacket(voicePacket []byte) (packet OpusPacket) {
	// TODO: Controlar largo de cada decode o posibles errores a la hora de obtener la trama.
	voicePacket = voicePacket[1:]
	session, n := varintDecode(voicePacket)
	voicePacket = voicePacket[n:]
	sequence, n := varintDecode(voicePacket)
	voicePacket = voicePacket[n:]
	len, n := varintDecode(voicePacket)
	voicePacket = voicePacket[n:]

	// Opus audio packets set the 13th bit in the size field as the terminator.
	audioLength := int(len) &^ 0x2000

	_ = session // Not used, se podria usar y no recibir del cliente.
	packet.sequence = sequence
	packet.payload = voicePacket[:audioLength]
	return
}

var GloLinea = 0
var GloSecuencia = int64(0)

func handleVoiceBroadcast(buffer []byte) {
	GloLinea++

	if len(buffer) < 30 {
		log.Fatal("len(buffer) < 30", buffer)
	}

	channelsList := []uint32{1}
	senderSession := uint32(1)

	// uso append para copiar los slices, sino quedan como referencia
	voicePacket := append([]byte{}, buffer...)

	packetBuff := getOpusPacket(voicePacket)

	if GloSecuencia == 0 {
		GloSecuencia = packetBuff.sequence - 4
	}

	if (packetBuff.sequence - GloSecuencia) != 4 {
		fmt.Printf("WARN %v sec %v dif %v \n", GloLinea, packetBuff.sequence, packetBuff.sequence-GloSecuencia)
	} else {
		fmt.Printf(".... %v sec %v dif %v \n", GloLinea, packetBuff.sequence, packetBuff.sequence-GloSecuencia)
	}

	GloSecuencia = packetBuff.sequence

	RecBuffMux.Lock()
	pb, exists := RecorderBuffer[senderSession]

	if !exists {
		RecorderBuffer[senderSession] = &PacketBuffer{
			firstTimestamp: time.Now(),
			lastTimestamp:  time.Now(),
			senderSession:  senderSession,
			channelsList:   channelsList,
			packetBuff:     []OpusPacket{packetBuff},
		}
	} else {
		pb.packetBuff = append(pb.packetBuff, packetBuff)
		pb.lastTimestamp = time.Now()
	}

	RecBuffMux.Unlock()
}
