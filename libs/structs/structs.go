package structs

import (
	"fmt"
	hash "go-exercise/hash"
	com "go-exercise/libs/comunication"
	"net/rpc"
	"sort"
	"strings"
	"sync"
	"unicode"
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

// SecondSlaveAddress
type SecondSlaveAddress struct {
	SeconeLevelAddress []JoinRequest
}

//SlaveConnected : list of slaves connected
var SlaveConnected []JoinRequest

// SecondLevelSlaveAddress : list of second slaves connected
var SecondLevelSlaveAddress []JoinRequest

//var Alphabet []string = ["a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z"]

// Join1 use by first level slave to contact master
func (j *JoinRequest) Join1(join *JoinRequest, result *SecondSlaveAddress) error {
	//s := infoSlave{join.Address, join.Port}
	SlaveConnected = append(SlaveConnected, *join)
	fmt.Printf("SecondLevelSlave %v\n", SecondLevelSlaveAddress)
	result.SeconeLevelAddress = SecondLevelSlaveAddress
	fmt.Printf("slave connected list %v", SlaveConnected)
	return nil
}

// Join2 use by second level slave to contact master
func (j *JoinRequest) Join2(join *JoinRequest, result *ResponseRequest) error {
	//s := infoSlave{join.Address, join.Port}
	SecondLevelSlaveAddress = append(SecondLevelSlaveAddress, *join)
	result.ResponseMessage = "Join done"
	fmt.Printf("slave connected list %v", SecondLevelSlaveAddress)
	return nil
}

////////slave

//Counter contatore per slave
type Counter struct {
	Count int
	lock  sync.RWMutex
}

// SlaveResponse send on channel from slave to master
type SlaveResponse struct {
	WordHashMap hash.ValueHashtable
	Counter     Counter
}

//SlaveData : data used from slave to count
type SlaveData struct {
	TextToParse string
}

//LavoroSlave ciao
func (s *SlaveData) LavoroSlave(text SlaveData, result *SlaveResponse) error {
	fmt.Printf("entro in lavoro slave\n")
	fmt.Printf("text: %s/n", text.TextToParse)
	slaveCall := new(rpc.Call)
	//result.WordHashMap.Items = TextParse(text.TextToParse)
	// 3 hashtable divise per i tre slave di secondo livello
	fmt.Printf("start text to pars\n")
	hashList := TextParse(text.TextToParse) //[]valueHashtable
	fmt.Printf("finito txt to parse\n")
	result.Counter.lock.Lock()
	result.Counter.Count++
	fmt.Printf("counter: %d\n", result.Counter.Count)
	defer result.Counter.lock.Unlock()
	secondSlaveResult := make([]SlaveResponse, NumberOfSlave)
	for i := 0; i < NumberOfSlave; i++ {
		add := SecondLevelSlaveAddress[i].Address
		port := SecondLevelSlaveAddress[i].Port

		server := com.ConnectToHost(add + port)
		slaveCall = server.Go("SlaveResponse.SortAndReduce", hashList[i], &secondSlaveResult[i], nil)
	}
	replyCall := <-slaveCall.Done
	fmt.Printf("replycall", replyCall)
	counter := 0
	// aspetta che tutti gli slave di secondo livello tornino tutte le hash table
	for {
		for i := 0; i < NumberOfSlave; i++ {
			if secondSlaveResult[i].Counter.Count == 1 {
				counter++
			}
		}
		if counter == 3 {
			break
		} else {
			counter = 0
		}
	}
	// fa il merge di tutte le hash table con tutte le lettere dalla a alla z
	for i := 0; i < len(secondSlaveResult); i++ {
		for k, v := range secondSlaveResult[i].WordHashMap.Items {
			result.WordHashMap.Put(hash.Key(k), v)
		}
	}
	fmt.Printf("esco da lavoro slave e restituisco il risultato al master\n")
	return nil
}

//DataOrder all the data in the hashmap
func DataOrder(h hash.ValueHashtable) {
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

// TextParse ciao  returned hash.ValueHashtable
func TextParse(text string) []hash.ValueHashtable {
	var splittedString []string
	word := ""
	h := make([]hash.ValueHashtable, NumberOfSlave)
	//h := hash.ValueHashtable{}
	for _, r := range text {

		if !unicode.IsLetter(r) && word != "" {
			//key := hash.Key(word)
			word = strings.ToLower(word)
			splittedString = append(splittedString, word)
			i := 0
			if int(r) <= 105 {
				i = 0
			} else if 106 <= int(r) && int(r) <= 114 {
				i = 1
			} else {
				i = 2
			}
			if h[i].IfWordExist(hash.Key(word)) != 0 {
				h[i].Increment(hash.Key(word))
			} else {
				//v := hash.Value{word, 1}

				h[i].Put(hash.Key(word), 1)
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

///////// second level slave

// CounterSecondLevel
type CounterSecondLevel struct {
	Counter        int
	lock           sync.RWMutex
	finalHashTable hash.ValueHashtable
}

// HashCounter
var HashCounter CounterSecondLevel

// MixShuffle order the elements of the hashtable
func (sr *SlaveResponse) SortAndReduce(partialHash *SlaveResponse, result *SlaveResponse) error {
	//contatore che conta gli accessi
	HashCounter.lock.Lock()
	// copiare hash table in quella globale e incrementare counter
	for k := range partialHash.WordHashMap.Items {
		if HashCounter.finalHashTable.IfWordExist(hash.Key(k)) != 0 {
			HashCounter.finalHashTable.Increment(hash.Key(k))
		} else {
			HashCounter.finalHashTable.Put(hash.Key(k), 1)
		}
	}
	HashCounter.Counter++
	defer HashCounter.lock.Unlock()
	//aspettare che l'hash table finale sia riempita da tutti i processi--> counter == 3
	for {
		if HashCounter.Counter == 3 {
			break
		}
	}
	// ordiniamo i valori per la hash table finale
	DataOrder(HashCounter.finalHashTable)
	//la rimandiamo indietro a tutti gli slave del primo livello
	result.Counter.Count++
	result.WordHashMap = HashCounter.finalHashTable
	return nil
}
