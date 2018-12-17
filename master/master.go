package main

import (
	com "go-exercise/libs/comunication"
	st "go-exercise/libs/structs"
	"log"
	"strings"
)

// Work ciao
type Work struct{}

// WordCount fnaculoi
func (t *Work) WordCount(msg *st.StringMsg, result *st.StringMsg) error {
	result.Text = strings.ToUpper(msg.Text)
	return nil
}

func AssignJob() {

	var numOfSlave = len(st.SlaveConnected)
	var lettersToAssign = 26 / numOfSlave
	for _, s := range st.SlaveConnected {
		job := st.SlaveData{}
	}
}

func main() {

	work := new(Work)

	//server := rpc.NewServer()
	server := com.RegisterRPCNamedService("Work", work)
	join := new(st.JoinRequest)
	server.RegisterName("JoinRequest", join)

	l := com.CreatePortListener(st.MasterAddress)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
		}
		go server.ServeConn(conn)
	}

}
