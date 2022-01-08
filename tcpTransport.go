package transports

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"
)

/** first step of encodding request (Data in string)*/
type reqDataString struct {
	Pattern string
	Data    string
}

/** second step of encodding request (Data in bite slice)*/
type reqEncoded struct {
	Pattern string
	Data    []byte
}

type RequestData interface {
}

/** RouterMap is a type for keys(patterns) and appropriate function handler
    recive:
	  conn - connection
	  options - specific for some function handler
	  data - in byte slice. Unmarshal in appropriate function handler
*/
type RouterMap map[string]func(conn net.Conn, options string, data []byte)

/** store for custom routes*/
var routes *RouterMap

/** Recive net.Listener and map with routes and appropriate function handler*/
func TCPRouter(ln net.Listener, routesDefault RouterMap) {

	routes = &routesDefault

	for {
		conn, err := ln.Accept()

		fmt.Println("TCPRouter conn - ", conn)

		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

/** Recive connection and map with routes and appropriate function handler*/
func handleConnection(conn net.Conn) {
	defer conn.Close()

	bs, err := ioutil.ReadAll(conn)

	if err != nil {
		fmt.Println(err)
	}

	/** Convert bite slice to string to get meta data*/
	incData := string(bs)

	fmt.Println("incData ---", incData)

	/** # is separator for meta data and data*/
	incDataSlice := strings.Split(incData, "#")

	/** @ is separator for message_length, pattern and optiopns in meta data*/
	messageLengthPatternOptionsMetaSlice := strings.Split(incDataSlice[0], "@")

	fmt.Println("messageLengthPatternOptionsMetaSlice0 ---", messageLengthPatternOptionsMetaSlice)

	/** get message_length message from meta data slice - that is length of data buffer*/
	messageLength := messageLengthPatternOptionsMetaSlice[0]

	/** get message_length from meta data slice - that is length of data buffer*/
	messageLengthSlice := strings.Split(messageLength, "=")

	/** get meta data length to separate meta data buffer from data buffer*/
	metaLength := len([]byte(incDataSlice[0])) + 1

	pattern := ""
	options := ""

	if len(messageLengthPatternOptionsMetaSlice) > 1 {
		/** get pattern and options from meta data slice*/
		pattern = strings.Split(messageLengthPatternOptionsMetaSlice[1], "=")[1]
		options = messageLengthPatternOptionsMetaSlice[2]
	}

	fmt.Println("options ---", options)

	/** determine what kind of handle use, stream or simple
	  if  "message_length" present inside income data than need to use handleStreamConnection
	*/
	switch messageLengthSlice[0] {
	case "message_length":
		handleStreamConnection(conn, pattern, options, bs[metaLength:])
	default:
		handleSimpleConnection(conn, bs[metaLength:])

	}
}

func handleSimpleConnection(conn net.Conn, data []byte) {

	// TODO: fix data from string to []byte
	jsonRequest := string(data)

	fmt.Println("jsonRequest ---", jsonRequest)

	var parsedRequestDataString reqDataString

	err := json.Unmarshal([]byte(jsonRequest), &parsedRequestDataString)

	if err != nil {
		fmt.Println("handleSimpleConnection: Request parsing error ---", err)

		conn.Write([]byte("Request parsing error"))
	}

	fmt.Println("handleConnection parsedRequest ---", parsedRequestDataString)

	var parsedRequest reqEncoded

	parsedRequest.Pattern = parsedRequestDataString.Pattern
	parsedRequest.Data = []byte(parsedRequestDataString.Data)

	router(conn, parsedRequest)

}

func handleStreamConnection(conn net.Conn, pattern string, optionsS string, data []byte) {
	defer conn.Close()

	streamRouter(conn, pattern, optionsS, data)

}

/** router for simple connection */
func router(conn net.Conn, parsedRequest reqEncoded) {
	useDefaultPattern := true

	r := *routes

	options := ""

	fmt.Println("parsedRequest router ---", parsedRequest)

	/** Determine what function handler to use for processing current connection
	  search appropriate key from RouterMap that suitable for pattern
	*/
	for k, v := range r {

		if k == parsedRequest.Pattern {
			v(conn, options, parsedRequest.Data)

			useDefaultPattern = false
			break
		}
	}

	if useDefaultPattern {
		conn.Write([]byte("Unknowen pattern"))
	}
}

/** router for stream connection */
func streamRouter(conn net.Conn, pattern string, optionsS string, data []byte) {
	useDefaultPattern := true

	r := *routes

	fmt.Println("Pattern ---", pattern, optionsS)

	/** Determine what function handler to use for processing current connection
	  search appropriate key from RouterMap that suitable for pattern
	*/
	for k, v := range r {

		if k == pattern {
			v(conn, optionsS, data)

			useDefaultPattern = false
			break
		}
	}

	if useDefaultPattern {
		conn.Write([]byte("Unknowen pattern"))
	}
}
