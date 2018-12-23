//https://golang.org/src/net/rpc/server.go
//https://golang.org/src/net/rpc/client.go

package main

import (
	"fmt"
	com "go-exercise/libs/comunication"
	st "go-exercise/libs/structs"
	"log"
	"time"
)

//BUFFERSIZE buffer chunk
const BUFFERSIZE = 1024

//Work ciao
type Work struct {
	textToParse string
}

/*
type Call struct {
	ServiceMethod string      // The name of the service and method to call.
	Args          interface{} // The argument to the function (*struct).
	Reply         interface{} // The reply from the function (*struct).
	Error         error       // After completion, the error status.
	Done          chan *Call  // Strobes when call is complete.
}*/

//Map divide work between slaves
func Map(fts *st.FileToSend, numOfSlave int) {
	fileSize := len(fts.File)
	//slaveCall := new(rpc.Call)
	c := st.Counter{}
	c.Count = 0
	//chunkSize := fileSize / numOfSlave
	fmt.Printf("numofslave: %d\n", numOfSlave)
	fmt.Printf("filesize: %d\n", fileSize)
	//slaveResult := make([]hash.ValueHashtable, numOfSlave)
	slaveResult := make([]st.SlaveResponse, numOfSlave)
	fmt.Printf("entro nel ciclo per mandare il lavoro agli slave\n")
	for i := 0; i < numOfSlave; i++ {
		//creazione parametri per lavoro master-slave
		slaveResult[i].Counter = c
		slaveText := st.SlaveData{}
		slaveText.TextToParse = fts.File //[i*chunkSize : chunkSize]

		//richiesta servizio slave
		server := com.ConnectToHost("localhost" + st.SlaveConnected[i].Port)
		slaveCall := server.Go("SlaveData.LavoroSlave", slaveText, &slaveResult[i], nil)
		/*		if slaveCall == nil {
				log.Fatal("Error in SlaveData.LavoroSlave: ", slaveCall)
			}*/
		replyCall := <-slaveCall.Done
		fmt.Printf("replycall", replyCall)
	}
	//wait group

	//	fmt.Printf("replycall", replyCall)
	fmt.Printf("entro in wait group\n")
	counter := 0
	for {
		time.Sleep(10 * time.Second)
		fmt.Printf("slaveresult counter: %d\n", slaveResult[0].Counter.Count)
		for i := 0; i < numOfSlave; i++ {
			if slaveResult[i].Counter.Count == 1 {
				counter++
				fmt.Printf("counter aggiornato: %d\n", counter)
			}
		}
		if counter == numOfSlave {
			fmt.Printf("esco dalla wait group\n")
			break
		}
		//jobCompleted := <-slaveCall.Done
	}
	st.DataOrder(slaveResult[0].WordHashMap)
	/*for h, v := range slaveResult {
		for
		a[k] = v
	}*/
}

// WordCount fnaculoi
func (t *Work) WordCount(fts *st.FileToSend, result *st.StringMsg) error {
	fmt.Printf("textToParse: %s\n", string(fts.File))
	for {
		if len(st.SlaveConnected) == st.NumberOfSlave {
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
