//https://golang.org/src/net/rpc/server.go
//https://golang.org/src/net/rpc/client.go

package main

import (
	"fmt"
	com "go-exercise/libs/comunication"
	st "go-exercise/libs/structs"
	"log"
	"net/rpc"
	"time"
	"unicode"
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

func parseSpace(text string, dimChunk int, start int) (int, string) {
	r := []rune(text)
	if start+dimChunk > len(text) {
		return start, text[start : len(text)-1]
	}
	for i := start + dimChunk; i < len(text); i++ {
		if !unicode.IsLetter(r[i]) {
			dimChunk = dimChunk + (i - dimChunk)
			break
		}
	}
	return dimChunk, text[start:dimChunk]
}

func splitText(text string, numOfSlave int) []string {
	s := make([]string, numOfSlave)
	dimChunk := len(text) / numOfSlave
	start := 0
	for i := 0; i < numOfSlave; i++ {
		//s1 := text[:dimChunk]
		pos, str := parseSpace(text, dimChunk, start)
		start = pos + 1
		s[i] = str
	}
	return s
}

//Map divide work between slaves
func Map(fts *st.FileToSend, numOfSlave int) {
	fileSize := len(fts.File)
	slaveCall := new(rpc.Call)
	c := st.Counter{}
	c.Count = 0
	//chunkSize := fileSize / numOfSlave
	fmt.Printf("numofslave: %d\n", numOfSlave)
	fmt.Printf("filesize: %d\n", fileSize)
	//slaveResult := make([]hash.ValueHashtable, numOfSlave)
	slaveResult := make([]st.SlaveResponse, numOfSlave)
	fmt.Printf("entro nel ciclo per mandare il lavoro agli slave\n")
	s := splitText(fts.File, numOfSlave)
	for i := 0; i < len(s); i++ {
		fmt.Printf("s: %s\n", s[i])
	}
	for i := 0; i < numOfSlave; i++ {
		//creazione parametri per lavoro master-slave
		slaveResult[i].Counter = c
		slaveText := st.SlaveData{}
		//slaveText.TextToParse = fts.File //[i*chunkSize : chunkSize]
		slaveText.TextToParse = s[i]
		//richiesta servizio slave
		server := com.ConnectToHost("localhost" + st.SlaveConnected[i].Port)
		slaveCall = server.Go("SlaveData.LavoroSlave", slaveText, &slaveResult[i], nil)
		/*if slaveCall == nil {
			log.Fatal("Error in SlaveData.LavoroSlave: ", slaveCall)
		}*/
	}
	//wait group
	replyCall := <-slaveCall.Done
	fmt.Printf("replycall", replyCall)
	fmt.Printf("entro in wait group\n")
	counter := 0
	for {
		time.Sleep(5 * time.Second)
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
		} else {
			counter = 0
		}
	}
	/*
	* questo va eliminato perché da result arrivano tutte le hash già ordinate
	* in slaveResult[]
	 */
	/*finalHash := hash.ValueHashtable{}
	for i := 0; i < numOfSlave; i++ {
		for k, v := range slaveResult[i].WordHashMap.Items {
			if finalHash.IfWordExist(hash.Key(k)) != 0 {
				finalHash.Increment1(hash.Key(k), int(v))
			} else {
				finalHash.Put(hash.Key(k), v)
			}
		}
	}
	st.DataOrder(finalHash)*/
	fmt.Printf("restituite tutte le hash table da tutti gli slave di primo livello\n")
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
	join1 := new(st.JoinRequest)
	join2 := new(st.JoinRequest)
	server.RegisterName("JoinRequest1", join1)
	server.RegisterName("JoinRequest2", join2)

	l := com.CreatePortListener(st.MasterAddress)
	n := 3
	y := 26 / n
	x := 102
	sot := 122 - x
	div := sot / y
	resto := sot % n
	if div == 0 {
		fmt.Printf("primo gruppo")
	} else if resto == 0 {
		fmt.Printf("resto = 0 --> gruppo numero: %d\n", div)
	} else if resto != 0 {
		fmt.Printf("resto != 0 --> gruppo numero: %d\n", div+1)
	}
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
