package main

import (
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

func handleVoiceBroadcast(buffer []byte) {

	if len(buffer) < 30 {
		log.Fatal("len(buffer) < 30", buffer)
	}

	channelsList := []uint32{1}
	senderSession := uint32(1)

	// uso append para copiar los slices, sino quedan como referencia
	voicePacket := append([]byte{}, buffer...)

	RecBuffMux.Lock()
	pb, exists := RecorderBuffer[senderSession]

	if !exists {
		RecorderBuffer[senderSession] = &PacketBuffer{
			firstTimestamp: time.Now(),
			lastTimestamp:  time.Now(),
			senderSession:  senderSession,
			channelsList:   channelsList,
			packetBuff:     []OpusPacket{getOpusPacket(voicePacket)},
		}
		log.Printf("Audio record started at:%v Sess:%v", time.Now(), senderSession)
	} else {
		pb.packetBuff = append(pb.packetBuff, getOpusPacket(voicePacket))
		pb.lastTimestamp = time.Now()
	}

	RecBuffMux.Unlock()
}
