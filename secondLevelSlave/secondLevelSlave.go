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

	secondSlave := new(st.SecondSlaveAddress)
	slaveServer := com.RegisterRPCNamedService("MixShuffle", secondSlave)
	//todo passare riferimento alle strutture di join del master
	jr := st.JoinRequest{"localhost", *slavePort}
	var msgFromServer = &st.SlaveResponse{}

	err := server.Call("JoinRequest2.Join2", jr, msgFromServer)
	if err != nil {
		log.Fatal("Error in JoinRequest2.Join2: ", err)
	}
	fmt.Printf("JoinRequest2.Join: %s\n")

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
