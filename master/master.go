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

func main() {

	work := new(Work)

	//server := rpc.NewServer()
	server := com.RegisterRPCNamedService("Work", work)

	l := com.CreatePortListener("localhost:1234")

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
		}
		go server.ServeConn(conn)
	}

}
