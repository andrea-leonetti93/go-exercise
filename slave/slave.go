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

	slaveData := new(st.SlaveData)
	slaveServer := com.RegisterRPCNamedService("SlaveData", slaveData)
	//todo passare riferimento alle strutture di join del master
	jr := st.JoinRequest{"localhost", *slavePort}
	var msgFromServer = &st.SecondSlaveAddress{}

	error := server.Call("JoinRequest1.Join1", jr, msgFromServer)
	if error != nil {
		log.Fatal("Error in JoinRequest1.Join1: ", error)
	}
	//fmt.Printf("JoinRequest1.Join1: %s\n", msgFromServer.ResponseMessage)
	fmt.Printf("Address slave: %s\n", msgFromServer)

	slaveAddress := "localhost" + *slavePort
	fmt.Printf("slave address: %s\n", slaveAddress)
	l := com.CreatePortListener(slaveAddress)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
		}
		go slaveServer.ServeConn(conn)
	}
}
