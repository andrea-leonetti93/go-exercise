package main

import (
	"flag"
	"fmt"
	com "go-exercise/libs/comunication"
	st "go-exercise/libs/structs"
	"log"
)

var serverAddress = flag.String("server", st.MasterAddress, "master server")
var slavePort = flag.String("slave", st.SlavePort, "slave port")

func main() {

	flag.Parse()

	server := com.ConnectToHost(*serverAddress)

	//todo passare riferimento alle strutture di join del master
	jr := st.JoinRequest{"localhost", *slavePort}
	var msgFromServer = &st.ResponseRequest{}

	err := server.Call("JoinRequest.Join", jr, msgFromServer)
	if err != nil {
		log.Fatal("Error in JoinRequest.Join: ", err)
	}
	fmt.Printf("JoinRequest.Join: %s\n", msgFromServer.ResponseMessage)

}
