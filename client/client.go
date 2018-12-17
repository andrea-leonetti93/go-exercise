package main

import (
	"flag"
	"fmt"
	"log"

	com "go-exercise/libs/comunication"

	st "go-exercise/libs/structs"
)

var textToParse = flag.String("text", "", "testo da parsare")
var serverAddress = flag.String("server", st.MasterAddress, "master server")

func main() {
	flag.Parse()

	fmt.Println("File to parse: ", *textToParse)

	//connection to a server (master server)
	server := com.ConnectToHost(*serverAddress)

	msgToSend := &st.StringMsg{*textToParse}
	var msgFromServer = &st.StringMsg{}

	err := server.Call("Work.WordCount", msgToSend, &msgFromServer)
	if err != nil {
		log.Fatal("Error in Work.WordCount: ", err)
	}
	fmt.Printf("Work.WordCount: %s %s\n", msgToSend.Text, msgFromServer.Text)

}
