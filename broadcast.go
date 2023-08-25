package main

import (
	"encoding/binary"
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
		return
	}

	// buffer
	// 0:4 session uint32
	// 5 cant canales byte
	// lista de canales []uint32

	var channelsList []uint32
	senderSession := binary.LittleEndian.Uint32(buffer[0:4])
	channelsInPacket := buffer[4]

	for i := byte(0); i < channelsInPacket; i++ {
		pos1 := i*4 + 5
		pos2 := pos1 + 4
		channelsList = append(channelsList, binary.LittleEndian.Uint32(buffer[pos1:pos2]))
	}

	headerEnd := channelsInPacket*4 + 5

	packet := buffer[headerEnd:]

	// uso append para copiar los slices, sino quedan como referencia
	voicePacket := append([]byte{}, packet...)

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
