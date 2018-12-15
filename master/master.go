package main

import (
	com "go-exercise/libs/comunication"
	st "go-exercise/libs/structs"
	"log"
	"strings"
)

// Work ciao
type Work struct{}

type infoSlave struct {
	address string
	port    string
}

var slaveConnected []infoSlave

// WordCount fnaculoi
func (t *Work) WordCount(msg *st.StringMsg, result *st.StringMsg) error {
	result.Text = strings.ToUpper(msg.Text)
	return nil
}

// JoinRequest message from slave to join master
type JoinRequest struct {
	address string
	port    string
}

// ResponseRequest send from master to slave
type ResponseRequest struct {
	responseMessage string
}

// Join use by slave to contact master
func (j *JoinRequest) Join(join *JoinRequest, result *ResponseRequest) error {
	s := infoSlave{join.address, join.port}
	slaveConnected = append(slaveConnected, s)
	result.responseMessage = "Join done"
	return nil
}

func main() {

	work := new(Work)

	//server := rpc.NewServer()
	server := com.RegisterRPCNamedService("Work", work)
	join := new(JoinRequest)
	server.RegisterName("Join", join)

	l := com.CreatePortListener("localhost:1234")

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
		}
		go server.ServeConn(conn)
	}

}
