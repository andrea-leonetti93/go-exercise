package structs

import (
	"fmt"
	"sort"
	"strings"
	"unicode"

	hash "go-exercise/hash"
)

// StringMsg fanc
type StringMsg struct {
	Text string
}

//FileToSend cioa
type FileToSend struct {
	File string
}

// Dictionary to map word found in the text
type Dictionary struct {
	name  string
	value int
}

/////////////////master

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

//var Alphabet []string = ["a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z"]

// Join use by slave to contact master
func (j *JoinRequest) Join(join *JoinRequest, result *ResponseRequest) error {
	//s := infoSlave{join.Address, join.Port}
	SlaveConnected = append(SlaveConnected, *join)
	result.ResponseMessage = "Join done"
	fmt.Printf("slave connected list %v", SlaveConnected)
	return nil
}

////////slave

// SlaveResponse send on channel from slave to master
type SlaveResponse struct {
	WordHashMap hash.ValueHashtable
}

//SlaveData : data used from slave to count
type SlaveData struct {
	textToParse string
}

//LavoroSlave ciao
func (s *SlaveData) LavoroSlave(text string, result SlaveResponse) {

	result.WordHashMap = TextParse(text)
}

func dataOrder(h hash.ValueHashtable) {
	//var keys []int
	keys := make([]string, 0, h.Size())
	for k := range h.Items {
		keys = append(keys, string(k))
	}
	//sort.Ints(keys)
	sort.Strings(keys)
	// To perform the opertion you want
	for _, k := range keys {
		fmt.Println(k, h.Items[hash.Key(k)])
	}

	/*for _, k := range keys {
		fmt.Println("Key:", k, "Value:", hash.Items[k])
	}*/
}

// TextParse ciao
func TextParse(text string) hash.ValueHashtable {
	var splittedString []string
	word := ""
	h := hash.ValueHashtable{}
	for _, r := range text {

		if !unicode.IsLetter(r) && word != "" {
			//key := hash.Key(word)
			word = strings.ToLower(word)
			splittedString = append(splittedString, word)
			if h.IfWordExist(hash.Key(word)) != 0 {
				h.Increment(hash.Key(word))
			} else {
				//v := hash.Value{word, 1}

				h.Put(hash.Key(word), 1)
			}
			word = ""
		} else {
			if !unicode.IsSpace(r) {
				word += string(r)
			}
		}
	}
	//dataOrder(h)
	return h
}

//SlaveJob slave work
func (s *SlaveData) SlaveJob(data SlaveData) []Dictionary {
	//d := []Dictionary

	return nil
}
