package comunication

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"reflect"
)

//se uso la prima lettera di una funzione maiuscola -> pubbliche
//se uso la prima lettera di una funzione minuscola -> private

// ConnectToHost connect to the specified host
func ConnectToHost(addr string) *rpc.Client {
	//esegue la connessione tramite chiamata rpc
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return client
}

// CreatePortListener create a listener port
func CreatePortListener(serverAddr string) net.Listener {
	ln, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatal("Listen error: ", err)
		return nil
	}
	return ln
}

// RegisterRPCNamedService register the RPC interface of reciver
func RegisterRPCNamedService(serviceName string, receiver interface{}) *rpc.Server {
	//newServer create a new istance of server
	server := rpc.NewServer()
	//registerName use to expose the rpc call on the server already created
	err := server.RegisterName(serviceName, receiver)
	if err != nil {
		rcvrType := reflect.TypeOf(receiver)
		log.Fatal("Format of ", rcvrType, " is not correct: ", err)
	}
	return server
}

// x = <-chan recive      chan <- x send     <-chan ingora risultato

//CommunicationMasterSlave with Vice
func CommunicationMasterSlave(ctx context.Context, countedWord <-chan []byte,
	fileChunk chan<- []byte, errs <-chan error) {

	for {
		select {
		case <-ctx.Done():
			log.Println("finished")
			return

		case err := <-errs:
			log.Println("an error occurred:", err)

		/*case name := <-countedWord:
			//greeting := "Hello " + string(name)
			var x = <-countedWord
		}*/

	}
}
