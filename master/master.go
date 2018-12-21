//https://golang.org/src/net/rpc/server.go
//https://golang.org/src/net/rpc/client.go

package main

import (
	"fmt"
	com "go-exercise/libs/comunication"
	st "go-exercise/libs/structs"
	"log"
)

//BUFFERSIZE buffer chunk
const BUFFERSIZE = 1024


//Work ciao
type Work struct {
	textToParse string
}

type Call struct {

	ServiceMethod string      // The name of the service and method to call.

	Args          interface{} // The argument to the function (*struct).

	Reply         interface{} // The reply from the function (*struct).

	Error         error       // After completion, the error status.

	Done          chan *Call  // Strobes when call is complete.

}

//Map divide work between slaves
func Map(fts *st.FileToSend, numOfSlave int) {
	fileSize := len(fts.File)
	slaveCall := new(Call)
	chunkSize := fileSize / numOfSlave
	slaveResult := [numOfSlave]hash.ValueHashTable
	for i := 0; i < numOfSlave; i++ {
		//creazione parametri per lavoro master-slave

		slateText := fts.File[i*chunkSize : chunkSize]

		//richiesta servizio slave
		server := com.ConnectToHost("localhost" + st.SlaveConnected[i].Port)
		slaveCall, err:= server.Go("SlaveData.LavoroSlave", slateText, slaveResult[i],nil)
		if err != nil {
			log.Fatal("Error in SlaveData.LavoroSlave: ", err)
		}
	}
	for{
		jobCompleted := slaveCall.Done()
	}
	//wait group
	/*for h, v := range slaveResult {
		for 
		a[k] = v
	}*/
}

// WordCount fnaculoi
func (t *Work) WordCount(fts *st.FileToSend, result *st.StringMsg) error {
	fmt.Printf("textToParse: %s\n", string(fts.File))
	for {
		if len(st.SlaveConnected) == 5 {
			break
		}
	}
	Map(fts, len(st.SlaveConnected))
	result.Text = "Text received"
	return nil
}

func main() {

	work := new(Work)

	//server := rpc.NewServer()
	server := com.RegisterRPCNamedService("Work", work)
	join := new(st.JoinRequest)
	server.RegisterName("JoinRequest", join)

	l := com.CreatePortListener(st.MasterAddress)

	timeout := 0

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
		}
		if timeout == 10 {
			//launch process
			timeout = 0
		}
		go server.ServeConn(conn)
	}

}
