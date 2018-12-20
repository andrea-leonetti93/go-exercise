package structs

import (
	"fmt"
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

//SlaveData : data used from slave to count
type SlaveData struct {
	lettersToCheck []string
	textToParse    string
}

// TextParse ciao
func TextParse(text string) []string {
	var splittedString []string
	word := ""
	h := hash.ValueHashtable{}
	for _, r := range text {

		if !unicode.IsLetter(r) && word != "" {
			key := hash.Key(word)
			splittedString = append(splittedString, word)
			if h.IfWordExist(key) != 0 {
				h.Increment(key)
			} else {
				v := hash.Value{word, 1}
				h.Put(key, v)
			}
			/*
				* if word exists: count +=1
				else key = word, count = 1
			*/
			word = ""
			continue
		} else {
			word += string(r)
		}
	}
	/*f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}*/
	//splittedString = strings.FieldsFunc(text, f)
	//fmt.Printf("Fields are: %q", strings.FieldsFunc("  foo1;bar2,baz3...", f))
	//fmt.Printf("splittedString: %s\n", splittedString[5])
	//splittedString1 := strings.Split(splittedString, " ")
	fmt.Println("Final Hash: ")
	for key, value := range h.Items {
		fmt.Printf("[key: %d, value : %s %d ]", key, value.Word, value.Count)
	}
	println("\n")
	return splittedString
}

//SlaveJob slave work
func (s *SlaveData) SlaveJob(data SlaveData) []Dictionary {
	//d := []Dictionary

	return nil
}
