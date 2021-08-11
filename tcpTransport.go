package transport

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type resS struct {
	Pattern string
	Data    string
}

func TCPRouter(ln net.Listener, routes map[string]func(conn net.Conn, data string)) {
	for {
		conn, err := ln.Accept()

		fmt.Println("TCPRouter conn - ", conn)

		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn, routes)
	}
}

func handleConnection(conn net.Conn, routes map[string]func(conn net.Conn, data string)) {
	defer conn.Close()

	buf := make([]byte, 10000000)

	n, err := conn.Read(buf)

	if err != nil {
		fmt.Println("conn.Read error", err)
	}

	println("handleConnection data from client 0 -", n)

	jsonRequest := string(buf[:n])
	println("HandleConnection data from client 1 -", jsonRequest, buf)

	// i := strings.Index(jsonRequest, "#")

	var parsedRequest resS

	// println("Data from client 2---", i, jsonRequest[i+1:], string(buf[:n]))

	// If use nest.js ClientPoxy
	// err = json.Unmarshal([]byte(jsonRequest[i+1:]), &parsedRequest)
	err = json.Unmarshal([]byte(jsonRequest), &parsedRequest)

	if err != nil {
		fmt.Println("Request parsing error ---", err)
	}

	fmt.Println("handleConnection parsedRequest ---", parsedRequest)

	router(conn, routes, parsedRequest)
}

func router(conn net.Conn, r map[string]func(conn net.Conn, data string), parsedRequest resS) {
	useDefaultPattern := true

	for k, v := range r {

		if k == parsedRequest.Pattern {
			v(conn, parsedRequest.Data)

			useDefaultPattern = false
			break
		}
	}

	if useDefaultPattern {
		conn.Write([]byte("Unknowen pattern"))
	}
}
