package main

import (
	"log"
	"time"
)

type OpusPacket struct {
	sequence int64
	payload  []byte
}

type PacketBuffer struct {
	firstTimestamp time.Time
	lastTimestamp  time.Time
	senderSession  uint32
	channelsList   []uint32
	packetBuff     []OpusPacket
}

func recordAudio() {

	pb := RecorderBuffer[1]
	fileName := "audio.opus"

	writer, err := New(fileName, 48000, 1) // sampleRate - channelCount

	if err != nil {
		log.Fatalf("Create outputdir folder: %v. Error: %s", fileName, err)
	}

	// INFO: El "numero magico" es el framesize (480 1920). Falta ver de donde sale ese numero.
	// 		   Una posibilidad es que sale de freq sampleo * packet time ya que verifica estos dos
	//         numeros verificados: con framesize de 1920 -> 40ms y con framesize de 480 -> 10ms.

	//TODO: Mejorar la forma de escribir?? mas paquetes por pagina.

	// Para calcular el "framesize" utilizo las secuencias (cuando son 10ms -> 0,1,2,3,4,5 , cuando son 40ms
	// -> 0,4,8,12,16,20. Esto lo multiplico por 10ms y por 48k (frec sampleo)).
	inc := (pb.packetBuff[1].sequence - pb.packetBuff[0].sequence) * 10 * 48
	inc_acumm := int64(0)
	for _, p := range pb.packetBuff {
		inc_acumm += inc
		writer.WritePacket(p.payload, uint32(inc_acumm))
	}
	writer.Close()
	log.Printf("Audio saved at:%v Sess:%v", time.Now(), pb.senderSession)

}
