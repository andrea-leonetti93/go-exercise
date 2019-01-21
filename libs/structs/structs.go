package structs

import (
	"fmt"
	hash "go-exercise/hash"
	com "go-exercise/libs/comunication"
	"net/rpc"
	"regexp"
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
	SecondLevelAddress []JoinRequest
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
	result.SecondLevelAddress = SecondLevelSlaveAddress
	fmt.Printf("slave connected list %v\n", SlaveConnected)
	return nil
}

// Join2 use by second level slave to contact master
func (j *JoinRequest) Join2(join *JoinRequest, result *ResponseRequest) error {
	//s := infoSlave{join.Address, join.Port}
	SecondLevelSlaveAddress = append(SecondLevelSlaveAddress, *join)
	result.ResponseMessage = "Join done"
	fmt.Printf("slave connected list %v\n", SecondLevelSlaveAddress)
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
	fmt.Printf("hash 1\n")
	hashList[0].PrintTable()
	fmt.Printf("hash 2\n")
	hashList[1].PrintTable()
	fmt.Printf("hash 3\n")
	hashList[2].PrintTable()
	result.Counter.lock.Lock()
	result.Counter.Count++
	fmt.Printf("counter: %d\n", result.Counter.Count)
	defer result.Counter.lock.Unlock()
	secondSlaveResult := make([]SlaveResponse, NumberOfSlave)
	fmt.Printf("number of second level slave: %d\n", len(SecondLevelSlaveAddress))
	for i := 0; i < NumberOfSlave; i++ {
		add := SecondLevelSlaveAddress[i].Address
		port := SecondLevelSlaveAddress[i].Port
		fmt.Printf("indirizzo slave secondo livello: N=%d, add=%s\n", i, add+port)
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
func DataOrder(h hash.ValueHashtable) hash.ValueHashtable {

	htOrdered := hash.ValueHashtable{}

	//slice di appoggio per l'ordinamento delle chiavi
	keys := make([]string, 0, h.Size())
	for k := range h.Items {
		keys = append(keys, string(k))
	}

	sort.Strings(keys)
	// To perform the opertion you want
	for _, k := range keys {
		htOrdered.Put(hash.Key(k), h.Items[hash.Key(k)])
		//fmt.Println(k, h.Items[hash.Key(k)])
	}
	return htOrdered
}

//ParseTextToSlice tokenize input text
func ParseTextToSlice(rawText string) []string {
	// torna true se c e' una lettera
	f := func(c rune) bool {
		return !unicode.IsLetter(c)
	}
	// crea uno slice di stringhe usando f come funzione di split
	tokensOfStrings := strings.FieldsFunc(rawText, f)
	return tokensOfStrings
}

//NormalizeString string to lower case
func NormalizeString(rawString string) string {
	return strings.ToLower(rawString)
}

//GetHashTableIndex Index finds the correct hashtable index to put word in
func GetHashTableIndex(word string) int {
	primoRangeRX := regexp.MustCompile("^[a-i]{1}")
	secondoRangeRX := regexp.MustCompile("^[j-r]{1}")
	terzoRangeRX := regexp.MustCompile("^[s-z]{1}")

	switch {
	case primoRangeRX.MatchString(word):
		return 0
	case secondoRangeRX.MatchString(word):
		return 1
	case terzoRangeRX.MatchString(word):
		return 2
	}
	return -1
}

//TextParse fills hashtable
func TextParse(rawText string) []hash.ValueHashtable {

	h := make([]hash.ValueHashtable, NumberOfSlave)
	tokens := ParseTextToSlice(rawText)

	for _, word := range tokens {
		//lower case string
		word = NormalizeString(word)
		tableIndex := GetHashTableIndex(word)

		//if table index is valid
		if tableIndex != -1 {
			if h[tableIndex].IfWordExist(hash.Key(word)) != 0 {
				h[tableIndex].Increment(hash.Key(word))
			} else {
				h[tableIndex].Put(hash.Key(word), 1)
			}
		}

	}
	return h
}

///////// second level slave

// CounterSecondLevel boh
type CounterSecondLevel struct {
	Counter        int
	lock           sync.RWMutex
	finalHashTable hash.ValueHashtable
}

// HashCounter boh
var HashCounter CounterSecondLevel

// SortAndReduce order the elements of the hashtable
func (sr *SlaveResponse) SortAndReduce(partialHash *hash.ValueHashtable, result *SlaveResponse) error {
	fmt.Printf("entro in SortAndReduce!!!!!!\n")
	fmt.Printf("printing hash \n")
	partialHash.PrintTable()
	//contatore che conta gli accessi
	HashCounter.lock.Lock()
	// copiare hash table in quella globale e incrementare counter
	for k, v := range partialHash.Items {
		if HashCounter.finalHashTable.IfWordExist(hash.Key(k)) != 0 {
			HashCounter.finalHashTable.IncrementByValue(hash.Key(k), v)
		} else {
			HashCounter.finalHashTable.Put(hash.Key(k), 1)
		}
	}
	HashCounter.Counter++
	fmt.Printf("creata hash finale!!\n")
	fmt.Printf("valore counter: %d", HashCounter.Counter)
	HashCounter.finalHashTable.PrintTable()
	defer HashCounter.lock.Unlock()
	//aspettare che l'hash table finale sia riempita da tutti i processi--> counter == 3
	fmt.Printf("entro nel for del counter")
	for {
		if HashCounter.Counter == 3 {
			break
		}
	}
	// ordiniamo i valori per la hash table finale
	HashCounter.finalHashTable = DataOrder(HashCounter.finalHashTable)
	//la rimandiamo indietro a tutti gli slave del primo livello
	result.Counter.Count++
	fmt.Printf("\nRisultato ordinamento chiavi:\n")
	HashCounter.finalHashTable.PrintTable()
	result.WordHashMap = HashCounter.finalHashTable
	return nil
}
