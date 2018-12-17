package structs

import "fmt"

// StringMsg fanc
type StringMsg struct {
	Text string
}

/////////////////

// JoinRequest message from slave to join master
type JoinRequest struct {
	Address string
	Port    string
}

// ResponseRequest send from master to slave
type ResponseRequest struct {
	ResponseMessage string
}

//SlaveConnected : list of slaves connected
var SlaveConnected []JoinRequest

//SlaveData : data used from slave to count
type SlaveData struct {
	lettersToCheck []string
	textToParse    string
}


var Alphabet []string = ["a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z"]

// Join use by slave to contact master
func (j *JoinRequest) Join(join *JoinRequest, result *ResponseRequest) error {
	//s := infoSlave{join.Address, join.Port}
	SlaveConnected = append(SlaveConnected, *join)
	result.ResponseMessage = "Join done"
	fmt.Printf("slave connected list %v", SlaveConnected)
	return nil
}
