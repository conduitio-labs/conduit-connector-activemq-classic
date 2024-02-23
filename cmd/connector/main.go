package main

import (
	"fmt"
	"time"

	"github.com/go-stomp/stomp"
)

func main() {
	go sendMessage()
	go receiveMessage("listener 1")
	go receiveMessage("listener 2")

	select {}
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

func sendMessage() {
	conn, err := stomp.Dial("tcp", "localhost:61613")
	assert(err)
	defer conn.Disconnect()

	for {
		err = conn.Send(
			"/queue/SampleQueue",
			"text/plain",
			[]byte("Message to Queue"), nil)
		assert(err)

		time.Sleep(1 * time.Second)
	}
}

func receiveMessage(name string) {
	conn, err := stomp.Dial("tcp", "localhost:61613")
	assert(err)

	defer conn.Disconnect()

	sub, err := conn.Subscribe("/queue/SampleQueue", stomp.AckAuto)
	assert(err)

	for {
		msg := <-sub.C
		fmt.Println(name, "received", string(msg.Body))
	}
}
