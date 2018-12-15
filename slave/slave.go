package main

import (
	"flag"
	"fmt"
	com "go-exercise/libs/comunication"
	"log"
)

var serverAddress = flag.String("server", "localhost:1234", "master server")

func main() {

	server := com.ConnectToHost(*serverAddress)

	//todo passare riferimento alle strutture di join del master

	err := server.Call("Work.WordCount", msgToSend, &msgFromServer)
	if err != nil {
		log.Fatal("Error in Work.WordCount: ", err)
	}
	fmt.Printf("Work.WordCount: %s %s\n", msgToSend.Text, msgFromServer.Text)

}
