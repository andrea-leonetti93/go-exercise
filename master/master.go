package main

import (
	"fmt"
	com "go-exercise/libs/comunication"
	st "go-exercise/libs/structs"
	"log"
)

//BUFFERSIZE buffer chunk
const BUFFERSIZE = 1024

// Work ciao
/*type Work struct{}

// WordCount fnaculoi
func (t *Work) WordCount(msg *st.StringMsg, result *st.StringMsg) error {
	result.Text = strings.ToUpper(msg.Text)
	return nil
}*/

func assignJob() {

	var numOfSlave = len(st.SlaveConnected)
	var lettersToAssign = 26 / numOfSlave
	for _, s := range st.SlaveConnected {
		job := st.SlaveData{}
	}
}

//Work ciao
type Work struct {
	textToParse string
}

// WordCount fnaculoi
func (t *Work) WordCount(fts *st.FileToSend, result *st.StringMsg) error {
	fmt.Printf("textToParse: %s\n", string(fts.File))
	for {
		if len(st.SlaveConnected) == 5 {
			break
		}
	}
	assignJob()
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
